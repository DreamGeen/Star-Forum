package client

import (
	"go-micro.dev/v4"
	"star/constant/str"
	"star/proto/sendSms/sendSmsPb"
	"star/proto/user/userPb"
)

var (
	userServiceClient    userPb.UserService
	sendSmsServiceClient sendSmsPb.SendMsgService
)

func Init() {
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userServiceClient = userPb.NewUserService(str.UserService, userMicroService.Client())

	//创建一个发送短信微服务客户端
	sendSmsMicroService := micro.NewService(micro.Name(str.SendSmsServiceClient))
	sendSmsServiceClient = sendSmsPb.NewSendMsgService(str.SendSmsService, sendSmsMicroService.Client())

}
