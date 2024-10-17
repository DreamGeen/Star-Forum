package client

import (
	"go-micro.dev/v4"
	"star/app/constant/str"
	"star/proto/collect/collectPb"
	"star/proto/comment/commentPb"
	"star/proto/community/communityPb"
	"star/proto/feed/feedPb"
	"star/proto/like/likePb"
	"star/proto/message/messagePb"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
)

var (
	userService      userPb.UserService
	commentService   commentPb.CommentService
	communityService communityPb.CommunityService
	feedService      feedPb.PostService
	messageService   messagePb.MessageService
	relationService  relationPb.RelationService
	likeService      likePb.LikeService
	collectService   collectPb.CollectService
)

func Init() {
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	// 创建一个评论微服务客户端
	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

	//创建一个社区微服务客户端
	communityMicroService := micro.NewService(micro.Name(str.CommunityServiceClient))
	communityService = communityPb.NewCommunityService(str.CommunityService, communityMicroService.Client())

	//创建一个feed服务客户端
	feedMicroService := micro.NewService(micro.Name(str.FeedServiceClient))
	feedService = feedPb.NewPostService(str.FeedService, feedMicroService.Client())

	//创建一个消息服务客户端
	messageMicroService := micro.NewService(micro.Name(str.MessageServiceClient))
	messageService = messagePb.NewMessageService(str.MessageService, messageMicroService.Client())

	//创建一个关系微服务客户端
	relationMicroService := micro.NewService(micro.Name(str.RelationServiceClient))
	relationService = relationPb.NewRelationService(str.RelationService, relationMicroService.Client())

	//创建一个点赞微服务客户端
	likeMicroService := micro.NewService(micro.Name(str.LikeServiceClient))
	likeService = likePb.NewLikeService(str.LikeService, likeMicroService.Client())

	//创建一个收藏微服务客户端
	collectMicroService := micro.NewService(micro.Name(str.CollectServiceClient))
	collectService = collectPb.NewCollectService(str.CollectService, collectMicroService.Client())

}
