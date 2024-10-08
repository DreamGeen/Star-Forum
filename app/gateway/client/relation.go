package client

import (
	"context"
	"star/proto/relation/relationPb"
)

func GetFollowList(ctx context.Context, req *relationPb.GetFollowRequest) (*relationPb.GetFollowResponse, error) {
	return relationService.GetFollowList(ctx, req)
}

func Follow(ctx context.Context, req *relationPb.FollowRequest) (*relationPb.FollowResponse, error) {
	return relationService.Follow(ctx, req)
}

func UnFollow(ctx context.Context, req *relationPb.UnFollowRequest) (*relationPb.UnFollowResponse, error) {
	return relationService.UnFollow(ctx, req)
}

func GetFansList(ctx context.Context, req *relationPb.GetFansListRequest) (*relationPb.GetFansListResponse, error) {
	return relationService.GetFansList(ctx, req)
}
