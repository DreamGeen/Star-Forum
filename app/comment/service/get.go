package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/proto/comment/commentPb"
)

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
			CreatedAt:   comment.CreatedAt,
		})
	}
	return nil
}
