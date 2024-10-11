package mysql

import (
	"database/sql"
	"errors"
	"star/app/constant/str"
	"star/app/models"
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
	tx, err := Client.Beginx()
	if err != nil {
		return str.ErrMessageError
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			err = str.ErrMessageError
		} else if err != nil {
			tx.Rollback()
		}
	}()
	_, err = Client.Exec(insertUserLoginSQL, u.UserId, u.Username, u.Password, u.Phone)
	if err != nil {
		return err
	}
	_, err = Client.Exec(insertUserSQL, u.UserId, u.Username, u.Img, u.Signature)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
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
