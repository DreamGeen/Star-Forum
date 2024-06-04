package rpc

import (
	"go-micro.dev/v4"

	"star/proto/sendSms/sendSmsPb"
	"star/proto/user/userPb"
)

var (
	UserService    userPb.UserService
	SendSmsService sendSmsPb.SendMsgService
)

func Init() {
	//创建一个用户微服务客户端
	UserMicroService := micro.NewService(micro.Name("UserService.client"))
	UserService = userPb.NewUserService("UserService", UserMicroService.Client())

	//创建一个发送短信微服务客户端
	SendSmsMicroService := micro.NewService(micro.Name("UserService.client"))
	SendSmsService = sendSmsPb.NewSendMsgService("SendSmsService", SendSmsMicroService.Client())

}
