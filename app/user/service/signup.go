package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"star/app/models"
	"star/app/user/dao/mysql"
	"star/app/user/dao/redis"
	"star/constant/str"
	"star/mq"
	"star/proto/user/userPb"
	"star/utils"
	"time"
)

func (u *UserSrv) Signup(ctx context.Context, req *userPb.LSRequest, resp *userPb.EmptyLSResponse) (err error) {

	//校验验证码
	if ok := utils.ValidateCaptcha(req.Phone, req.Captcha); !ok {
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
	//先把数据存入redis
	redisKey := fmt.Sprintf("user_signup:%d", user.UserId)
	if err = redis.SetUser(ctx, redisKey, user, 10*time.Minute); err != nil {
		log.Println("将注册信息储存进redis失败", err)
		return str.ErrSignupError
	}
	//向消息队列发送redisKey
	if err = mq.SendMessage("user_signup", []byte(redisKey)); err != nil {
		log.Println("发送消息失败", err)
		return str.ErrSignupError
	}
	log.Println("用户注册成功", user.UserId, user.Username)
	return
}

func insert() {
	msgs, err := mq.ConsumeMessage("user_signup")
	if err != nil {
		log.Println("获取消息失败", err)
	}
	for msg := range msgs {
		var user *models.User
		err := redis.GetUser(context.Background(), string(msg.Body), user)
		if err != nil {
			log.Println("获取注册用户失败", err)
			continue
		}
		if err = mysql.InsertUser(user); err != nil {
			log.Println("插入用户失败", err)
			return
		}
	}
}
