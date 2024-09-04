package utils

import (
	"context"
	"log"
	"star/app/storage/cached"
)

// ValidateCaptcha 校验验证码是否正确
func ValidateCaptcha(ctx context.Context, phone, captcha string) bool {
	cacheKey := "captcha:" + phone
	storedCaptcha, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		log.Println("get dao saved captcha fail", err)
		return false
	}
	if !ok {
		log.Println("captcha not exist")
		return false
	}
	if storedCaptcha != captcha {
		log.Println("captcha is wrong", phone)
		return false
	}
	return true
}
