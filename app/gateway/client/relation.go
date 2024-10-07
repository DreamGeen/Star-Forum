package client

import (
	"context"
	"star/proto/relation/relationPb"
)

func GetFriendList(ctx context.Context, req *relationPb.GetFollowerRequest) (*relationPb.GetFollowerResponse, error) {
	return relationService.GetFollowerList(ctx, req)
}
