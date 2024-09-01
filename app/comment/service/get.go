package service

import (
	"context"
	"fmt"
	"star/app/comment/dao/mysql"
	"star/models"
	"star/proto/comment/commentPb"
)

// GetComments 获取一个帖子的评论
// 根据页面获取，第几页，每一页多少个评论
func (s *CommentService) GetComments(ctx context.Context, req *commentPb.GetCommentsRequest, rsp *commentPb.GetCommentsResponse) error {
	// 检查帖子是否存在
	if err := mysql.CheckPost(req.PostId); err != nil {
		return err
	}

	// 按照点赞数排序，获取所有评论
	comments, err := mysql.GetCommentsStar(req.PostId)
	if err != nil {
		return fmt.Errorf("获取评论失败, err: %v", err)
	}

	// 构建评论树
	commentTree, err := buildCommentTree(comments)
	if err != nil {
		return fmt.Errorf("构建评论树失败, err: %v", err)
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
			Star:        comment.Star,
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
