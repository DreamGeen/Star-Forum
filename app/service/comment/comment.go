package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis_rate/v10"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/proto/comment/commentPb"
	"star/proto/feed/feedPb"
	"strconv"
)

type CommentService struct {
	commentPb.CommentService
}

const redisCommentQPS = 3

var (
	commentSrvIns *CommentService
	postService   feedPb.FeedService
)

func (s *CommentService) New() {
	postMicroService := micro.NewService(micro.Name(str.FeedServiceClient))
	postService = feedPb.NewFeedService(str.FeedService, postMicroService.Client())
}

func commentLimitKey(userId int64) string {
	return fmt.Sprintf("redis_post_limiter:%d", userId)
}

// PostComment 发布评论
func (s *CommentService) PostComment(ctx context.Context, req *commentPb.PostCommentRequest, rsp *commentPb.PostCommentResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "PostCommentService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CommentService.PostComment")

	//redis limit
	limiter := redis_rate.NewLimiter(redis.Client)
	limiterKey := commentLimitKey(req.UserId)
	limitRes, err := limiter.Allow(ctx, limiterKey, redis_rate.PerSecond(redisCommentQPS))
	if err != nil {
		logger.Error("comment limiter error",
			zap.Error(err),
			zap.Int64("actorId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	if limitRes.Allowed == 0 {
		logger.Error("user feed comment too frequently",
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrRequestTooFrequently
	}

	comment := &models.Comment{
		PostId:      req.PostId,
		UserId:      req.UserId,
		Content:     req.Content,
		BeCommentId: req.BeCommentId,
	}
	// 存储评论
	if err := mysql.CreateComment(comment); err != nil {
		logger.Error("feed service error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.String("content", req.Content),
			zap.Int64("BeCommentId", req.BeCommentId))
		logging.SetSpanError(span, err)
		return err
	}
	key := fmt.Sprintf("CountComment:%d", req.PostId)
	cached.Delete(ctx, key)
	rsp.Content = comment.Content

	return nil
}

// GetComments 获取一个帖子的评论
// 根据页面获取，第几页，每一页多少个评论
func (s *CommentService) GetComments(ctx context.Context, req *commentPb.GetCommentsRequest, rsp *commentPb.GetCommentsResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetCommentsService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CommentService.GetComments")

	// 检查帖子是否存在
	postExistResp, err := postService.QueryPostExist(ctx,
		&feedPb.QueryPostExistRequest{
			PostId: req.PostId,
		})
	if err != nil {
		logger.Error("query feed exist error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId))
		logging.SetSpanError(span, err)
		return err
	}
	if !postExistResp.Exist {
		logger.Error("feed not exist",
			zap.Int64("post_id", req.PostId))
		logging.SetSpanError(span, err)
		return str.ErrPostNotExists
	}
	// 按照点赞数排序，获取所有评论
	comments, err := mysql.GetCommentsStar(req.PostId)
	if err != nil {
		logger.Error("mysql get comment star error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCommentError
	}

	// 构建评论树
	commentTree, err := buildCommentTree(comments)
	if err != nil {
		logger.Error("create comment tree error",
			zap.Error(err))
		logging.SetSpanError(span, err)
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
	ctx, span := tracing.Tracer.Start(ctx, "DeleteCommentService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CommentService.DeleteComment")

	// 检查评论是否存在
	if err := mysql.CheckComment(req.CommentId); err != nil {
		logger.Error("check comment error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCommentError
	}

	// 清除对应评论的Redis缓存
	if err := redis.Client.Del(ctx, fmt.Sprintf("comment:star:%d", req.CommentId)).Err(); err != nil {
		logger.Error("delete redis comment buffer error",
			zap.Error(err),
			zap.Int64("commentId", req.CommentId))
		logging.SetSpanError(span, err)
	}

	//删除评论
	if err := mysql.DeleteComment(req.CommentId); err != nil {
		logger.Error("delete comment error",
			zap.Error(err),
			zap.Int64("commentId", req.CommentId))
		logging.SetSpanError(span, err)
		return str.ErrCommentError
	}
	return nil
}

func (s *CommentService) CountComment(ctx context.Context, req *commentPb.CountCommentRequest, resp *commentPb.CountCommentResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CountCommentService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CommentService.CountComment")

	key := fmt.Sprintf("CountComment:%d", req.PostId)
	countStr, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		return mysql.CountComment(req.PostId)
	})
	if err != nil {
		logger.Error("get count comment error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId),
			zap.Int64("actorId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrCommentError
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logger.Error("parse comment count error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId),
			zap.Int64("actorId", req.ActorId),
			zap.String("countStr", countStr))
		logging.SetSpanError(span, err)
		return str.ErrCommentError
	}
	resp.Count = count
	return nil
}

func (s *CommentService) QueryCommentExist(ctx context.Context, req *commentPb.QueryCommentExistRequest, resp *commentPb.QueryCommentExistResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "QueryCommentExistService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CommentService.QueryCommentExist")

	err := mysql.CheckComment(req.CommentId)
	if err != nil {
		if errors.Is(err, str.ErrCommentNotExists) {
			resp.Result = false
			return nil
		}
		logger.Error("query comment exist error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCommentError
	}
	resp.Result = true
	return nil
}
