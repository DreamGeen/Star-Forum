package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/models"
	"star/proto/comment/commentPb"
)

type CommentService struct{}

func (s *CommentService) PostComment(ctx context.Context, req *commentPb.PostCommentRequest, rsp *commentPb.PostCommentResponse) error {
	comment := &models.Comment{
		PostId:      req.PostId,
		UserId:      req.UserId,
		Content:     req.Content,
		BeCommentId: &req.BeCommentId,
	}
	if err := mysql.CreateComment(comment); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	rsp.Success = true
	rsp.Message = "评论发布成功"
	return nil
}
