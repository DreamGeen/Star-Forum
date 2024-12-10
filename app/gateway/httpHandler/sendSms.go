package httpHandler

import (
	"github.com/gin-gonic/gin"
	"star/app/constant/settings"
	"star/app/constant/str"
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
		str.Response(c, str.ErrPhoneEmpty, nil)
		return
	}
	if err := sendSms.HandleSendSms(c, phone, templateCode); err != nil {
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
}
