package httpHandler

import (
	"github.com/mojocn/base64Captcha"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/models"
	"star/app/utils/jwt"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/user/userPb"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"star/app/gateway/client"
)

// 图片验证码生成
var store = &models.CaptchaStore{}
var driver = &base64Captcha.DriverMath{
	Height: 42,
	Width:  140,
}

// LoginHandler 用户名或手机号或邮箱和密码进行登录
func LoginHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "LoginHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.Login")

	//参数校验
	u := new(models.LoginPassword)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.Error("login error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	if !store.Verify(u.CheckCodeKey, u.CheckCode, true) {
		logger.Warn("img captcha error",
			zap.String("user", u.User))
		str.Response(c, str.ErrInvalidImgCaptcha, nil)
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
		str.Response(c, err, nil)
		return
	}

	//成功响应
	str.Response(c, nil, map[string]interface{}{
		"accessToken":  resp.Token.AccessToken,
		"refreshToken": resp.Token.RefreshToken,
		"userInfo":     resp.UserInfo,
	})
}

// LoginWithCaptchaHandler 手机验证码登录
func LoginWithCaptchaHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "LoginWithCaptchaHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.LoginWithCaptcha")

	//参数校验
	u := new(models.LoginCaptcha)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.Error("login error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidImgCaptcha, nil)
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
		str.Response(c, err, nil)
		return
	}
	//成功响应
	str.Response(c, nil, map[string]interface{}{
		"token": resp.Token,
	})
}

func AutoLoginHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "AutoLoginHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.AutoLoginHandler")

	oldToken := new(models.Token)
	if err := c.ShouldBindJSON(oldToken); err != nil {
		logger.Error("auto login error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	userId, accessToken, err := jwt.AutoLogin(oldToken.AccessToken, oldToken.RefreshToken)
	if err != nil {
		logger.Error("auto login refresh token error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	userInfoResp, err := client.GetUserInfo(c.Request.Context(), &userPb.GetUserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error("auto login error because get use info error",
			zap.Error(err))
		str.Response(c, str.ErrAutoLoginError, nil)
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"accessToken": accessToken,
		"userInfo":    userInfoResp.User,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "RefreshTokenHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.RefreshToken")

	oldToken := new(models.Token)
	if err := c.ShouldBindJSON(oldToken); err != nil {
		logger.Error("refresh error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	accessToken, err := jwt.RefreshAccessToken(oldToken.AccessToken, oldToken.RefreshToken)
	if err != nil {
		logger.Error("refresh token error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"accessToken": accessToken,
	})
}

func SignupHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "SignupHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.Signup")

	//参数校验
	u := new(models.SignupUser)
	if err := c.ShouldBindJSON(u); err != nil {
		logging.Logger.Error("sign up error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	if !store.Verify(u.CheckCodeKey, u.CheckCode, true) {
		logger.Warn("img captcha error",
			zap.String("phone", u.Phone))
		str.Response(c, str.ErrInvalidImgCaptcha, nil)
		return
	}
	remoteAddr := c.RemoteIP()
	//校验用户名
	if err := validateUsername(u.Username); err != nil {
		str.Response(c, err, nil)
		return
	}
	//注册处理
	req := &userPb.LSRequest{
		User:     u.Username,
		Password: u.Password,
		Phone:    u.Phone,
		Captcha:  u.Captcha,
		Ip:       remoteAddr,
	}
	if _, err := client.Signup(c.Request.Context(), req); err != nil {
		logger.Error("sign up error",
			zap.String("phone", req.Phone))
		str.Response(c, err, nil)
		return
	}
	//返回成功响应
	str.Response(c, nil, nil)
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

func GetUserInfoHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetUserInfoHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetUserInfo")

	userIdStr := c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 64, 10)
	if err != nil || userId == 0 {
		logger.Error("invalid param",
			zap.Error(err),
			zap.String("userIdStr", userIdStr))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	actorId, err := request.GetUserId(c)
	if err != nil {
		actorId = 0
	}
	resp, err := client.GetUserInfo(c.Request.Context(), &userPb.GetUserInfoRequest{
		UserId:  userId,
		ActorId: actorId,
	})
	if err != nil {
		logger.Error("get community info service error",
			zap.Error(err))
		str.Response(c, err, nil)
	}

	str.Response(c, nil, map[string]interface{}{
		"userInfo": resp.User,
	})
}

// GetCaptchaHandler  生成并返回图片验证码
func GetCaptchaHandler(c *gin.Context) {
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, encodedImage, answer, err := captcha.Generate()
	if err != nil {
		str.Response(c, str.ErrServiceBusy, nil)
		return
	}
	if err := store.Set(id, answer); err != nil {
		str.Response(c, str.ErrServiceBusy, nil)
		return
	}

	// 返回Base64编码的图像
	str.Response(c, nil, map[string]interface{}{
		"captcha":   encodedImage,
		"captchaId": id,
	})
}
