package client

import (
	"go-micro.dev/v4"
	"star/constant/str"
	"star/proto/comment/commentPb"
	"star/proto/community/communityPb"
	"star/proto/message/messagePb"
	"star/proto/post/postPb"
	"star/proto/relation/relationPb"
	"star/proto/sendSms/sendSmsPb"
	"star/proto/user/userPb"
)

var (
	userService      userPb.UserService
	sendSmsService   sendSmsPb.SendMsgService
	commentService   commentPb.CommentService
	communityService communityPb.CommunityService
	postService      postPb.PostService
	messageService   messagePb.MessageService
	relationService  relationPb.RelationService
)

func Init() {
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	//创建一个发送短信微服务客户端
	sendSmsMicroService := micro.NewService(micro.Name(str.SendSmsServiceClient))
	sendSmsService = sendSmsPb.NewSendMsgService(str.SendSmsService, sendSmsMicroService.Client())

	// 创建一个评论微服务客户端
	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

	//创建一个社区微服务客户端
	communityMicroService := micro.NewService(micro.Name(str.CommunityServiceClient))
	communityService = communityPb.NewCommunityService(str.CommunityService, communityMicroService.Client())

	//创建一个帖子服务客户端
	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
	postService = postPb.NewPostService(str.PostService, postMicroService.Client())

	//创建一个消息服务客户端
	messageMicroService := micro.NewService(micro.Name(str.MessageServiceClient))
	messageService = messagePb.NewMessageService(str.MessageService, messageMicroService.Client())

	//创建一个关系微服务客户端
	relationMicroService := micro.NewService(micro.Name(str.RelationServiceClient))
	relationService = relationPb.NewRelationService(str.RelationService, relationMicroService.Client())
}
