package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/constant/str"
	"star/models"
	"star/proto/community/communityPb"
	"star/proto/user/userPb"
	"star/utils"
	"strconv"
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
		communityInfo := &communityPb.Community{
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
	communityInfo := new(communityPb.Community)
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

func (c *CommunitySrv) FollowCommunity(ctx context.Context, req *communityPb.FollowCommunityRequest, resp *communityPb.FollowCommunityResponse) error {
	err := updateCommunityFollowListCache(ctx, req.ActorId, true, req.CommunityId)
	if err != nil {
		utils.Logger.Error("update community follow list error",
			zap.Error(err), zap.Int64("userId", req.ActorId))
		return str.ErrCommunityError
	}
	err = updateCommunityFollowCountCache(ctx, req.ActorId, true)
	if err != nil {
		utils.Logger.Error("update community follow count error",
			zap.Error(err), zap.Int64("userId", req.ActorId))
		return str.ErrCommunityError
	}
	isFollow, err := isFollowCommunity(ctx, req.ActorId, req.CommunityId)
	if err != nil {
		utils.Logger.Error("isFollow community error", zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		return str.ErrCommunityError
	}
	if isFollow {
		utils.Logger.Warn("user follow community isFollow",
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		return nil
	}
	err = mysql.FollowCommunity(req.CommunityId, req.CommunityId)
	if err != nil {
		utils.Logger.Error("follow community error", zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		return str.ErrCommunityError
	}
	cached.Delete(ctx, fmt.Sprintf("IsFollowCommunity_%d_%d", req.ActorId, req.CommunityId))
	return nil
}

func (c *CommunitySrv) UnFollowCommunity(ctx context.Context, req *communityPb.UnFollowCommunityRequest, resp *communityPb.UnFollowCommunityResponse) error {
	err := updateCommunityFollowListCache(ctx, req.ActorId, false, req.CommunityId)
	if err != nil {
		utils.Logger.Error("update community follow list error", zap.Error(err))
		return err
	}
	err = updateCommunityFollowCountCache(ctx, req.ActorId, false)
	if err != nil {
		utils.Logger.Error("update community follow count error", zap.Error(err))
		return err
	}
	isFollow, err := isFollowCommunity(ctx, req.ActorId, req.CommunityId)
	if err != nil {
		utils.Logger.Error("isFollow community error", zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		return str.ErrCommunityError
	}
	if !isFollow {
		utils.Logger.Warn("user follow community is not follow,cancel following",
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		return nil
	}
	err = mysql.UnFollowCommunity(req.CommunityId, req.CommunityId)
	if err != nil {
		utils.Logger.Error("unfollow community error", zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		return str.ErrCommunityError
	}
	cached.Delete(ctx, fmt.Sprintf("IsFollowCommunity_%d_%d", req.ActorId, req.CommunityId))
	return nil
}

func (c *CommunitySrv) IsFollowCommunity(ctx context.Context, req *communityPb.IsFollowCommunityRequest, resp *communityPb.IsFollowCommunityResponse) error {
	result, err := isFollowCommunity(ctx, req.ActorId, req.CommunityId)
	if err != nil {
		utils.Logger.Error("get is Follow community result error", zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		resp.Result = false
		return str.ErrRelationError
	}
	resp.Result = result
	return nil
}
func isFollowCommunity(ctx context.Context, userId, communityId int64) (bool, error) {
	cacheKey := fmt.Sprintf("IsFollowCommunity_%d_%d", userId, communityId)

	countStr, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		return mysql.IsFollowCommunity(userId, communityId)
	})
	if err != nil {
		utils.Logger.Error("is follow community error", zap.Error(err),
			zap.Int64("userId", userId))
		return false, err
	}
	return countStr != "0", nil
}

func (c *CommunitySrv) CountCommunityFollow(ctx context.Context, req *communityPb.CountCommunityFollowRequest, resp *communityPb.CountCommunityFollowResponse) error {
	cacheKey := fmt.Sprintf("CountCommunityFollow:%d", req.UserId)
	countStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get user community follow count from cache error", zap.Error(err),
			zap.Int64("userId", req.UserId))
		return str.ErrRelationError
	}
	if ok {
		//cache命中
		count, err := strconv.ParseInt(countStr, 64, 10)
		if err != nil {
			utils.Logger.Error("parse community follow count error", zap.Error(err),
				zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		resp.Count = count
		return nil
	} else {
		//cache未命中
		//mysql里查询
		count, err := mysql.CountCommunityFollow(req.UserId)
		if err != nil {
			utils.Logger.Error("get community follow count error", zap.Error(err),
				zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		resp.Count = count
		countStr = strconv.FormatInt(count, 10)
		cached.Write(ctx, cacheKey, countStr, true)
	}
	return nil
}

func (c *CommunitySrv) GetFollowCommunityList(ctx context.Context, req *communityPb.GetFollowCommunityListRequest, resp *communityPb.GetFollowCommunityListResponse) error {
	key := fmt.Sprintf("GetCommunityFollowList:%d", req.UserId)
	var communityIdList []int64
	communityIdStrList, err := redis.Client.SMembers(ctx, key).Result()
	var db bool
	if err != nil {
		utils.Logger.Error("redis get communityIdStrList error",
			zap.Error(err), zap.Int64("userId", req.UserId))
		db = true
	} else {
		for _, communityIdStr := range communityIdStrList {
			communityId, err := strconv.ParseInt(communityIdStr, 10, 64)
			if err != nil {
				utils.Logger.Error("parse communityId error", zap.Error(err),
					zap.Int64("userId", req.UserId),
					zap.String("communityIdStr", communityIdStr))
				//如果有follower的id不合法，则删除缓存
				if _, err := redis.Client.Del(ctx, key).Result(); err != nil {
					utils.Logger.Error("delete redis communityList error", zap.Error(err),
						zap.Int64("userId", req.UserId),
						zap.String("key", key))
				}
				//去mysql里查询
				db = true
				break
			}
			communityIdList = append(communityIdList, communityId)
		}
		if db {
			communityIdList, err = mysql.GetCommunityFollowId(req.UserId)
			if err != nil {
				utils.Logger.Error("get communityIdList error", zap.Error(err),
					zap.Int64("userId", req.UserId))
				return str.ErrRelationError
			}
			//将communityId存入redis中去
			_, err = redis.Client.Pipelined(ctx, func(pipe redis2.Pipeliner) error {
				for _, communityId := range communityIdList {
					pipe.SAdd(ctx, key, communityId)
				}
				return nil
			})
			if err != nil {
				utils.Logger.Error("redis add followIdList error", zap.Error(err),
					zap.Int64("userId", req.UserId))
			}
		}
	}
	var communityList []*communityPb.Community
	for _, communityId := range communityIdList {
		communityInfoResp := new(communityPb.GetCommunityInfoResponse)
		err = c.GetCommunityInfo(ctx, &communityPb.GetCommunityInfoRequest{
			CommunityId: communityId,
		}, communityInfoResp)
		if err != nil {
			utils.Logger.Error("get community Info error", zap.Error(err),
				zap.Int64("userId", req.UserId),
				zap.Int64("communityId", communityId))
		} else {
			communityList = append(communityList, communityInfoResp.Community)
		}
	}
	resp.CommunityList = communityList
	return nil
}

// followOp   true ->follow
// followOp   false->unfollow
func updateCommunityFollowCountCache(ctx context.Context, userId int64, followOp bool) error {
	cacheKey := fmt.Sprintf("CountCommunityFollow:%d", userId)
	cacheCountStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get user comunity follow count from cache error", zap.Error(err),
			zap.Int64("userId", userId))
		return err
	}
	var count int64
	if ok {
		//cache命中
		count, err = strconv.ParseInt(cacheCountStr, 64, 10)
		if err != nil {
			utils.Logger.Error("parse community follow count error", zap.Error(err),
				zap.Int64("userId", userId))
			return str.ErrCommunityError
		}

	} else {
		//cache未命中
		//mysql里查询
		count, err = mysql.CountCommunityFollow(userId)
		if err != nil {
			utils.Logger.Error("get  community follow count error", zap.Error(err),
				zap.Int64("userId", userId))
			return str.ErrCommunityError
		}
	}
	if followOp {
		count++
	} else {
		count--
	}
	cacheCountStr = strconv.FormatInt(count, 10)
	cached.Write(ctx, cacheKey, cacheCountStr, true)
	return nil
}

func updateCommunityFollowListCache(ctx context.Context, userId int64, followOp bool, communityId int64) (err error) {
	key := fmt.Sprintf("GetCommunityFollowList:%d", userId)
	if followOp {
		err = redis.Client.SAdd(ctx, key, communityId).Err()
	} else {
		err = redis.Client.SRem(ctx, key, communityId).Err()
	}
	if err != nil {
		utils.Logger.Error("update user community follow list error", zap.Error(err),
			zap.Int64("userId", userId))
		return err
	}
	return nil
}
