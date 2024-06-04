package models

// SignupUser 校验用户注册结构体
type SignupUser struct {
	Username string `json:"username"  binding:"required,excludes=@"`
	Password string `json:"password"  binding:"required"`
	RePasswd string `json:"repassword" binding:"required,eqfield=Password"`
	Phone    string `json:"phone" binding:"required"`
	Captcha  string `json:"captcha" binding:"required"`
}

// LoginPassword  校验用户密码登录结构体
type LoginPassword struct {
	User     string `json:"usr" binding:"required"`
	Password string `json:"password" binding:"required" `
}

// LoginCaptcha  校验用户验证码登录结构体
type LoginCaptcha struct {
	Phone   string `json:"phone" binding:"required"`
	Captcha string `json:"captcha" binding:"required"`
}

// User 用户结构体
type User struct {
	UserId     int64   `db:"userid"`    //用户id
	Username   string  `db:"username"`  //用户名
	Password   string  `db:"password"`  //密码
	Phone      string  `db:"phone"`     //手机号
	Email      string  `db:"email"`     //邮箱
	Gender     string  `db:"gender"`    //性别
	Signature  string  `db:"sign"`      //签名
	Img        string  `db:"img"`       //头像
	Birth      string  `db:"birth"`     //生日
	Grade      uint8   `db:"grade"`     //等级
	Exp        int     `db:"exp"`       //经验值
	DeleteTime *string `db:"deletedAt"` //删除时间
}
