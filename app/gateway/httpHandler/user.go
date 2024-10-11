package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/proto/user/userPb"
	"unicode"

	"star/app/gateway/client"
)

// LoginHandler 用户名或手机号或邮箱和密码进行登录
func LoginHandler(c *gin.Context) {
	//参数校验
	u := new(models.LoginPassword)
	if err := c.ShouldBindJSON(u); err != nil {
		logging.Logger.Error("login error invalid param",
			zap.Error(err))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		User:     u.User,
		Password: u.Password,
	}
	resp, err := client.LoginPassword(c, req)
	if err != nil {
		logging.Logger.Error("login error",
			zap.Error(err),
			zap.String("user", req.User),
			zap.String("password", req.Password))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	//成功响应
	str2.Response(c, nil, "token", resp.Token)
}

// LoginWithCaptcha 手机验证码登录
func LoginWithCaptcha(c *gin.Context) {
	//参数校验
	u := new(models.LoginCaptcha)
	if err := c.ShouldBindJSON(u); err != nil {
		logging.Logger.Error("login error invalid param",
			zap.Error(err))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		Phone:   u.Phone,
		Captcha: u.Captcha,
	}
	resp, err := client.LoginCaptcha(c, req)
	if err != nil {
		logging.Logger.Error("login error",
			zap.Error(err),
			zap.String("phone", req.Phone),
			zap.String("captcha", req.Captcha))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	//成功响应
	str2.Response(c, nil, "token", resp.Token)
}

func SignupHandler(c *gin.Context) {
	//参数校验
	u := new(models.SignupUser)
	if err := c.ShouldBindJSON(u); err != nil {
		logging.Logger.Error("sign up error invalid param",
			zap.Error(err))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	//校验用户名
	if err := validateUsername(u.Username); err != nil {
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	//注册处理
	req := &userPb.LSRequest{
		User:     u.Username,
		Password: u.Password,
		Phone:    u.Phone,
		Captcha:  u.Captcha,
	}
	if _, err := client.Signup(c, req); err != nil {
		logging.Logger.Error("sign up error",
			zap.String("phone", req.Phone))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	//返回成功响应
	str2.Response(c, nil, str2.Empty, nil)
}

// validateUsername 校验用户名的长度和开头是否为数字
func validateUsername(username string) error {
	runes := []rune(username)
	if len(runes) > 20 {
		return str2.ErrUsernameMustLess
	}
	if unicode.IsDigit(runes[0]) {
		return str2.ErrUsernameStartWith
	}
	return nil
}
