package str

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	UserNotExistsCode = 10001 + iota
	UsernameExistsCode
	UsernameMustLess
	UsernameStartWith
	InvalidAccessTokenCode
	InvalidPasswordCode
	InvalidParamCode
	InvalidImgCaptchaCode
	InvalidCaptchaCode
	NotLoginCode
	PhoneEmptyCode
	PhoneRegisteredCode
	PhoneUnregisteredCode
	CommentNotExistsCode
	HubExistsCode
	HubNotExistCode
	MessageInadequateCode
	DescriptionShortCode
	DescriptionLongCode
	CommunityNameEmptyCode
	CommunityNameLongCode
	CommunityNameExistsCode
	CommunityNotExitsCode
	PostNotExistsCode
	UploadCode
	RequestTooFrequentlyCode
	AutoLoginErrorCode
	CategoryNotExistsCode
	CategoryIdExistsCode
)

const (
	ServiceBusyCode = 50001 + iota
	UserErrorCode
	LoginErrorCode
	SignupErrorCode
	SendSmsErrorCode
	CommentErrorCode
	CommunityErrorCode
	MessageErrorCode
	RelationErrorCode
	LikeErrorCode
	CollectErrorCode
	PublishErrorCode
	FeedErrorCode
	CategoryErrorCode
)

var (
	ErrUserNotExists        = errors.New("用户不存在")
	ErrUsernameExists       = errors.New("用户名已存在")
	ErrUsernameMustLess     = errors.New("用户名长度必须小于20")
	ErrUsernameStartWith    = errors.New("用户名不能以数字开头")
	ErrInvalidAccessToken   = errors.New("非法的accessToken")
	ErrInvalidPassword      = errors.New("用户名或密码错误")
	ErrInvalidParam         = errors.New("请求参数错误")
	ErrInvalidImgCaptcha    = errors.New("图形验证码错误")
	ErrInvalidCaptcha       = errors.New("验证码错误")
	ErrNotLogin             = errors.New("请登录")
	ErrPhoneEmpty           = errors.New("手机号为空")
	ErrPhoneRegistered      = errors.New("手机号已注册")
	ErrPhoneUnregistered    = errors.New("该手机号未注册")
	ErrCommentNotExists     = errors.New("评论不存在或已被删除")
	ErrHubExists            = errors.New("聊天频道已存在")
	ErrHubNotExists         = errors.New("聊天频道不存在")
	ErrMessageInadequate    = errors.New("消息不足")
	ErrDescriptionShort     = errors.New("简介字数不能少于2")
	ErrDescriptionLong      = errors.New("简介字数不能大于50")
	ErrCommunityNameEmpty   = errors.New("社区名不能为空")
	ErrCommunityNameLong    = errors.New("社区名不能大于10")
	ErrCommunityNameExists  = errors.New("该社区名已存在")
	ErrCommunityNotExists   = errors.New("该社区不存在")
	ErrPostNotExists        = errors.New("该帖子不存在")
	ErrUpload               = errors.New("上传失败")
	ErrRequestTooFrequently = errors.New("请求频繁")
	ErrAutoLoginError       = errors.New("自动登录失败")
	ErrCategoryNotExists    = errors.New("分类不存在")
	ErrCategoryIdExists     = errors.New("分类编号已存在")
)

var (
	ErrServiceBusy    = errors.New("服务繁忙")
	ErrUserError      = errors.New("用户服务出现内部错误，请稍后再试！")
	ErrLoginError     = errors.New("登录服务出现内部错误，请稍后再试！")
	ErrSignupError    = errors.New("注册服务出现内部错误，请稍后再试！")
	ErrSendSmsError   = errors.New("发送短信失败，请稍后再试！")
	ErrCommentError   = errors.New("评论服务出现内部错误，请稍后再试！")
	ErrCommunityError = errors.New("社区服务内部出现错误，请稍后再试！")
	ErrFeedError      = errors.New("feed流服务内部出现问题，请稍后再试! ")
	ErrMessageError   = errors.New("消息服务内部出现问题，请稍后再试！")
	ErrRelationError  = errors.New("关系服务内部出现问题，请稍后再试！")
	ErrLikeError      = errors.New("点赞服务内部出现问题，请稍后再试！")
	ErrCollectError   = errors.New("收藏服务内部出现问题，请稍后再试！")
	ErrPublishError   = errors.New("发帖服务内部出现问题，请稍后再试！")
	ErrCategoryError  = errors.New("分类服务内部出现问题，请稍后再试! ")
)

