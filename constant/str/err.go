package str

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	UserNotExistsCode = 10001 + iota
	UsernameExistsCode
	UsernameMustLess
	UsernameStartWith
	InvalidPasswordCode
	InvalidParamCode
	InvalidCaptchaCode
	NotLoginCode
	PhoneEmptyCode
	PhoneRegisteredCode
	PhoneUnregisteredCode
)

const (
	ServiceBusyCode = 50001 + iota
	LoginErrorCode
	SignupErrorCode
	SendSmsErrorCode
)

var (
	ErrUserNotExists     = errors.New("用户不存在")
	ErrUsernameExists    = errors.New("用户名已存在")
	ErrUsernameMustLess  = errors.New("用户名长度必须小于20")
	ErrUsernameStartWith = errors.New("用户名不能以数字开头")
	ErrInvalidPassword   = errors.New("用户名或密码错误")
	ErrInvalidParam      = errors.New("请求参数错误")
	ErrInvalidCaptcha    = errors.New("验证码错误")
	ErrNotLogin          = errors.New("请登录")
	ErrPhoneEmpty        = errors.New("手机号为空")
	ErrPhoneRegistered   = errors.New("手机号已注册")
	ErrPhoneUnregistered = errors.New("该手机号未注册")
)

var (
	ErrLoginError   = errors.New("登录服务出现内部错误，请稍后再试！")
	ErrSignupError  = errors.New("注册服务出现内部错误，请稍后再试！")
	ErrSendSmsError = errors.New("发送短信失败，请稍后再试！")
)

var codeMap = map[error]int32{
	ErrUserNotExists:     UserNotExistsCode,
	ErrUsernameExists:    UsernameExistsCode,
	ErrUsernameMustLess:  UsernameMustLess,
	ErrUsernameStartWith: UsernameStartWith,
	ErrInvalidPassword:   InvalidPasswordCode,
	ErrInvalidParam:      InvalidParamCode,
	ErrInvalidCaptcha:    InvalidCaptchaCode,
	ErrNotLogin:          NotLoginCode,
	ErrPhoneEmpty:        PhoneEmptyCode,
	ErrPhoneRegistered:   PhoneRegisteredCode,
	ErrPhoneUnregistered: PhoneUnregisteredCode,

	ErrLoginError:   LoginErrorCode,
	ErrSignupError:  SignupErrorCode,
	ErrSendSmsError: SendSmsErrorCode,
}

func getCode(err error) int32 {
	code, ok := codeMap[err]
	if !ok {
		return ServiceBusyCode
	}
	return code
}

func Response(c *gin.Context, err error, dataFieldName string, data interface{}) {
	statusMsg := Success
	if err != nil {
		statusMsg = err.Error()
	}
	statusCode := getCode(err)
	// 构建响应 JSON
	response := gin.H{
		"statusCode": statusCode,
		"statusMsg":  statusMsg,
	}

	// 根据传入的 dataFieldName 动态设置字段名
	if dataFieldName != "" {
		response[dataFieldName] = data
	} else {
		response["data"] = data // 默认使用 "data" 字段名
	}
	c.JSON(http.StatusOK, response)
	return
}
