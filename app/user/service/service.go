package service

import "sync"

type UserSrv struct {
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
