package service

import (
	"context"
	"errors"
	"log"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/proto/user/userPb"
	"star/utils"
)

func (u *UserSrv) Signup(ctx context.Context, req *userPb.LSRequest, resp *userPb.EmptyLSResponse) (err error) {

	//校验验证码
	if ok := validateCaptcha(ctx, req.Phone, req.Captcha); !ok {
		return str.ErrInvalidCaptcha
	}
	//检查用户名和手机号是否已经注册过
	user := createUser(0, req.User, req.Password, req.Phone, "")
	if err = mysql.QueryUserByUsername(user); err == nil || !errors.Is(err, str.ErrUserNotExists) {
		log.Println("用户名已注册", user.Username, err)
		return str.ErrUsernameExists
	}
	if err = mysql.QueryUserByPhone(user); err == nil || !errors.Is(err, str.ErrUserNotExists) {
		log.Println("手机号已注册", user.Phone, err)
		return str.ErrPhoneRegistered
	}
	//对密码进行加密
	if user.Password, err = utils.EncryptPassword(user.Password); err != nil {
		log.Println("密码加密失败", err)
		return str.ErrSignupError
	}
	//生成用户id
	user.UserId = utils.GetID()
	//用户默认签名
	user.Signature = str.DefaultSignature
	//用户默认头像
	user.Img = str.DefaultImg
	//将用户插入mysql当中
	if err := mysql.InsertUser(user); err != nil {
		return err
	}
	log.Println("用户注册成功", user.UserId, user.Username)
	return
}
