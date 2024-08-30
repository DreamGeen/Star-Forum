package httpHandler

import (
	"github.com/gin-gonic/gin"
	"log"
	"star/app/gateway/models"
	"star/constant/str"
	"star/proto/user/userPb"
	"unicode"

	"star/app/gateway/client"
)

// LoginHandler 用户名或手机号或邮箱和密码进行登录
func LoginHandler(c *gin.Context) {
	//参数校验
	u := new(models.LoginPassword)
	if err := c.ShouldBindJSON(u); err != nil {
		log.Println("参数效验失败", err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		User:     u.User,
		Password: u.Password,
	}
	resp, err := client.LoginPassword(c, req)
	if err != nil {
		log.Println("登录失败", err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	//成功响应
	str.Response(c, nil, "token", resp.Token)
}

// LoginWithCaptcha 手机验证码登录
func LoginWithCaptcha(c *gin.Context) {
	//参数校验
	u := new(models.LoginCaptcha)
	if err := c.ShouldBindJSON(u); err != nil {
		log.Println("参数错误", err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		Phone:   u.Phone,
		Captcha: u.Captcha,
	}
	resp, err := client.LoginCaptcha(c, req)
	if err != nil {
		log.Println("登录失败", err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	//成功响应
	str.Response(c, nil, "token", resp.Token)
}

func SignupHandler(c *gin.Context) {
	//参数校验
	u := new(models.SignupUser)
	if err := c.ShouldBindJSON(u); err != nil {
		log.Println("invalid param", err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	//校验用户名
	if err := validateUsername(u.Username); err != nil {
		str.Response(c, err, str.Empty, nil)
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
		log.Println("注册失败", err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	//返回成功响应
	str.Response(c, nil, str.Empty, nil)
}

// validateUsername 校验用户名的长度和开头是否为数字
func validateUsername(username string) error {
	runes := []rune(username)
	if len(runes) > 20 {
		return str.ErrUsernameMustLess
	}
	if unicode.IsDigit(runes[0]) {
		return str.ErrUsernameStartWith
	}
	return nil
}