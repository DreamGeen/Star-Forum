package httpHandler

import (
	"github.com/gin-gonic/gin"
	"star/app/constant/settings"
	str2 "star/app/constant/str"
	"star/app/utils/sendSms"
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
	phone := c.Query("phone")
	if phone == "" {
		str2.Response(c, str2.ErrPhoneEmpty, str2.Empty, nil)
		return
	}
	if err := sendSms.HandleSendSms(c, phone, templateCode); err != nil {
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
}
