package service

import (
	"context"
	"log"
	"regexp"
	"star/app/models"
	"star/app/user/dao/mysql"
	"star/app/user/dao/redis"
	"star/constant/str"
	"star/proto/user/userPb"
	"star/utils"
	"strings"
	"time"
)

// 正则表达式用于匹配手机号
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// LoginPassword 用密码的方式登录
func (u *UserSrv) LoginPassword(ctx context.Context, req *userPb.LSRequest, resp *userPb.LoginResponse) (err error) {
	user := createUser(0, str.Empty, req.Password, str.Empty, str.Empty)
	err = determineLoginMethod(ctx, req.User, user)
	if err != nil {
		return
	}
	accessToken, refreshToken, err := utils.GetToken(user)
	if err != nil {
		log.Println("获取JWT令牌失败", err)
		return str.ErrLoginError
	}

	resp.Token = &userPb.LoginResponse_Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}

// determineLoginMethod 确定用户登录方法是通过手机、邮箱还是用户名
func determineLoginMethod(ctx context.Context, userInput string, user *models.User) error {
	if isPhoneNumber(userInput) {
		user.Phone = userInput
		return loginByPhone(ctx, user)
	} else if isEmail(userInput) {
		user.Email = userInput
		return loginByEmail(ctx, user)
	} else {
		user.Username = userInput
		return loginByUsername(ctx, user)
	}
}

// 判断输入字符串是否为手机号
func isPhoneNumber(input string) bool {
	return phoneRegex.MatchString(input)
}

// 判断输入字符串是否为邮箱
func isEmail(input string) bool {
	return strings.Contains(input, "@")
}

// 通过手机号登录
func loginByPhone(ctx context.Context, u *models.User) error {
	return validatePassword(ctx, u, mysql.QueryUserByPhone)
}

// 通过用户名登录
func loginByUsername(ctx context.Context, u *models.User) error {
	return validatePassword(ctx, u, mysql.QueryUserByUsername)
}

// 通过邮箱登录
func loginByEmail(ctx context.Context, u *models.User) error {
	return validatePassword(ctx, u, mysql.QueryUserByEmail)
}

// 验证密码，检查用户是否存在并验证密码
func validatePassword(ctx context.Context, user *models.User, queryFunc func(*models.User) error) error {
	password := user.Password
	// 先尝试从缓存中获取用户信息
	cacheKey := "user:" + getUserIdentifier(user)
	if err := redis.GetUser(ctx, cacheKey, user); err == nil {
		// 从缓存中获取到了用户信息
	} else {
		// 缓存中没有，从数据库中获取
		if err := queryFunc(user); err != nil {
			log.Println("query user err:", err)
			return str.ErrUserNotExists
		}
		// 将用户信息缓存到Redis中并设置超时时间
		_ = redis.SetUser(ctx, cacheKey, user, 2*24*time.Hour)
	}
	if err := utils.EqualsPassword(password, user.Password); err != nil {
		log.Println("密码错误,err:", err, user.Username, user.UserId)
		return str.ErrInvalidPassword
	}
	return nil
}

// 获取用户唯一标识
func getUserIdentifier(u *models.User) string {
	if u.Phone != "" {
		return u.Phone
	} else if u.Email != "" {
		return u.Email
	}
	return u.Username
}

// LoginCaptcha 用验证码的方式登录
func (u *UserSrv) LoginCaptcha(ctx context.Context, req *userPb.LSRequest, resp *userPb.LoginResponse) (err error) {

	if ok := utils.ValidateCaptcha(req.Phone, req.Captcha); !ok {
		return str.ErrInvalidCaptcha
	}
	// 先尝试从缓存中获取用户信息
	user := createUser(0, "", "", req.Phone, "")
	cacheKey := "user:" + user.Phone
	if err = redis.GetUser(ctx, cacheKey, user); err == nil {
		// 从缓存中获取到了用户信息
	} else {
		// 缓存中没有，从数据库中获取
		if err = mysql.QueryUserByPhone(user); err != nil {
			log.Println("该手机号未注册", user.Phone)
			return str.ErrPhoneEmpty
		}
		// 将用户信息缓存到Redis中并设置超时时间
		_ = redis.SetUser(ctx, cacheKey, user, 2*24*time.Hour)
	}
	accessToken, refreshToken, err := utils.GetToken(user)
	if err != nil {
		log.Println("获取JWT令牌失败", err)
		return str.ErrLoginError
	}

	resp.Token = &userPb.LoginResponse_Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}

// createUser 创建一个新的用户对象
func createUser(userid int64, username, password, phone, email string) *models.User {
	return &models.User{
		UserId:   userid,
		Username: username,
		Password: password,
		Phone:    phone,
		Email:    email,
	}
}
