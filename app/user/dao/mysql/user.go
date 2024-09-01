package mysql

import (
	"database/sql"
	"errors"
	"log"
	"star/constant/str"
	"star/models"
)

const (
	queryUserByPhoneSQL    = "SELECT userid,username,phone,password,deletedAt FROM userLogin WHERE phone=?"
	queryUserByUsernameSQL = "SELECT userid, username,password,deletedAt FROM userLogin WHERE username=?"
	queryUserByEmailSQL    = "SELECT userid,username, email, password,deletedAt FROM userLogin WHERE email=?"
	insertUserLoginSQL     = "INSERT INTO userLogin(userId, username, password, phone) VALUES (?, ?, ?, ?)"
	insertUserSQL          = "INSERT INTO user(userId, username,img,sign) VALUES (?, ?, ?,?)"
)

// QueryUserByPhone 通过手机号查询用户密码
func QueryUserByPhone(u *models.User) error {
	return queryUser(u, queryUserByPhoneSQL, u.Phone)
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
	err := db.Get(u, sqlStr, args...)
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
	transaction, err := db.Beginx()
	if err != nil {
		log.Println("开启事务失败", err)
		return err
	}
	_, err = db.Exec(insertUserLoginSQL, u.UserId, u.Username, u.Password, u.Phone)
	if err != nil {
		if err := transaction.Rollback(); err != nil {
			log.Println("回滚事务失败", err)
		}
		log.Println("插入数据失败", err)
		return err
	}
	_, err = db.Exec(insertUserSQL, u.UserId, u.Username, u.Img, u.Signature)
	if err != nil {
		if err := transaction.Rollback(); err != nil {
			log.Println("回滚事务失败", err)
		}
		log.Println("插入数据失败", err)
		return err
	}
	if err := transaction.Commit(); err != nil {
		if err := transaction.Rollback(); err != nil {
			log.Println("回滚事务失败", err)
			return err
		}
		log.Println("插入数据失败", err)
		return err
	}
	return nil
}
