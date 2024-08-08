package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	logger "star/app/gateway/logger"

	"star/app/gateway/client"
	"star/models"
	"star/proto/user/userPb"
	"star/utils"
)

// LoginHandler 用户名或手机号或邮箱和密码进行登录
func LoginHandler(c *gin.Context) {
	//参数校验
	u := new(models.LoginPassword)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.GatewayLogger.Error("参数错误", zap.Error(err))
		utils.ResponseMessage(c, utils.CodeInvalidParam)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		User:     u.User,
		Password: u.Password,
	}
	resp, err := client.LoginPassword(c, req)
	if err != nil {
		logger.GatewayLogger.Error("登录失败", zap.Error(err))
		utils.ResponseErr(c, err)
		return
	}
	//成功响应
	utils.ResponseMessageWithData(c, utils.CodeLoginSuccess, resp.Token)
}

// LoginWithCaptcha 手机验证码登录
func LoginWithCaptcha(c *gin.Context) {
	//参数校验
	u := new(models.LoginCaptcha)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.GatewayLogger.Error("参数错误", zap.Error(err))
		utils.ResponseMessage(c, utils.CodeInvalidParam)
		return
	}
	//登录处理
	req := &userPb.LSRequest{
		Phone:   u.Phone,
		Captcha: u.Captcha,
	}
	resp, err := client.LoginCaptcha(c, req)
	if err != nil {
		logger.GatewayLogger.Error("登录失败", zap.Error(err))
		utils.ResponseErr(c, err)
		return
	}
	//成功响应
	utils.ResponseMessageWithData(c, utils.CodeLoginSuccess, resp.Token)
}
