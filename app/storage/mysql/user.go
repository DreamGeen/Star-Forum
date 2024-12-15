package mysql

import (
	"database/sql"
	"errors"
	"star/app/constant/str"
	"star/app/models"
	"time"
)

const (
	queryUserByPhoneSQL     = "SELECT user_id,username,phone,password FROM user_info WHERE phone=?"
	queryUserByUsernameSQL  = "SELECT user_id, username,password FROM user_info WHERE username=?"
	queryUserByEmailSQL     = "SELECT user_id,username, email, password FROM user_info WHERE email=?"
	insertUserSQL           = "INSERT INTO user_info(user_id, username,password,phone,avatar,person_introduction,sex,join_time,total_coin_count,current_coin_count) VALUES (?,?, ?,?, ?,?,?,?,?,?)"
	queryUserInfoSQL        = "select  user_id,username,phone,email,person_introduction,avatar,birthday,school,notice_info,last_login_ip,total_coin_count,current_coin_count,theme,sex,status,last_login_time,join_time from user_info  where user_id=?"
	updateLoginTimeAndIpSQL = "update user_info set last_login_time=?,last_login_ip=? where user_id=? "
)

// QueryUserByPhone 通过手机号查询用户密码
func QueryUserByPhone(u *models.User) error {
	return queryUser(u, queryUserByPhoneSQL, u.Username)
}

// QueryUserByUsername 通过用户名查询用户密码
func QueryUserByUsername(u *models.User) error {
	return queryUser(u, queryUserByUsernameSQL, u.Username)
}

// QueryUserByEmail 通过邮箱查询用户密码
func QueryUserByEmail(u *models.User) error {
	return queryUser(u, queryUserByEmailSQL, u.Email)
}

func queryUser(u *models.User, sqlStr string, args ...interface{}) error {
	err := Client.Get(u, sqlStr, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.ErrUserNotExists
		}
		return err
	}
	return nil
}

// InsertUser 将用户信息插入mysql
func InsertUser(u *models.User) error {
	//将用户信息插入mysql
	_, err := Client.Exec(insertUserSQL, u.UserId, u.Username, u.Password, u.Phone, u.Avatar, u.Introduction, u.Sex, u.JoinTime, u.TotalCoinCount, u.CurrentCoinCount)
	if err != nil {
		return err
	}
	return nil
}

func QueryUserInfo(user *models.User, userId int64) error {
	if err := Client.Get(user, queryUserInfoSQL, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.ErrUserNotExists
		}
		return err
	}
	return nil
}

func UpdateLoginTimeAndIp(loginTime time.Time, loginIp string, userId int64) error {
	if _, err := Client.Exec(updateLoginTimeAndIpSQL, loginTime, loginIp, userId); err != nil {
		return err
	}
	return nil
}
