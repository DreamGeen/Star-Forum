package utils

import "errors"

type ResponseCode int64

// 定义响应码
const (
	CodeSuccess ResponseCode = 1000 + iota
	CodeLoginSuccess
	CodeSignupSuccess
	CodeSendSmsSuccess

	CodeUserExists = 2000 + iota
	CodeUserNotExists
	CodeUsernameStartWithNumber
	CodeUsernameMustLess
	CodeInvalidPassword
	CodeInvalidParam
	CodeGetCaptchaFailed
	CodeInvalidCaptcha
	CodeServiceBusy
	CodeNotLogin
	CodeSendSmsFailed
	CodePhoneEmpty
)

// 定义错误
var (
	ErrUserNotExists           = errors.New("用户名或密码错误")
	ErrServiceBusy             = errors.New("服务繁忙")
	ErrUsernameStartWithNumber = errors.New("用户名不能以数字字符开头")
	ErrUsernameMustLess        = errors.New("用户名长度必须小于20")
	ErrGetCaptchaFailed        = errors.New("获取验证码失败")
	ErrInvalidCaptcha          = errors.New("验证码错误或失效")
	ErrNotLogin                = errors.New("请登录")
	ErrSendSmsFailed           = errors.New("发送短信失败")
	ErrPhoneEmpty              = errors.New("手机号为空")
)

// 构建map使响应码对应响应信息
var codeMsg = map[ResponseCode]string{
	CodeSuccess:        "成功",
	CodeLoginSuccess:   "登录成功",
	CodeSignupSuccess:  "注册成功",
	CodeSendSmsSuccess: "发送短信成功",

	CodeUserExists:              "用户已存在",
	CodeUserNotExists:           "用户不存在",
	CodeUsernameStartWithNumber: "用户名不能以数字字符开头",
	CodeUsernameMustLess:        "用户名长度必须小于20",
	CodeInvalidPassword:         "用户名或密码错误",
	CodeInvalidParam:            "请求参数错误",
	CodeGetCaptchaFailed:        "获取验证码失败",
	CodeInvalidCaptcha:          "验证码错误或失效",
	CodeServiceBusy:             "服务繁忙",
	CodeNotLogin:                "请登录",
	CodeSendSmsFailed:           "发送短信失败",
	CodePhoneEmpty:              "手机号为空",
}

var errMsg = map[error]ResponseCode{
	ErrUserNotExists:           CodeUserNotExists,
	ErrServiceBusy:             CodeServiceBusy,
	ErrUsernameStartWithNumber: CodeUsernameStartWithNumber,
	ErrUsernameMustLess:        CodeUsernameMustLess,
	ErrGetCaptchaFailed:        CodeGetCaptchaFailed,
	ErrInvalidCaptcha:          CodeInvalidCaptcha,
	ErrNotLogin:                CodeNotLogin,
	ErrSendSmsFailed:           CodeSendSmsFailed,
	ErrPhoneEmpty:              CodePhoneEmpty,
}

// 通过响应码获取对应响应信息
func (r ResponseCode) getMsg() string {
	msg, ok := codeMsg[r]
	if !ok {
		msg = codeMsg[CodeServiceBusy]
	}
	return msg
}

// 通过错误获取响应码
func getCode(err error) ResponseCode {
	code, ok := errMsg[err]
	if !ok {
		code = errMsg[ErrServiceBusy]
	}
	return code
}
