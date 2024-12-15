package models

import "time"

// User 用户结构体
type User struct {
	UserId           int64      `db:"user_id" redis:"user_id"`                         //用户id
	Username         string     `db:"username" redis:"username"`                       //用户名
	Password         string     `db:"password" redis:"password"`                       //密码
	Phone            string     `db:"phone" redis:"phone"`                             //手机号
	Introduction     string     `db:"person_introduction" redis:"person_introduction"` //个人简介
	Avatar           string     `db:"avatar" redis:"avatar"`                           //头像
	School           string     `db:"school" redis:"school"`                           //学校
	LastLoginIp      string     `db:"last_login_ip" redis:"last_login_ip"`             //最后登录ip
	Email            string     `db:"email" redis:"email"`                             //邮箱
	Birthday         string     `db:"birthday" redis:"birthday"`                       //生日
	NoticeInfo       string     `db:"notice_info" redis:"notice_info"`                 //空间公告
	TotalCoinCount   uint32     `db:"total_coin_count" redis:"total_coin_count"`       //硬币总数量
	CurrentCoinCount uint32     `db:"current_coin_count" redis:"current_coin_count"`   //当前硬币数
	Sex              uint32     `db:"sex" redis:"sex"`                                 //性别  0女 1男 2未知
	Status           uint32     `db:"status" redis:"status"`                           //是否被禁用  0禁用 1正常
	Theme            uint32     `db:"theme" redis:"theme"`                             //主题
	LastLoginTime    *time.Time `db:"last_login_time" `                                //最后登录时间
	JoinTime         *time.Time `db:"join_time"  `                                     //加入时间
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
