package service

import (
	"context"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"log"
	"star/app/comment/rabbitMQ"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/constant/str"
	"star/models"
	"star/proto/comment/commentPb"
	"star/proto/post/postPb"
	"star/utils"
	"strconv"
)

type CommentService struct {
	commentPb.CommentService
}

var (
	commentSrvIns *CommentService
	postService   postPb.PostService
)

func (s *CommentService) New() {
	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
	postService = postPb.NewPostService(str.PostService, postMicroService.Client())
}

// PostComment 发布评论
func (s *CommentService) PostComment(ctx context.Context, req *commentPb.PostCommentRequest, rsp *commentPb.PostCommentResponse) error {
	comment := &models.Comment{
		PostId:      req.PostId,
		UserId:      req.UserId,
		Content:     req.Content,
		BeCommentId: req.BeCommentId,
	}
	// 存储评论
	if err := mysql.CreateComment(comment); err != nil {
		return err
	}
	rsp.Content = comment.Content

	return nil
}

// StarComment 点赞评论
func (s *CommentService) StarComment(ctx context.Context, req *commentPb.StarCommentRequest, rsp *commentPb.StarCommentResponse) error {
	// 检查评论是否存在
	if err := mysql.CheckComment(req.CommentId); err != nil {
		zap.L().Error("点赞评论: 检查评论是否存在返回错误", zap.Error(err))
		return err
	}
	// 尝试从Redis中获取点赞数
	star, err := redis.GetCommentStar(req.CommentId)
	// 缓存未命中或已过期
	if err != nil || star == 0 {
		// 如果Redis获取失败或返回0，则从MySQL中获取点赞数
		if dbStar, err := mysql.GetStar(req.CommentId); err != nil {
			return err
		} else {
			// 将从MySQL获取的点赞数更新回Redis缓存
			err = redis.SetCommentStar(req.CommentId, dbStar)
			if err != nil {
				log.Println("更新点赞数至Redis失败", err)
			}
			rsp.Star = dbStar
		}
	} else {
		rsp.Star = star
	}

	// Redis中点赞
	if err := redis.IncrementCommentStar(req.CommentId); err != nil {
		log.Println("Redis中点赞失败", err)
		return str.ErrCommentError
	}

	// 使用RabbitMQ异步存储至MySQL数据库中
	// 生产者发布点赞消息
	if err := rabbitMQ.PublishStarEvent(req.CommentId); err != nil {
		if err := redis.Client.Del(ctx, fmt.Sprintf("comment:star:%d", req.CommentId)).Err(); err != nil {
			log.Println("删除Redis中点赞数缓存失败", err)
		}
		return str.ErrCommentError
	}

	rsp.Star = rsp.Star + 1
	return nil
}

// GetComments 获取一个帖子的评论
// 根据页面获取，第几页，每一页多少个评论
func (s *CommentService) GetComments(ctx context.Context, req *commentPb.GetCommentsRequest, rsp *commentPb.GetCommentsResponse) error {
	// 检查帖子是否存在
	postExistResp, err := postService.QueryPostExist(ctx, &postPb.QueryPostExistRequest{PostId: req.PostId})
	if err != nil {
		utils.Logger.Error("query post exist error", zap.Error(err), zap.Int64("post_id", req.PostId))
		return err
	}
	if !postExistResp.Exist {
		utils.Logger.Error("post not exist", zap.Int64("post_id", req.PostId))
		return str.ErrPostNotExists
	}
	// 按照点赞数排序，获取所有评论
	comments, err := mysql.GetCommentsStar(req.PostId)
	if err != nil {
		return str.ErrCommentError
	}

	// 构建评论树
	commentTree, err := buildCommentTree(comments)
	if err != nil {
		log.Println("构建评论树失败")
		return str.ErrCommentError
	}

	// 将评论树转换为 protobuf 格式
	rsp.Comments = convertCommentsToPB(commentTree)

	return nil
}

// buildCommentTree 构建评论树
func buildCommentTree(comments []*models.Comment) ([]*models.Comment, error) {
	// 创建一个map，用于存储所有评论
	commentsMap := make(map[int64]*models.Comment)
	for _, comment := range comments {
		commentsMap[comment.CommentId] = comment
	}

	// 创建一个切片，用于存储顶级评论
	var rootComments []*models.Comment

	// 遍历所有评论，构建树结构
	for _, comment := range comments {
		if comment.BeCommentId == 0 {
			// 如果是顶级评论，直接添加到顶级评论切片中
			rootComments = append(rootComments, comment)
		} else if parent, ok := commentsMap[comment.BeCommentId]; ok {
			// 如果有父评论，将当前评论添加到父评论的子评论列表
			if parent.ChildComments == nil {
				parent.ChildComments = []*models.Comment{}
			}
			parent.ChildComments = append(parent.ChildComments, comment)
		}
	}

	return rootComments, nil
}

// convertCommentsToPB 将评论树转换为protobuf格式
func convertCommentsToPB(comments []*models.Comment) []*commentPb.Comment {
	var pbComments []*commentPb.Comment

	for _, comment := range comments {
		pbComment := &commentPb.Comment{
			CommentId:   comment.CommentId,
			PostId:      comment.PostId,
			UserId:      comment.UserId,
			Content:     comment.Content,
			Star:        int64(comment.Star),
			BeCommentId: comment.BeCommentId,
			Reply:       comment.Reply,
			CreatedAt:   comment.CreatedAt,
		}

		if comment.ChildComments != nil {
			pbComment.ChildComments = convertCommentsToPB(comment.ChildComments)
		}

		pbComments = append(pbComments, pbComment)
	}

	return pbComments
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(ctx context.Context, req *commentPb.DeleteCommentRequest, rsp *commentPb.DeleteCommentResponse) error {
	// 检查评论是否存在
	if err := mysql.CheckComment(req.CommentId); err != nil {
		return err
	}

	// 清除对应评论的Redis缓存
	if err := redis.Client.Del(ctx, fmt.Sprintf("comment:star:%d", req.CommentId)).Err(); err != nil {
		log.Println("删除Redis中点赞数缓存失败", err)
	}

	// 使用RabbitMQ异步在MySQL数据库中删除
	// 生产者发布删除消息
	if err := rabbitMQ.PublishDeleteEvent(req.CommentId); err != nil {
		return err
	}
	return nil
}

func (s *CommentService) CountComment(ctx context.Context, req *commentPb.CountCommentRequest, resp *commentPb.CountCommentResponse) error {
	key := fmt.Sprintf("CountComment:%d", req.PostId)
	countStr, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		return mysql.CountComment(req.PostId)
	})
	if err != nil {
		utils.Logger.Error("get count comment error", zap.Error(err), zap.Int64("post_id", req.PostId), zap.Int64("actorId", req.ActorId))
		return str.ErrCommentError
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		utils.Logger.Error("parse comment count error", zap.Error(err), zap.Int64("post_id", req.PostId), zap.Int64("actorId", req.ActorId), zap.String("countStr", countStr))
		return str.ErrCommentError
	}
	resp.Count = count
	return nil
}

func (s *CommentService) QueryCommentExist(ctx context.Context, req *commentPb.QueryCommentExistRequest, resp *commentPb.QueryCommentExistResponse) error {
	err := mysql.CheckComment(req.CommentId)
	if err != nil {
		if errors.Is(err, str.ErrCommentNotExists) {
			resp.Result = false
			return nil
		}
		return str.ErrCommentError
	}
	resp.Result = true
	return nil
}
