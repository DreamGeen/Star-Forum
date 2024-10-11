package password

import "golang.org/x/crypto/bcrypt"

// Encrypt 对密码进行加密
func Encrypt(password string) (string, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encryptPassword), nil
}

// Equals 匹配密码与加密密码是否相等
func Equals(password string, encryptPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password))
	return err
}
