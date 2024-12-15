package models

// SignupUser 校验用户注册结构体
type SignupUser struct {
	Username     string `json:"username"  binding:"required,excludes=@"`
	Password     string `json:"registerPassword"  binding:"required"`
	RePasswd     string `json:"reRegisterPassword" binding:"required,eqfield=Password"`
	Phone        string `json:"phone" binding:"required"`
	CheckCodeKey string `json:"checkCodeKey" binding:"required"`
	CheckCode    string `json:"checkCode" bind:"required"`
	Captcha      string `json:"captcha" binding:"required"`
}

// LoginPassword  校验用户密码登录结构体
type LoginPassword struct {
	User         string `json:"user" binding:"required"`
	Password     string `json:"password" binding:"required" `
	CheckCodeKey string `json:"checkCodeKey" binding:"required"`
	CheckCode    string `json:"checkCode" bind:"required"`
}

// LoginCaptcha  校验用户验证码登录结构体
type LoginCaptcha struct {
	Phone   string `json:"phone" binding:"required"`
	Captcha string `json:"captcha" binding:"required"`
}
type Token struct {
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}
