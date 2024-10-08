package client

import (
	"context"
	"star/proto/collect/collectPb"
)

func CollectList(ctx context.Context, req *collectPb.CollectListRequest) (*collectPb.CollectListResponse, error) {
	return collectService.CollectList(ctx, req)
}

func CollectAction(ctx context.Context, req *collectPb.CollectActionRequest) (*collectPb.CollectActionResponse, error) {
	return collectService.CollectAction(ctx, req)
}
