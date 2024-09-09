package models

// User 用户结构体
type User struct {
	UserId     int64   `db:"userId"`    //用户id
	Username   string  `db:"username"`  //用户名
	Password   string  `db:"password"`  //密码
	Phone      string  `db:"phone"`     //手机号
	Email      string  `db:"email"`     //邮箱
	Gender     string  `db:"gender"`    //性别
	Signature  string  `db:"sign"`      //签名
	Img        string  `db:"img"`       //头像
	Birth      string  `db:"birth"`     //生日
	Grade      uint32  `db:"grade"`     //等级
	Exp        int64   `db:"exp"`       //经验值
	DeleteTime *string `db:"deletedAt"` //删除时间
}

type LoginCheck struct {
	UserId   int64  `json:"userid"`
	Password string `json:"password"`
}

func (u *User) GetID() int64 {
	return u.UserId
}

func (u *User) IsDirty() bool {
	return u.Username != ""
}
