package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"regexp"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
	"star/utils"
	"strings"
	"sync"
)

type UserSrv struct {
}

var relationService relationPb.RelationService

// 正则表达式用于匹配手机号
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func New() {
	relationMicroService := micro.NewService(micro.Name(str.RelationServiceClient))
	relationService = relationPb.NewRelationService(str.RelationService, relationMicroService.Client())
}

// GetUserInfo 获取用户具体信息
func (u *UserSrv) GetUserInfo(ctx context.Context, req *userPb.GetUserInfoRequest, resp *userPb.GetUserInfoResponse) error {

	key := fmt.Sprintf("GetUserInfo:%d", req.UserId)
	user := new(models.User)
	found, err := cached.ScanGetUser(ctx, key, user)
	if err != nil {
		utils.Logger.Error("GetUserInfo failed", zap.Error(err))
		return err
	}
	if !found {
		utils.Logger.Info("GetUserInfo err:user not found", zap.Int64("userId", req.UserId))
		return str.ErrUserNotExists
	}
	resp.User = &userPb.User{
		UserId:   user.UserId,
		Exp:      user.Exp,
		Grade:    user.Grade,
		Gender:   user.Gender,
		UserName: user.Username,
		Img:      user.Img,
		Sign:     user.Signature,
		Birth:    user.Birth,
		IsFollow: false,
	}
	var wg sync.WaitGroup
	var isErr bool
	wg.Add(1)
	go func() {
		defer wg.Done()
		isFollowResp, err := relationService.IsFollow(ctx, &relationPb.IsFollowRequest{
			UserId:   req.UserId,
			FollowId: req.ActorId,
		})
		if err != nil {
			utils.Logger.Error("get is follow failed", zap.Error(err), zap.Int64("userId", req.UserId), zap.Any("followId", req.ActorId))
			isErr = true
			return
		}
		resp.User.IsFollow = isFollowResp.Result
	}()
	wg.Wait()
	if isErr {
		return str.ErrUserError
	}
	//返回user信息
	return nil
}

// LoginPassword 用密码的方式登录
func (u *UserSrv) LoginPassword(ctx context.Context, req *userPb.LSRequest, resp *userPb.LoginResponse) (err error) {
	user := createUser(0, str.Empty, req.Password, str.Empty, str.Empty)
	err = determineLoginMethod(ctx, req.User, user)
	if err != nil {
		return
	}
	accessToken, refreshToken, err := utils.GetToken(user)
	if err != nil {
		utils.Logger.Error("get  token  error",
			zap.Error(err),
			zap.String("user", req.User),
			zap.String("password", req.Password))
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
	// 获取用户密码和用户Id
	cacheKey := "user:" + getUserIdentifier(user)
	checkJson, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		if err := queryFunc(user); err != nil {
			return "", err
		}
		check := &models.LoginCheck{
			UserId:   user.UserId,
			Password: user.Password,
		}
		checkJson, err := json.Marshal(check)
		if err != nil {
			utils.Logger.Error("validatePassword error,json marshal checkJson error",
				zap.Error(err))
			return "", str.ErrLoginError
		}
		return string(checkJson), nil
	})
	if err != nil {
		utils.Logger.Error("validatePassword error,query user error",
			zap.Error(err))
		return err
	}
	check := new(models.LoginCheck)
	if err := json.Unmarshal([]byte(checkJson), &check); err != nil {
		utils.Logger.Error("json unmarshal checkJson error",
			zap.Error(err),
			zap.String("checkJson", checkJson),
			zap.Int64("user", user.UserId))
		return str.ErrLoginError
	}
	if err := utils.EqualsPassword(password, check.Password); err != nil {
		utils.Logger.Error("password error ,err:",
			zap.Error(err),
			zap.String("username", user.Username),
			zap.Int64("userId", user.UserId))
		return str.ErrInvalidPassword
	}
	//将获取的userId赋值给user
	user.UserId = check.UserId
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

	if ok := validateCaptcha(ctx, req.Phone, req.Captcha); !ok {
		return str.ErrInvalidCaptcha
	}
	// 查询用户是否存在
	user := createUser(0, "", "", req.Phone, "")
	cacheKey := "user:" + user.Phone
	checkJson, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		if err := mysql.QueryUserByPhone(user); err != nil {
			return "", err
		}
		check := &models.LoginCheck{
			UserId:   user.UserId,
			Password: user.Password,
		}
		checkJson, err := json.Marshal(check)
		if err != nil {
			utils.Logger.Error("loginCaptcha  service error ,json marshal checkJson error",
				zap.Error(err),
				zap.String("phone", req.Phone),
				zap.Any("check", check))
			return "", str.ErrLoginError
		}
		return string(checkJson), nil
	})
	if err != nil {
		utils.Logger.Error("login captcha query user error",
			zap.Error(err),
			zap.String("phone", req.Phone))
		return err
	}
	check := new(models.LoginCheck)
	if err := json.Unmarshal([]byte(checkJson), &check); err != nil {
		utils.Logger.Error("json unmarshal checkJson error",
			zap.Error(err),
			zap.String("phone", req.Phone),
			zap.String("checkJson", checkJson))
		return str.ErrLoginError
	}
	user.UserId = check.UserId
	accessToken, refreshToken, err := utils.GetToken(user)
	if err != nil {
		utils.Logger.Error("get token error",
			zap.Error(err),
			zap.String("phone", req.Phone))
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

// 校验验证码是否正确
func validateCaptcha(ctx context.Context, phone, captcha string) bool {
	cacheKey := "captcha:" + phone
	storedCaptcha, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get dao saved captcha fail",
			zap.Error(err),
			zap.String("phone", phone))
		return false
	}
	if !ok {
		utils.Logger.Error("captcha not exist",
			zap.String("phone", phone))
		return false
	}
	if storedCaptcha != captcha {
		utils.Logger.Error("captcha is wrong",
			zap.String("phone", phone))
		return false
	}
	return true
}

func (u *UserSrv) Signup(ctx context.Context, req *userPb.LSRequest, resp *userPb.EmptyLSResponse) (err error) {

	//校验验证码
	if ok := validateCaptcha(ctx, req.Phone, req.Captcha); !ok {
		return str.ErrInvalidCaptcha
	}
	//检查用户名和手机号是否已经注册过
	user := createUser(0, req.User, req.Password, req.Phone, "")
	if err = mysql.QueryUserByUsername(user); err == nil || !errors.Is(err, str.ErrUserNotExists) {
		utils.Logger.Warn("username is registered",
			zap.Error(err),
			zap.String("username", user.Username),
			zap.String("phone", user.Phone))
		return str.ErrUsernameExists
	}
	if err = mysql.QueryUserByPhone(user); err == nil || !errors.Is(err, str.ErrUserNotExists) {
		utils.Logger.Warn("phone is registered",
			zap.Error(err),
			zap.String("phone", user.Phone))
		return str.ErrPhoneRegistered
	}
	//对密码进行加密
	if user.Password, err = utils.EncryptPassword(user.Password); err != nil {
		utils.Logger.Error("encrypt password error",
			zap.Error(err),
			zap.String("password", user.Password),
			zap.String("phone", user.Phone))
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
		utils.Logger.Error("sign up error",
			zap.Error(err),
			zap.String("phone", req.Phone))
		return err
	}
	return
}
