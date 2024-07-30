package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
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

func (s *CommentService) DeleteComment(ctx context.Context, req *commentPb.DeleteCommentRequest, rsp *commentPb.DeleteCommentResponse) error {
	if err := mysql.DeleteComment(req.CommentId); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	rsp.Success = true
	rsp.Message = "评论删除成功"
	return nil
}

func (s *CommentService) GetComments(ctx context.Context, req *commentPb.GetCommentsRequest, rsp *commentPb.GetCommentsResponse) error {
	comments, err := mysql.GetComments(req.PostId, req.Page, req.PageSize)
	if err != nil {
		return err
	}
	for _, comment := range comments {
		rsp.Comments = append(rsp.Comments, &commentPb.Comment{
			CommentId:   comment.CommentId,
			PostId:      comment.PostId,
			UserId:      comment.UserId,
			Content:     comment.Content,
			Star:        comment.Star,
			BeCommentId: *comment.BeCommentId,
		})
	}
	return nil
}

func (s *CommentService) StarComment(ctx context.Context, req *commentPb.StarCommentRequest, rsp *commentPb.StarCommentResponse) error {
	if err := redis.IncrementCommentStar(req.CommentId); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	if err := mysql.UpdateStar(req.CommentId, 1); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	rsp.Success = true
	rsp.Message = "点赞成功"
	return nil
}
