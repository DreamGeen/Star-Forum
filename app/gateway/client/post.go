package client

import (
	"context"
	"star/proto/post/postPb"
)

func QueryPostExist(ctx context.Context, req *postPb.QueryPostExistRequest) (*postPb.QueryPostExistResponse, error) {
	return postService.QueryPostExist(ctx, req)
}
