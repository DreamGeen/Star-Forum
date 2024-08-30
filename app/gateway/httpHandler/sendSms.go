package httpHandler

import (
	"github.com/gin-gonic/gin"
	"star/app/gateway/client"
	"star/constant/settings"
	"star/constant/str"
	"star/proto/sendSms/sendSmsPb"
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
	req, empty := validatePhone(c, templateCode)
	if empty {
		str.Response(c, str.ErrPhoneEmpty, str.Empty, nil)
		return
	}
	if _, err := client.HandleSendSms(c, req); err != nil {
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
}

// 验证手机号是否为空
func validatePhone(c *gin.Context, templateCode string) (*sendSmsPb.SendRequest, bool) {
	phone := c.Query("phone")
	return &sendSmsPb.SendRequest{
		Phone:        phone,
		TemplateCode: templateCode,
	}, phone == ""
}
