package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/user/userPb"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"star/app/gateway/client"
)

// LoginHandler 用户名或手机号或邮箱和密码进行登录
func LoginHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"LoginHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.Login")
	
	//参数校验
	u := new(models.LoginPassword)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.Error("login error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		User:     u.User,
		Password: u.Password,
	}
	resp, err := client.LoginPassword(c.Request.Context(), req)
	if err != nil {
		logger.Error("login error",
			zap.Error(err),
			zap.String("user", req.User),
			zap.String("password", req.Password))
		str.Response(c, err, str.Empty, nil)
		return
	}
	//成功响应
	str.Response(c, nil, "token", resp.Token)
}

// LoginWithCaptchaHandler 手机验证码登录
func LoginWithCaptchaHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"LoginWithCaptchaHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.LoginWithCaptcha")
	
	//参数校验
	u := new(models.LoginCaptcha)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.Error("login error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		Phone:   u.Phone,
		Captcha: u.Captcha,
	}
	resp, err := client.LoginCaptcha(c.Request.Context(), req)
	if err != nil {
		logger.Error("login error",
			zap.Error(err),
			zap.String("phone", req.Phone),
			zap.String("captcha", req.Captcha))
		str.Response(c, err, str.Empty, nil)
		return
	}
	//成功响应
	str.Response(c, nil, "token", resp.Token)
}

func SignupHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"SignupHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.Signup")

	//参数校验
	u := new(models.SignupUser)
	if err := c.ShouldBindJSON(u); err != nil {
		logging.Logger.Error("sign up error invalid param",
			zap.Error(err))
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
	if _, err := client.Signup(c.Request.Context(), req); err != nil {
		logger.Error("sign up error",
			zap.String("phone", req.Phone))
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
func GetUserInfoHandler(c *gin.Context){
	_,span:=tracing.Tracer.Start(c.Request.Context(),"GetUserInfoHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.GetUserInfo")

    userIdStr:=c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 64, 10)
	if err != nil || userId == 0 {
		logger.Error("invalid param",
			zap.Error(err),
			zap.String("userIdStr", userIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	actorId,err:=request.GetUserId(c)
	if err!=nil{
		actorId=0
	}
    resp,err:=client.GetUserInfo(c.Request.Context(),&userPb.GetUserInfoRequest{
		UserId: userId,
		ActorId: actorId,
	})
	if err!=nil{
		logger.Error("get community info service error",
		   zap.Error(err)) 
        str.Response(c,err,str.Empty,nil)
	}
	str.Response(c,nil,"userInfo",resp.User)
}
