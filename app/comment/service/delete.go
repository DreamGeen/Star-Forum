package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/proto/comment/commentPb"
)

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
