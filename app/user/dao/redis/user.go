package redis

import (
	"context"
	"fmt"
	"star/models"
	"time"
)

// GetCaptcha 获取验证码
func GetCaptcha(phone string) (string, error) {
	return rdb.Get(context.Background(), fmt.Sprintf("captcha:%s", phone)).Result()
}

// GetUser 获取用户
func GetUser(ctx context.Context, key string, user *models.User) error {
	return rdb.Get(ctx, key).Scan(user)
}

// SetUser 储存用户
func SetUser(ctx context.Context, key string, user *models.User, overtime time.Duration) error {
	return rdb.Set(ctx, key, user, overtime).Err()
}
