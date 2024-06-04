package redis

import (
	"context"
	"fmt"
	"time"
)

const overtime = 300 * time.Second

func SaveCaptcha(captcha string, phone string) {
	rdb.SetEx(context.Background(), fmt.Sprintf("captcha:%s", phone), captcha, overtime)
}
