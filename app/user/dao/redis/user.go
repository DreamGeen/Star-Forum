package redis

import (
	"context"
	"fmt"
	"time"

	"star/models"
)

//GetCaptcha 获取验证码
func GetCaptcha(phone string) (string, error) {
	return rdb.Get(context.Background(), fmt.Sprintf("captcha:%s", phone)).Result()
}

//Get 获取用户
func Get(ctx context.Context, key string, user *models.User) error {
	return rdb.Get(ctx, key).Scan(user)
}

//Set 储存用户
func Set(ctx context.Context, key string, user *models.User, overtime time.Duration) {
	rdb.Set(ctx, key, user, overtime)
}
