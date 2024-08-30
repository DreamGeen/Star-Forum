package utils

import (
	"log"
	"star/app/user/dao/redis"
)

// ValidateCaptcha 校验验证码是否正确
func ValidateCaptcha(phone, captcha string) bool {
	storedCaptcha, err := redis.GetCaptcha(phone)
	if err != nil {
		log.Println("get dao saved captcha fail", err)
		return false
	}
	if storedCaptcha != captcha {
		log.Println("captcha is wrong", phone)
		return false
	}
	return true
}
