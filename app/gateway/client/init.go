package client

import (
	"go-micro.dev/v4"
	"star/proto/comment/commentPb"
	"star/proto/sendSms/sendSmsPb"
	"star/proto/user/userPb"
)

var (
	userService    userPb.UserService
	sendSmsService sendSmsPb.SendMsgService
	commentService commentPb.CommentService
)

func Init() {
	// 创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name("UserService.client"))
	userService = userPb.NewUserService("UserService", userMicroService.Client())

	// 创建一个发送短信微服务客户端
	sendSmsMicroService := micro.NewService(micro.Name("SendSmsService.client"))
	sendSmsService = sendSmsPb.NewSendMsgService("SendSmsService", sendSmsMicroService.Client())

	// 创建一个评论微服务客户端
	commentMicroService := micro.NewService(micro.Name("CommentService.client"))
	commentService = commentPb.NewCommentService("CommentService", commentMicroService.Client())
}
