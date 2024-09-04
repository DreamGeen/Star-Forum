package utils

import (
	"errors"
	"star/constant/settings"
	"star/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// token过期时间
const (
	expireAccessToken  = 10 * time.Hour
	expireRefreshToken = 7 * 24 * time.Hour
)

// MyClaims token配置结构体
type MyClaims struct {
	UserID int64 `json:"userid"`
	jwt.RegisteredClaims
}

// 密钥
var key = []byte("冯宇萌")
var (
	tokenExpired = errors.New("jwtAuth is expired")
	tokenInValid = errors.New("jwtAuth invalid")
)

// GetToken 获取token
func GetToken(user *models.User) (accessTokenString string, refreshTokenString string, err error) {
	//获取accessToken
	accessTokenString, err = generateToken(user.UserId, expireAccessToken)
	if err != nil {
		return "", "", err
	}
	//获取refreshToken
	refreshTokenString, err = generateToken(0, expireRefreshToken)
	if err != nil {
		return "", "", err
	}

	return
}

// RefreshAccessToken 刷新accessToken
func RefreshAccessToken(accessToken, refreshToken string) (string, error) {
	//判断refresh token是否合法
	if _, err := ParseToken(refreshToken); err != nil {
		return "", err
	}
	//判断access jwtAuth 是否是因为过期而失效并解析出负载信息
	claims, err := ParseToken(accessToken)
	if !errors.Is(err, tokenExpired) {
		return "", tokenInValid
	}
	return generateToken(claims.UserID, expireAccessToken)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	//基于密钥解析token
	claims := new(MyClaims)
	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, tokenInValid
	}
	//基于过期时间判断token是否合法
	if !jwtToken.Valid {
		return nil, tokenExpired
	}
	return claims, nil
}

// generateToken 生成JWT token
func generateToken(userID int64, expiration time.Duration) (string, error) {
	claims := &MyClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    settings.Conf.AliyunConfig.SignName,            //发行人
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 //发行时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), //过期时间
		},
	}
	//使用jwt签名算法生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//将token进行盐加密
	return token.SignedString(key)
}
