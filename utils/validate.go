package utils

import (
	"errors"
	"log"

	"star/app/user/dao/redis"
)

var (
	ErrCaptchaWrong = errors.New("验证码错误")
)

// ValidateCaptcha 校验验证码是否正确
func ValidateCaptcha(phone, captcha string) error {
	storedCaptcha, err := redis.GetCaptcha(phone)
	if err != nil {
		log.Println("get dao saved captcha fail", err)
		return ErrCaptchaWrong
	}
	if storedCaptcha != captcha {
		log.Println("captcha is wrong", phone)
		return ErrCaptchaWrong
	}
	return nil
}
