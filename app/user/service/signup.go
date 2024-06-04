package service

import (
	"context"
	"errors"
	"log"

	"star/app/user/dao/mysql"
	"star/proto/user/userPb"
	"star/utils"
)

// 默认签名
var defaultSignature = "签名是一种态度，我想我可以很酷"

// 默认头像
var defaultImg = ""

func (u *UserSrv) Signup(ctx context.Context, req *userPb.LSRequest, resp *userPb.EmptyLSResponse) (err error) {
	//校验验证码
	if err = utils.ValidateCaptcha(req.Phone, req.Captcha); err != nil {
		return
	}
	//检查用户名和手机号是否已经注册过
	user := createUser(0, req.User, req.Password, req.Phone, "")
	if err = mysql.QueryUserByUsername(user); !errors.Is(err, utils.ErrUserNotExists) {
		log.Println("用户名已注册", user.Username, err)
		return
	}
	if err = mysql.QueryUserByPhone(user); !errors.Is(err, utils.ErrUserNotExists) {
		log.Println("手机号已注册", user.Phone, err)
		return
	}
	//对密码进行加密
	if user.Password, err = utils.EncryptPassword(user.Password); err != nil {
		log.Println("密码加密失败", err)
		return err
	}
	//生成用户id
	user.UserId = utils.GetID()
	//用户默认签名
	user.Signature = defaultSignature
	//用户默认头像
	user.Img = defaultImg

	//最终
	//先把数据存入redis返回响应
	//后用消息队列异步将数据存入mysql

	//暂定
	//将用户信息储存到mysql
	if err = mysql.InsertUser(user); err != nil {
		log.Println("插入用户失败", err)
		return
	}
	log.Println("用户注册成功", user.UserId, user.Username)
	return
}
