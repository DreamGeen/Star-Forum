package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	"star/proto/comment/commentPb"
)

// StarComment 点赞评论
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
