package service

import (
	"context"
	"fmt"
	"log"
	"star/app/storage/cached"
	"star/constant/str"
	"star/models"
	"star/proto/user/userPb"
)

// GetUserInfo 获取用户具体信息
func (u *UserSrv) GetUserInfo(ctx context.Context, req *userPb.GetUserInfoRequest, resp *userPb.GetUserInfoResponse) error {

	key := fmt.Sprintf("GetUserInfo:%d", req.UserId)
	user := new(models.User)
	found, err := cached.ScanGetUser(ctx, key, user)
	if err != nil {
		log.Println("GetUserInfo err:", err)
		return err
	}
	if !found {
		log.Println("GetUserInfo err:user not found")
		return str.ErrUserNotExists
	}

	//返回user信息
	resp.User = &userPb.User{
		UserId:   user.UserId,
		Exp:      user.Exp,
		Grade:    user.Grade,
		Gender:   user.Gender,
		UserName: user.Username,
		Img:      user.Img,
		Sign:     user.Signature,
		Birth:    user.Birth,
	}
	return nil
}
