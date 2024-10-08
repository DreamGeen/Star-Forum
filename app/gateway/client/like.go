package client

import (
	"context"
	"star/proto/like/likePb"
)

func LikeAction(ctx context.Context, req *likePb.LikeActionRequest) (*likePb.LikeActionResponse, error) {
	return likeService.LikeAction(ctx, req)
}
func LikeList(ctx context.Context, req *likePb.LikeListRequest) (*likePb.LikeListResponse, error) {
	return likeService.LikeList(ctx, req)
}
