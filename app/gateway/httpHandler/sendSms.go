package httpHandler

import (
	"log"

	"github.com/gin-gonic/gin"

	"star/app/gateway/client"
	"star/proto/sendSms/sendSmsPb"
	"star/settings"
	"star/utils"
)

// SendSetupHandler 发送注册短信
func SendSetupHandler(c *gin.Context) {
	sendHandler(c, settings.Conf.SignupTemplateCode)
}

// SendLoginHandler 发送登录短信
func SendLoginHandler(c *gin.Context) {
	sendHandler(c, settings.Conf.LoginTemplateCode)
}

// 发送短信处理
func sendHandler(c *gin.Context, templateCode string) {
	req, err := validatePhone(c, templateCode)
	if err != nil {
		utils.ResponseErr(c, utils.ErrPhoneEmpty)
		return
	}
	if _, err := client.HandleSendSms(c, req); err != nil {
		utils.ResponseErr(c, err)
		return
	}
	utils.ResponseMessage(c, utils.CodeSendSmsSuccess)
}

// 验证手机号是否为空
func validatePhone(c *gin.Context, templateCode string) (*sendSmsPb.SendRequest, error) {
	phone := c.Query("phone")
	if phone == "" {
		log.Println("手机号为空")
		return nil, utils.ErrPhoneEmpty
	}
	return &sendSmsPb.SendRequest{
		Phone:        phone,
		TemplateCode: templateCode,
	}, nil
}
