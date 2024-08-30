package utils

import "golang.org/x/crypto/bcrypt"

// EncryptPassword 对密码进行加密
func EncryptPassword(password string) (string, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encryptPassword), nil
}

// EqualsPassword 匹配密码与加密密码是否相等
func EqualsPassword(password string, encryptPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password))
	return err
}
