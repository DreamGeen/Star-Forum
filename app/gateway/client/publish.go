package client

import (
	"context"
	"star/proto/publish/publishPb"
)

func CreatePost(ctx context.Context,req *publishPb.CreatePostRequest)(*publishPb.CreatePostResponse,error){
    return  publishService.CreatePost(ctx,req)
}