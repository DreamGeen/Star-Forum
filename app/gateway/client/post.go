package client

import (
	"context"
	"star/proto/post/postPb"
)

func QueryPostExist(ctx context.Context, req *postPb.QueryPostExistRequest) (*postPb.QueryPostExistResponse, error) {
	return postService.QueryPostExist(ctx, req)
}

func CreatePost(ctx context.Context, req *postPb.CreatePostRequest) (*postPb.CreatePostResponse, error) {
	return postService.CreatePost(ctx, req)
}

func GetPostByPopularity(ctx context.Context, req *postPb.GetPostListByPopularityRequest) (*postPb.GetPostListByPopularityResponse, error) {
	return postService.GetPostByPopularity(ctx, req)
}

func GetPostByTime(ctx context.Context, req *postPb.GetPostListByTimeRequest) (*postPb.GetPostListByTimeResponse, error) {
	return postService.GetPostByTime(ctx, req)
}
