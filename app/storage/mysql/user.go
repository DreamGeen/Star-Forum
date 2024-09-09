package mysql

import (
	"database/sql"
	"errors"
	"log"
	"star/constant/str"
	"star/models"
)

const (
	queryUserByPhoneSQL    = "SELECT userId,username,phone,password,deletedAt FROM userLogin WHERE phone=?"
	queryUserByUsernameSQL = "SELECT userId, username,password,deletedAt FROM userLogin WHERE username=?"
	queryUserByEmailSQL    = "SELECT userId,username, email, password,deletedAt FROM userLogin WHERE email=?"
	insertUserLoginSQL     = "INSERT INTO userLogin(userId, username, password, phone) VALUES (?, ?, ?, ?)"
	insertUserSQL          = "INSERT INTO user(userId, username,img,sign) VALUES (?, ?, ?,?)"
	queryUserInfoSQL       = "select  userId,userName,gender,sign,birth,grade,exp,img,updatedAt,createdAt,deletedAt from user  where userId=?"
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
		log.Println("查询用户信息失败", err)
		if errors.Is(err, sql.ErrNoRows) {
			return str.ErrUserNotExists
		}
		return err
	}
	if u.DeleteTime != nil {
		return str.ErrUserNotExists
	}
	return nil
}

// InsertUser 将用户信息插入mysql
func InsertUser(u *models.User) error {
	//将用户信息插入mysql
	transaction, err := Client.Beginx()
	if err != nil {
		log.Println("开启事务失败", err)
		return str.ErrSignupError
	}
	_, err = Client.Exec(insertUserLoginSQL, u.UserId, u.Username, u.Password, u.Phone)
	if err != nil {
		if err := transaction.Rollback(); err != nil {
			log.Println("回滚事务失败", err)
		}
		log.Println("插入数据失败", err)
		return str.ErrSignupError
	}
	_, err = Client.Exec(insertUserSQL, u.UserId, u.Username, u.Img, u.Signature)
	if err != nil {
		if err := transaction.Rollback(); err != nil {
			log.Println("回滚事务失败", err)
		}
		log.Println("插入数据失败", err)
		return str.ErrSignupError
	}
	if err := transaction.Commit(); err != nil {
		if err := transaction.Rollback(); err != nil {
			log.Println("回滚事务失败", err)
			return str.ErrSignupError
		}
		log.Println("插入数据失败", err)
		return str.ErrSignupError
	}
	return nil
}

func QueryUserInfo(user *models.User, userId int64) error {
	if err := Client.Get(user, queryUserInfoSQL, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.ErrUserNotExists
		}
		return str.ErrUserError
	}
	return nil
}
