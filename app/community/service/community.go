package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/constant/str"
	"star/models"
	"star/proto/community/communityPb"
	"star/proto/user/userPb"
	"star/utils"
	"time"
)

type CommunitySrv struct {
}

var userService userPb.UserService

func (c *CommunitySrv) New() {
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

}

func (c *CommunitySrv) CreateCommunity(ctx context.Context, req *communityPb.CreateCommunityRequest, resp *communityPb.EmptyCommunityResponse) error {
	//检查该社区名是否已经存在
	if err := mysql.CheckCommunity(req.CommunityName); err != nil {
		return err
	}
	//构建社区结构体
	community := &models.Community{
		CommunityId:   utils.GetID(),
		CommunityName: req.CommunityName,
		Description:   req.Description,
		LeaderId:      req.LeaderId,
		Img:           str.DefaultCommunityImg,
		Member:        1,
	}
	//将社区插入mysql
	if err := mysql.InsertCommunity(community); err != nil {
		return err
	}
	return nil
}

func (c *CommunitySrv) GetCommunityInfo(ctx context.Context, req *communityPb.GetCommunityInfoRequest, resp *communityPb.GetCommunityInfoResponse) error {
	key := fmt.Sprintf("GetCommunityInfo:%d", req.CommunityId)
	val, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get community info error", zap.Error(err))
		return str.ErrCommunityError
	}
	if errors.Is(err, redis2.Nil) {
		community, err := mysql.GetCommunityInfo(req.CommunityId)
		if err != nil {
			utils.Logger.Error("get community info error", zap.Error(err))
			return str.ErrCommunityError
		}
		userResp, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
			UserId: community.LeaderId,
		})
		if err != nil {
			utils.Logger.Error("get leader info error", zap.Error(err))
			return str.ErrCommunityError
		}
		communityInfo := &communityPb.CommunityInfo{
			CommunityName: community.CommunityName,
			CommunityImg:  community.Img,
			Description:   community.Description,
			Member:        community.Member,
			LeaderName:    userResp.User.UserName,
			LeaderImg:     userResp.User.Img,
		}
		resp.Community = communityInfo
		communityInfoJSON, err := json.Marshal(communityInfo)
		if err != nil {
			utils.Logger.Error("marshal community info error", zap.Error(err))
			return nil
		}
		if err = redis.Client.Set(ctx, key, string(communityInfoJSON), 120*time.Hour).Err(); err != nil {
			utils.Logger.Error("redis set community info error", zap.Error(err))
			return nil
		}
		return nil
	}
	communityInfo := new(communityPb.CommunityInfo)
	if err := json.Unmarshal([]byte(val), communityInfo); err != nil {
		utils.Logger.Error("unmarshal community info error", zap.Error(err))
		return str.ErrCommunityError
	}
	resp.Community = communityInfo
	return nil
}

func (c *CommunitySrv) GetCommunityList(ctx context.Context, req *communityPb.EmptyCommunityRequest, resp *communityPb.GetCommunityListResponse) error {
	//查询community列表

	return nil
}

func (c *CommunitySrv) ShowCommunity(ctx context.Context, req *communityPb.ShowCommunityRequest, resp *communityPb.ShowCommunityResponse) error {
	return nil
}
