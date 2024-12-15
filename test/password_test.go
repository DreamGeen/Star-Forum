package test

import (
	"star/app/utils/password"
	"testing"
)

func TestPassword(t *testing.T) {
	encryptPassword := "$2a$04$gdvhEzvEqyEQFuNA9hMbSuHphcfyQwGM4TuA.mGomUY8fN9j9fwhy"
	mypassword := "8176445027a"
	t.Log(password.Equals(mypassword, encryptPassword))
}
