package client

import (
	"context"
	"star/proto/publish/publishPb"
)

func CreatePost(ctx context.Context, req *publishPb.CreatePostRequest) (*publishPb.CreatePostResponse, error) {
	return publishService.CreatePost(ctx, req)
}

func PreUploadVideos(ctx context.Context, req *publishPb.PreUploadVideosRequest) (*publishPb.PreUploadVideosResponse, error) {
	return publishService.PreUploadVideos(ctx, req)
}
