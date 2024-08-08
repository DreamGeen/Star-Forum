package httpHandler

import (
	"go.uber.org/zap"
	logger "star/app/gateway/logger"
	"unicode"

	"github.com/gin-gonic/gin"

	"star/app/gateway/client"
	"star/models"
	"star/proto/user/userPb"
	"star/utils"
)

func SignupHandler(c *gin.Context) {
	//参数校验
	u := new(models.SignupUser)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.GatewayLogger.Error("invalid param", zap.Error(err))
		utils.ResponseMessage(c, utils.CodeInvalidParam)
		return
	}
	//校验用户名
	if err := validateUsername(u.Username); err != nil {
		utils.ResponseErr(c, err)
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
		logger.GatewayLogger.Error("注册失败", zap.Error(err))
		utils.ResponseMessage(c, utils.CodeUserExists)
		return
	}
	//返回成功响应
	utils.ResponseMessage(c, utils.CodeSignupSuccess)
}

// validateUsername 校验用户名的长度和开头是否为数字
func validateUsername(username string) error {
	runes := []rune(username)
	if len(runes) > 20 {
		return utils.ErrUsernameMustLess
	}
	if unicode.IsDigit(runes[0]) {
		return utils.ErrUsernameMustLess
	}
	return nil
}
