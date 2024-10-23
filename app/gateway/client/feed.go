package client

import (
	"context"
	"star/proto/feed/feedPb"
)

func QueryPostExist(ctx context.Context, req *feedPb.QueryPostExistRequest) (*feedPb.QueryPostExistResponse, error) {
	return feedService.QueryPostExist(ctx, req)
}


func GetCommunityPostByTime(ctx context.Context,req *feedPb.GetCommunityPostByTimeRequest)(*feedPb.GetCommunityPostByTimeResponse,error){
	 return  feedService.GetCommunityPostByTime(ctx,req)
}

func GetCommunityPostByNewRely(ctx context.Context,req *feedPb.GetCommunityPostByNewReplyRequest)(*feedPb.GetCommunityPostByNewReplyResponse,error){
	return  feedService.GetCommunityPostByNewReply(ctx,req)
}

func GetPostByRelation(ctx context.Context,req *feedPb.GetPostByRelationRequest)(*feedPb.GetPostByRelationResponse,error){
	return  feedService.GetPostByRelation(ctx,req)
}


