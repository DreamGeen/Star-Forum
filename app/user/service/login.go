package service

import (
	"context"
	"log"
	"regexp"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/proto/user/userPb"
	"star/utils"
	"strings"
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
	// 获取用户密码
	cacheKey := "user:" + getUserIdentifier(user)
	checkedPassword, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		if err := queryFunc(user); err != nil {
			return "", err
		}
		return user.Password, nil
	})
	if err != nil {
		log.Println("query user error", err)
		return err
	}
	if err := utils.EqualsPassword(password, checkedPassword); err != nil {
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

	if ok := utils.ValidateCaptcha(ctx, req.Phone, req.Captcha); !ok {
		return str.ErrInvalidCaptcha
	}
	// 查询用户是否存在
	user := createUser(0, "", "", req.Phone, "")
	cacheKey := "user:" + user.Phone
	_, err = cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		if err := mysql.QueryUserByPhone(user); err != nil {
			return "", err
		}
		return user.Password, nil
	})
	if err != nil {
		log.Println("query user error", err)
		return err
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
