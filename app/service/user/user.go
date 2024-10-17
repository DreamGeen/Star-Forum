package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"regexp"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/utils/jwt"
	"star/app/utils/logging"
	"star/app/utils/password"
	"star/app/utils/snowflake"
	"star/proto/collect/collectPb"
	"star/proto/like/likePb"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
	"strings"
	"sync"
)

type UserSrv struct {
}

var relationService relationPb.RelationService
var likeService likePb.LikeService
var collectService collectPb.CollectService
var userIns = new(UserSrv)

// 正则表达式用于匹配手机号
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func (u *UserSrv) New() {
	relationMicroService := micro.NewService(micro.Name(str.RelationServiceClient))
	relationService = relationPb.NewRelationService(str.RelationService, relationMicroService.Client())

	likeMicroService := micro.NewService(micro.Name(str.LikeServiceClient))
	likeService = likePb.NewLikeService(str.LikeService, likeMicroService.Client())
}

// GetUserInfo 获取用户具体信息
func (u *UserSrv) GetUserInfo(ctx context.Context, req *userPb.GetUserInfoRequest, resp *userPb.GetUserInfoResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserInfoService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "UserService.GetUserInfo")

	key := fmt.Sprintf("GetUserInfo:%d", req.UserId)
	user := new(models.User)
	found, err := cached.ScanGetUser(ctx, key, user)
	if err != nil {
		logger.Error("GetUserInfo failed",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return err
	}
	if !found {
		logger.Info("GetUserInfo err:user not found",
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrUserNotExists
	}
	resp.User = &userPb.User{
		UserId:   user.UserId,
		Exp:      user.Exp,
		Grade:    user.Grade,
		Gender:   &user.Gender,
		UserName: user.Username,
		Img:      &user.Img,
		Sign:     &user.Signature,
		Birth:    &user.Birth,
		IsFollow: false,
	}
	var wg sync.WaitGroup
	var isErr bool
	wg.Add(6)
	go func() {
		defer wg.Done()
		isFollowResp, err := relationService.IsFollow(ctx, &relationPb.IsFollowRequest{
			UserId:   req.UserId,
			FollowId: req.ActorId,
		})
		if err != nil {
			logger.Error("get is follow failed",
				zap.Error(err),
				zap.Int64("userId", req.UserId),
				zap.Any("followId", req.ActorId))
			isErr = true
			return
		}
		resp.User.IsFollow = isFollowResp.Result
	}()

	go func() {
		defer wg.Done()
		countFollowResp, err := relationService.CountFollow(ctx, &relationPb.CountFollowRequest{
			UserId: req.UserId,
		})
		if err != nil {
			logger.Error("get user follow count error",
				zap.Error(err),
				zap.Int64("user", req.UserId))
			isErr = true
			return
		}
		resp.User.FollowCount = &countFollowResp.Count
	}()
	go func() {
		defer wg.Done()
		countFansResp, err := relationService.CountFans(ctx, &relationPb.CountFansRequest{
			UserId: req.UserId,
		})
		if err != nil {
			logger.Error("get user fans count error",
				zap.Error(err),
				zap.Int64("user", req.UserId))
			isErr = true
			return
		}
		resp.User.FansCount = &countFansResp.Count
	}()
	go func() {
		defer wg.Done()
		getUserLikeCountResp, err := likeService.GetUserLikeCount(ctx, &likePb.GetUserLikeCountRequest{
			UserId: req.UserId,
		})
		if err != nil {
			logger.Error("get user like count error",
				zap.Error(err),
				zap.Int64("user", req.UserId))
			isErr = true
			return
		}
		resp.User.LikeCount = &getUserLikeCountResp.Count
	}()
	go func() {
		getUserTotalLikeResp, err := likeService.GetUserTotalLike(ctx, &likePb.GetUserTotalLikeRequest{
			UserId: req.UserId,
		})
		if err != nil {
			logger.Error("get user total liked count error",
				zap.Error(err),
				zap.Int64("user", req.UserId))
			isErr = true
			return
		}
		resp.User.TotalLiked = &getUserTotalLikeResp.Count
	}()
	go func() {
		getUserCollectCount, err := collectService.GetUserCollectCount(ctx, &collectPb.GetUserCollectCountRequest{
			UserId: req.UserId,
		})
		if err != nil {
			logger.Error("get user  collect count error",
				zap.Error(err),
				zap.Int64("user", req.UserId))
			isErr = true
			return
		}
		resp.User.CollectCount = &getUserCollectCount.Count
	}()
	wg.Wait()
	if isErr {
		return str.ErrUserError
	}
	//返回user信息
	return nil
}

// GetUserExistInformation 检查用户是否存在
func (u *UserSrv) GetUserExistInformation(ctx context.Context, req *userPb.GetUserExistInformationRequest, resp *userPb.GetUserExistInformationResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserExisted")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "UserService.GetUserExisted")

	user := new(models.User)
	key := fmt.Sprintf("GetUserInfo:%d", req.UserId)
	found, err := cached.ScanGetUser(ctx, key, user)
	if err != nil {
		logger.Error("GetUserInfo failed",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrUserError
	}
	if !found {
		logger.Info("GetUserInfo err:user not found",
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)

		resp.Existed = false
		return nil
	}
	resp.Existed = true
	return nil
}

// LoginPassword 用密码的方式登录
func (u *UserSrv) LoginPassword(ctx context.Context, req *userPb.LSRequest, resp *userPb.LoginResponse) (err error) {
	ctx, span := tracing.Tracer.Start(ctx, "LoginPasswordService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "userService.LoginPassword")

	user := createUser(0, str.Empty, req.Password, str.Empty, str.Empty)
	err = determineLoginMethod(ctx, span, logger, req.User, user)
	if err != nil {
		return
	}
	accessToken, refreshToken, err := jwt.GetToken(user)
	if err != nil {
		logger.Error("get  token  error",
			zap.Error(err),
			zap.String("user", req.User),
			zap.String("password", req.Password))
		logging.SetSpanError(span, err)
		return str.ErrLoginError
	}

	resp.Token = &userPb.LoginResponse_Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}

// determineLoginMethod 确定用户登录方法是通过手机、邮箱还是用户名
func determineLoginMethod(ctx context.Context, span trace.Span, logger *zap.Logger, userInput string, user *models.User) error {
	if isPhoneNumber(userInput) {
		user.Phone = userInput
		return loginByPhone(ctx, span, logger, user)
	} else if isEmail(userInput) {
		user.Email = userInput
		return loginByEmail(ctx, span, logger, user)
	} else {
		user.Username = userInput
		return loginByUsername(ctx, span, logger, user)
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
func loginByPhone(ctx context.Context, span trace.Span, logger *zap.Logger, u *models.User) error {
	return validatePassword(ctx, span, logger, u, mysql.QueryUserByPhone)
}

// 通过用户名登录
func loginByUsername(ctx context.Context, span trace.Span, logger *zap.Logger, u *models.User) error {
	return validatePassword(ctx, span, logger, u, mysql.QueryUserByUsername)
}

// 通过邮箱登录
func loginByEmail(ctx context.Context, span trace.Span, logger *zap.Logger, u *models.User) error {
	return validatePassword(ctx, span, logger, u, mysql.QueryUserByEmail)
}

// 验证密码，检查用户是否存在并验证密码
func validatePassword(ctx context.Context, span trace.Span, logger *zap.Logger, user *models.User, queryFunc func(*models.User) error) error {
	truePassword := user.Password
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
			logger.Error("validatePassword error,json marshal checkJson error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return "", str.ErrLoginError
		}
		return string(checkJson), nil
	})
	if err != nil {
		logger.Error("validatePassword error,query user error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return err
	}
	check := new(models.LoginCheck)
	if err := json.Unmarshal([]byte(checkJson), &check); err != nil {
		logger.Error("json unmarshal checkJson error",
			zap.Error(err),
			zap.String("checkJson", checkJson),
			zap.Int64("user", user.UserId))
		logging.SetSpanError(span, err)
		return str.ErrLoginError
	}
	if err := password.Equals(truePassword, check.Password); err != nil {
		logger.Error("password error ,err:",
			zap.Error(err),
			zap.String("username", user.Username),
			zap.Int64("userId", user.UserId))
		logging.SetSpanError(span, err)
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
	ctx, span := tracing.Tracer.Start(ctx, "LoginCaptchaService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "UserService.LoginCaptcha")

	if ok := validateCaptcha(ctx, span, logger, req.Phone, req.Captcha); !ok {
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
			logger.Error("json marshal checkJson error",
				zap.Error(err),
				zap.String("phone", req.Phone),
				zap.Any("check", check))
			logging.SetSpanError(span, err)
			return "", str.ErrLoginError
		}
		return string(checkJson), nil
	})
	if err != nil {
		logger.Error("login captcha query user error",
			zap.Error(err),
			zap.String("phone", req.Phone))
		logging.SetSpanError(span, err)
		return err
	}
	check := new(models.LoginCheck)
	if err := json.Unmarshal([]byte(checkJson), &check); err != nil {
		logger.Error("json unmarshal checkJson error",
			zap.Error(err),
			zap.String("phone", req.Phone),
			zap.String("checkJson", checkJson))
		logging.SetSpanError(span, err)
		return str.ErrLoginError
	}
	user.UserId = check.UserId
	accessToken, refreshToken, err := jwt.GetToken(user)
	if err != nil {
		logger.Error("get token error",
			zap.Error(err),
			zap.String("phone", req.Phone))
		logging.SetSpanError(span, err)
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
func validateCaptcha(ctx context.Context, span trace.Span, logger *zap.Logger, phone, captcha string) bool {
	cacheKey := "captcha:" + phone
	storedCaptcha, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("get dao saved captcha fail",
			zap.Error(err),
			zap.String("phone", phone))
		logging.SetSpanError(span, err)
		return false
	}
	if !ok {
		logger.Error("captcha not exist",
			zap.String("phone", phone))
		logging.SetSpanError(span, err)
		return false
	}
	if storedCaptcha != captcha {
		logger.Error("captcha is wrong",
			zap.String("phone", phone))
		logging.SetSpanError(span, err)
		return false
	}
	return true
}

func (u *UserSrv) Signup(ctx context.Context, req *userPb.LSRequest, resp *userPb.EmptyLSResponse) (err error) {
	ctx, span := tracing.Tracer.Start(ctx, "SignupService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "UserService.Signup")

	//校验验证码
	if ok := validateCaptcha(ctx, span, logger, req.Phone, req.Captcha); !ok {
		return str.ErrInvalidCaptcha
	}
	//检查用户名和手机号是否已经注册过
	user := createUser(0, req.User, req.Password, req.Phone, "")
	if err = mysql.QueryUserByUsername(user); err == nil || !errors.Is(err, str.ErrUserNotExists) {
		logger.Warn("username is registered",
			zap.Error(err),
			zap.String("username", user.Username),
			zap.String("phone", user.Phone))
		logging.SetSpanError(span, err)
		return str.ErrUsernameExists
	}
	if err = mysql.QueryUserByPhone(user); err == nil || !errors.Is(err, str.ErrUserNotExists) {
		logger.Warn("phone is registered",
			zap.Error(err),
			zap.String("phone", user.Phone))
		logging.SetSpanError(span, err)
		return str.ErrPhoneRegistered
	}
	//对密码进行加密
	if user.Password, err = password.Encrypt(user.Password); err != nil {
		logger.Error("encrypt password error",
			zap.Error(err),
			zap.String("password", user.Password),
			zap.String("phone", user.Phone))
		logging.SetSpanError(span, err)
		return str.ErrSignupError
	}
	//生成用户id
	user.UserId = snowflake.GetID()
	//用户默认签名
	user.Signature = str.DefaultSignature
	//用户默认头像
	user.Img = str.DefaultImg
	//将用户插入mysql当中
	if err := mysql.InsertUser(user); err != nil {
		logger.Error("sign up error",
			zap.Error(err),
			zap.String("phone", req.Phone))
		logging.SetSpanError(span, err)
		return err
	}
	return
}
