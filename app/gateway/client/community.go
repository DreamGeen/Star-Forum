package client

import (
	"context"
	"star/proto/community/communityPb"
)

func CreateCommunity(ctx context.Context, req *communityPb.CreateCommunityRequest) (*communityPb.EmptyCommunityResponse, error) {
	return communityService.CreateCommunity(ctx, req)
}

func GetCommunityList(ctx context.Context, req *communityPb.EmptyCommunityRequest) (*communityPb.GetCommunityListResponse, error) {
	return communityService.GetCommunityList(ctx, req)
}

func ShowCommunity(ctx context.Context, req *communityPb.ShowCommunityRequest) (*communityPb.ShowCommunityResponse, error) {
	return communityService.ShowCommunity(ctx, req)
}
