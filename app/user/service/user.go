package service

import (
	"context"
	"fmt"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/storage/cached"
	"star/constant/str"
	"star/models"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
	"star/utils"
	"sync"
)

var relationService relationPb.RelationService

func New() {
	relationMicroService := micro.NewService(micro.Name(str.RelationServiceClient))
	relationService = relationPb.NewRelationService(str.RelationService, relationMicroService.Client())
}

// GetUserInfo 获取用户具体信息
func (u *UserSrv) GetUserInfo(ctx context.Context, req *userPb.GetUserInfoRequest, resp *userPb.GetUserInfoResponse) error {

	key := fmt.Sprintf("GetUserInfo:%d", req.UserId)
	user := new(models.User)
	found, err := cached.ScanGetUser(ctx, key, user)
	if err != nil {
		utils.Logger.Error("GetUserInfo failed", zap.Error(err))
		return err
	}
	if !found {
		utils.Logger.Info("GetUserInfo err:user not found", zap.Int64("userId", req.UserId))
		return str.ErrUserNotExists
	}
	resp.User = &userPb.User{
		UserId:   user.UserId,
		Exp:      user.Exp,
		Grade:    user.Grade,
		Gender:   user.Gender,
		UserName: user.Username,
		Img:      user.Img,
		Sign:     user.Signature,
		Birth:    user.Birth,
		IsFollow: false,
	}
	var wg sync.WaitGroup
	var isErr bool
	wg.Add(1)
	go func() {
		defer wg.Done()
		isFollowResp, err := relationService.IsFollow(ctx, &relationPb.IsFollowRequest{
			UserId:   req.UserId,
			FollowId: req.ActorId,
		})
		if err != nil {
			utils.Logger.Error("get is follow failed", zap.Error(err), zap.Int64("userId", req.UserId), zap.Any("followId", req.ActorId))
			isErr = true
			return
		}
		resp.User.IsFollow = isFollowResp.Result
	}()
	wg.Wait()
	if isErr {
		return str.ErrUserError
	}
	//返回user信息
	return nil
}
