package client

import (
	"context"
	"star/proto/feed/feedPb"
)

func QueryPostExist(ctx context.Context, req *feedPb.QueryPostExistRequest) (*feedPb.QueryPostExistResponse, error) {
	return feedService.QueryPostExist(ctx, req)
}
