package client

import (
	"context"
	"star/proto/comment/commentPb"
)

func PostComment(ctx context.Context, req *commentPb.PostCommentRequest) (*commentPb.PostCommentResponse, error) {
	return commentService.PostComment(ctx, req)
}

func DeleteComment(ctx context.Context, req *commentPb.DeleteCommentRequest) (*commentPb.DeleteCommentResponse, error) {
	return commentService.DeleteComment(ctx, req)
}

func GetComments(ctx context.Context, req *commentPb.GetCommentsRequest) (*commentPb.GetCommentsResponse, error) {
	return commentService.GetComments(ctx, req)
}

func StarComment(ctx context.Context, req *commentPb.StarCommentRequest) (*commentPb.StarCommentResponse, error) {
	return commentService.StarComment(ctx, req)
}
