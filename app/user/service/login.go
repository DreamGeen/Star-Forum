package service

import (
	"context"
	"regexp"
	logger "star/app/user/logger"
	"strings"
	"time"

	"go.uber.org/zap"

	"star/app/user/dao/mysql"
	"star/app/user/dao/redis"
	"star/models"
	"star/proto/user/userPb"
	"star/utils"
)

// 正则表达式用于匹配手机号
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// LoginPassword 用密码的方式登录
func (u *UserSrv) LoginPassword(ctx context.Context, req *userPb.LSRequest, resp *userPb.LoginResponse) error {
	user := createUser(0, "", req.Password, "", "")
	err := determineLoginMethod(ctx, req.User, user)
	if err != nil {
		return err
	}
	accessToken, refreshToken, err := utils.GetToken(user)
	if err != nil {
		logger.UserLogger.Error("获取JWT令牌失败", zap.Error(err))
		return err
	}
	resp.Token = &userPb.LoginResponse_Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return nil
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
	if err := redis.Get(ctx, cacheKey, user); err == nil {
		// 从缓存中获取到了用户信息
	} else {
		// 缓存中没有，从数据库中获取
		if err := queryFunc(user); err != nil {
			return err
		}
		// 将用户信息缓存到Redis中并设置超时时间
		redis.Set(ctx, cacheKey, user, 24*time.Hour)
	}
	if !utils.EqualsPassword(password, user.Password) {
		logger.UserLogger.Error("密码错误", zap.String("userName", user.Username), zap.String("userPassword", user.Password))
		return utils.ErrUserNotExists
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
func (u *UserSrv) LoginCaptcha(ctx context.Context, req *userPb.LSRequest, resp *userPb.LoginResponse) error {
	if err := utils.ValidateCaptcha(req.Phone, req.Captcha); err != nil {
		return err
	}
	// 先尝试从缓存中获取用户信息
	user := createUser(0, "", "", req.Phone, "")
	cacheKey := "user:" + user.Phone
	if err := redis.Get(ctx, cacheKey, user); err == nil {
		// 从缓存中获取到了用户信息
	} else {
		// 缓存中没有，从数据库中获取
		if err := mysql.QueryUserByPhone(user); err != nil {
			logger.UserLogger.Error("密码错误", zap.String("userPhone", user.Phone), zap.Error(err))
			return err
		}
		// 将用户信息缓存到Redis中并设置超时时间
		redis.Set(ctx, cacheKey, user, 24*time.Hour)
	}
	accessToken, refreshToken, err := utils.GetToken(user)
	if err != nil {
		logger.UserLogger.Error("获取JWT令牌失败", zap.Error(err))
		return err
	}

	resp.Token = &userPb.LoginResponse_Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return nil
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
