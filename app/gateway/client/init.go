package client

import (
	"go-micro.dev/v4"

	"star/proto/sendSms/sendSmsPb"
	"star/proto/user/userPb"
)

var (
	userService    userPb.UserService
	sendSmsService sendSmsPb.SendMsgService
)

func Init() {
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name("UserService.client"))
	userService = userPb.NewUserService("UserService", userMicroService.Client())

	//创建一个发送短信微服务客户端
	sendSmsMicroService := micro.NewService(micro.Name("UserService.client"))
	sendSmsService = sendSmsPb.NewSendMsgService("SendSmsService", sendSmsMicroService.Client())

}