var codeMap = map[error]int32{
	ErrUserNotExists:        UserNotExistsCode,
	ErrUsernameExists:       UsernameExistsCode,
	ErrUsernameMustLess:     UsernameMustLess,
	ErrUsernameStartWith:    UsernameStartWith,
	ErrInvalidAccessToken:   InvalidAccessTokenCode,
	ErrInvalidImgCaptcha:    InvalidImgCaptchaCode,
	ErrInvalidPassword:      InvalidPasswordCode,
	ErrInvalidParam:         InvalidParamCode,
	ErrInvalidCaptcha:       InvalidCaptchaCode,
	ErrNotLogin:             NotLoginCode,
	ErrPhoneEmpty:           PhoneEmptyCode,
	ErrPhoneRegistered:      PhoneRegisteredCode,
	ErrPhoneUnregistered:    PhoneUnregisteredCode,
	ErrCommentNotExists:     CommentNotExistsCode,
	ErrHubExists:            HubExistsCode,
	ErrHubNotExists:         HubNotExistCode,
	ErrMessageInadequate:    MessageInadequateCode,
	ErrDescriptionShort:     DescriptionShortCode,
	ErrDescriptionLong:      DescriptionLongCode,
	ErrCommunityNameEmpty:   CommunityNameEmptyCode,
	ErrCommunityNameLong:    CommunityNameLongCode,
	ErrCommunityNameExists:  CommunityNameExistsCode,
	ErrCommunityNotExists:   CommunityNotExitsCode,
	ErrPostNotExists:        PostNotExistsCode,
	ErrUpload:               UploadCode,
	ErrRequestTooFrequently: RequestTooFrequentlyCode,
	ErrAutoLoginError:       AutoLoginErrorCode,
	ErrCategoryNotExists:    CategoryNotExistsCode,
	ErrCategoryIdExists:     CategoryIdExistsCode,

	ErrServiceBusy:    ServiceBusyCode,
	ErrUserError:      UserErrorCode,
	ErrLoginError:     LoginErrorCode,
	ErrSignupError:    SignupErrorCode,
	ErrSendSmsError:   SendSmsErrorCode,
	ErrCommentError:   CommentErrorCode,
	ErrCommunityError: CommunityErrorCode,
	ErrFeedError:      FeedErrorCode,
	ErrMessageError:   MessageErrorCode,
	ErrRelationError:  RelationErrorCode,
	ErrLikeError:      LikeErrorCode,
	ErrCollectError:   CollectErrorCode,
	ErrPublishError:   PublishErrorCode,
	ErrCategoryError:  CategoryErrorCode,
}

func getCode(err error) int32 {
	code, ok := codeMap[err]
	if !ok {
		return ServiceBusyCode
	}
	return code
}

func Response(c *gin.Context, err error, dataFields map[string]interface{}) {
	statusCode := SuccessCode
	statusMsg := Success
	if err != nil {
		statusMsg = err.Error()
		statusCode = getCode(err)
	}

	// 构建响应 JSON
	response := gin.H{
		"statusCode": statusCode,
		"statusMsg":  statusMsg,
	}

	// 动态设置多个字段
	if dataFields != nil {
		for key, value := range dataFields {
			response[key] = value
		}
	}

	c.JSON(http.StatusOK, response)
	return
}
