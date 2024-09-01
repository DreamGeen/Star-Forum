package client

import (
	"go-micro.dev/v4"
	"star/constant/str"
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
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	//创建一个发送短信微服务客户端
	sendSmsMicroService := micro.NewService(micro.Name(str.SendSmsServiceClient))
	sendSmsService = sendSmsPb.NewSendMsgService(str.SendSmsService, sendSmsMicroService.Client())

	// 创建一个评论微服务客户端
	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

}
