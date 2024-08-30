package service

import (
	"star/proto/user/userPb"
	"sync"
)

type UserSrv struct {
	userPb.UserService
}

var (
	userSrvIns *UserSrv
	once       sync.Once
)

func GetUserSrv() *UserSrv {
	once.Do(func() {
		userSrvIns = &UserSrv{}
	})
	return userSrvIns
}
