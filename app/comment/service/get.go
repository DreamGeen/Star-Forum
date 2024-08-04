package service

import (
	"context"
	"fmt"
	"star/app/comment/dao/mysql"
	"star/proto/comment/commentPb"
)

// GetComments 获取一个帖子的评论
// 根据页面获取，第几页，每一页多少个评论
func (s *CommentService) GetComments(ctx context.Context, req *commentPb.GetCommentsRequest, rsp *commentPb.GetCommentsResponse) error {
	// 按照点赞数排序
	comments, err := mysql.GetCommentsStar(req.PostId, req.Page, req.PageSize)
	if err != nil {
		return fmt.Errorf("GetComments err: %v", err)
	}
	for _, comment := range comments {
		rsp.Comments = append(rsp.Comments, &commentPb.Comment{
			CommentId:   comment.CommentId,
			PostId:      comment.PostId,
			UserId:      comment.UserId,
			Content:     comment.Content,
			Star:        comment.Star,
			BeCommentId: comment.BeCommentId,
			Reply:       comment.Reply,
			CreatedAt:   comment.CreatedAt,
		})
	}
	return nil
}
