package service

import (
	"context"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/constant/str"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
	"star/utils"
	"strconv"
)

type RelationSrv struct {
}

var userService userPb.UserService

func New() {
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())
}

func (r *RelationSrv) Follow(ctx context.Context, req *relationPb.FollowRequest, resp *relationPb.FollowResponse) error {
	if err := updateFollowCountCache(ctx, req.UserId, true); err != nil {
		utils.Logger.Error("update follow count error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		return str.ErrRelationError
	}
	if err := updateFollowerListCache(ctx, req.UserId, false, req.BeFollowerId); err != nil {
		utils.Logger.Error("update follow list error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		return str.ErrRelationError
	}
	if err := updateFansCountCache(ctx, req.UserId, false); err != nil {
		utils.Logger.Error("update fans count error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		return str.ErrRelationError
	}
	if err := updateFansListCache(ctx, req.UserId, false, req.BeFollowerId); err != nil {
		utils.Logger.Error("update fans count error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		return str.ErrRelationError
	}
	beFollowerStatus, err := isFollow(ctx, req.BeFollowerId, req.UserId)
	if err != nil {
		utils.Logger.Error("get is follow error", zap.Error(err),
			zap.Int64("user_id", req.BeFollowerId),
			zap.Int64("be_follower_id", req.UserId))
		return str.ErrRelationError
	}
	if err := mysql.Follow(req.UserId, req.BeFollowerId, beFollowerStatus); err != nil {
		utils.Logger.Error("Follow error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		return str.ErrRelationError
	}
	cached.Delete(ctx, fmt.Sprintf("IsFollow_%d_%d", req.UserId, req.BeFollowerId))
	return nil
}
func (r *RelationSrv) UnFollow(ctx context.Context, req *relationPb.UnFollowRequest, resp *relationPb.UnFollowResponse) error {
	if err := updateFollowCountCache(ctx, req.UserId, true); err != nil {
		utils.Logger.Error("update follow count error",
			zap.Error(err), zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		return str.ErrRelationError
	}
	if err := updateFollowerListCache(ctx, req.UserId, false, req.UnBeFollowerId); err != nil {
		utils.Logger.Error("update follow list error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		return str.ErrRelationError
	}
	if err := updateFansCountCache(ctx, req.UserId, false); err != nil {
		utils.Logger.Error("update fans count error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		return str.ErrRelationError
	}
	if err := updateFansListCache(ctx, req.UserId, false, req.UnBeFollowerId); err != nil {
		utils.Logger.Error("update fans count error",
			zap.Error(err), zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
	}
	unFollowerStatus, err := isFollow(ctx, req.UnBeFollowerId, req.UserId)
	if err != nil {
		utils.Logger.Error("get is follow error", zap.Error(err),
			zap.Int64("user_id", req.UnBeFollowerId),
			zap.Int64("be_follower_id", req.UserId))
		return str.ErrRelationError
	}
	if err := mysql.Unfollow(req.UserId, req.UnBeFollowerId, unFollowerStatus); err != nil {
		utils.Logger.Error("UnFollow error", zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		return str.ErrRelationError
	}
	cached.Delete(ctx, fmt.Sprintf("IsFollow_%d_%d", req.UserId, req.UnBeFollowerId))
	return nil
}

func (r *RelationSrv) GetFollowList(ctx context.Context, req *relationPb.GetFollowRequest, resp *relationPb.GetFollowResponse) error {
	key := fmt.Sprintf("GetFollowerList:%d", req.UserId)
	var followIdList []int64
	followIdStrList, err := redis.Client.SMembers(ctx, key).Result()
	var db bool
	if err != nil {
		utils.Logger.Error("redis "+
			"get followIdStrList error", zap.Error(err), zap.Int64("userId", req.UserId))
		db = true
	} else {
		for _, followIdStr := range followIdStrList {
			followerId, err := strconv.ParseInt(followIdStr, 10, 64)
			if err != nil {
				utils.Logger.Error("parse follower userId error", zap.Error(err), zap.Int64("userId", req.UserId), zap.String("followerId", followIdStr))
				//如果有follower的id不合法，则删除缓存
				if _, err := redis.Client.Del(ctx, key).Result(); err != nil {
					utils.Logger.Error("delete redis followerList error", zap.Error(err), zap.Int64("userId", req.UserId), zap.String("key", key))
				}
				//去mysql里查询
				db = true
				break
			}
			followIdList = append(followIdList, followerId)
		}
		if db {
			followIdList, err = mysql.GetFollowIdList(req.UserId)
			if err != nil {
				utils.Logger.Error("get followerIdList error", zap.Error(err), zap.Int64("userId", req.UserId))
				return str.ErrRelationError
			}
			//将follow存入redis中去
			_, err = redis.Client.Pipelined(ctx, func(pipe redis2.Pipeliner) error {
				for _, followId := range followIdList {
					pipe.SAdd(ctx, key, followId)
				}
				return nil
			})
			if err != nil {
				utils.Logger.Error("redis add followIdList error", zap.Error(err), zap.Int64("userId", req.UserId))
			}
		}
	}
	var followList []*userPb.User
	for _, followId := range followIdList {
		userResponse, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
			UserId: followId,
		})
		if err != nil {
			utils.Logger.Error("get follower Info error", zap.Error(err), zap.Int64("userId", req.UserId), zap.Int64("followerId", followId))
		} else {
			followList = append(followList, userResponse.User)
		}
	}
	resp.FollowList = followList
	return nil
}

func (r *RelationSrv) GetFansList(ctx context.Context, req *relationPb.GetFansListRequest, resp *relationPb.GetFansListResponse) error {
	key := fmt.Sprintf("GetFansList:%d", req.UserId)
	var fansIdList []int64
	fansIdStrList, err := redis.Client.SMembers(ctx, key).Result()
	if err != nil {
		utils.Logger.Error("redis "+
			"get fansStrList error", zap.Error(err), zap.Int64("userId", req.UserId))
		fansIdList, err = mysql.GetFansIdList(req.UserId)
		if err != nil {
			utils.Logger.Error("get fansIdList error", zap.Error(err), zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		//将fansId存入redis中去
		for followerId, _ := range fansIdList {
			redis.Client.SAdd(ctx, key, followerId)
		}
	} else {
		var db bool
		for _, fansIdStr := range fansIdStrList {
			fansId, err := strconv.ParseInt(fansIdStr, 10, 64)
			if err != nil {
				utils.Logger.Error("parse fans userId error", zap.Error(err), zap.Int64("userId", req.UserId), zap.String("fansIdStr", fansIdStr))
				//如果有fans的id不合法，则删除缓存
				if _, err := redis.Client.Del(ctx, key).Result(); err != nil {
					utils.Logger.Error("delete redis fansList error", zap.Error(err), zap.Int64("userId", req.UserId), zap.String("key", key))
				}
				//去mysql里查询
				db = true
				break
			}
			fansIdList = append(fansIdList, fansId)
		}
		if db {
			fansIdList, err = mysql.GetFansIdList(req.UserId)
			if err != nil {
				utils.Logger.Error("get fansIdList error", zap.Error(err), zap.Int64("userId", req.UserId))
				return str.ErrRelationError
			}
			//将fansId存入redis中去
			for fansId, _ := range fansIdList {
				redis.Client.SAdd(ctx, key, fansId)
			}
		}
	}
	var fansList []*userPb.User
	for _, fansId := range fansIdList {
		userResponse, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
			UserId: fansId,
		})
		if err != nil {
			utils.Logger.Error("get fans Info error", zap.Error(err), zap.Int64("userId", req.UserId), zap.Int64("fansId", fansId))
		}
		fansList = append(fansList, userResponse.User)
	}
	resp.FansList = fansList
	return nil
}
func (r *RelationSrv) CountFollow(ctx context.Context, req *relationPb.CountFollowRequest, resp *relationPb.CountFollowResponse) error {
	cacheKey := fmt.Sprintf("CountFollower:%d", req.UserId)
	countStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get follower count from cache error", zap.Error(err), zap.Int64("userId", req.UserId))
		return str.ErrRelationError
	}
	if ok {
		//cache命中
		count, err := strconv.ParseInt(countStr, 64, 10)
		if err != nil {
			utils.Logger.Error("parse follower count error", zap.Error(err), zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		resp.Count = count
		return nil
	} else {
		//cache未命中
		//mysql里查询
		count, err := mysql.GetFollowCount(req.UserId)
		if err != nil {
			utils.Logger.Error("get follower count error", zap.Error(err), zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		resp.Count = count
		countStr = strconv.FormatInt(count, 10)
		cached.Write(ctx, cacheKey, countStr, true)
	}
	return nil
}

func (r *RelationSrv) CountFans(ctx context.Context, req *relationPb.CountFansRequest, resp *relationPb.CountFansResponse) error {
	cacheKey := fmt.Sprintf("CountFans:%d", req.UserId)
	countStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get fans count from cache error", zap.Error(err), zap.Int64("userId", req.UserId))
		return str.ErrRelationError
	}
	if ok {
		//cache命中
		count, err := strconv.ParseInt(countStr, 64, 10)
		if err != nil {
			utils.Logger.Error("parse fans count error", zap.Error(err), zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		resp.Count = count
		return nil
	} else {
		//cache未命中
		//mysql里查询
		count, err := mysql.GetFansCount(req.UserId)
		if err != nil {
			utils.Logger.Error("get fans count error", zap.Error(err), zap.Int64("userId", req.UserId))
			return str.ErrRelationError
		}
		resp.Count = count
		countStr = strconv.FormatInt(count, 10)
		cached.Write(ctx, cacheKey, countStr, true)
	}
	return nil
}

// followOp   true ->follow
// followOp   false->unfollow
func updateFollowCountCache(ctx context.Context, userId int64, followOp bool) error {
	cacheKey := fmt.Sprintf("CountFollower:%d", userId)
	cacheCountStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get follow count from cache error", zap.Error(err), zap.Int64("userId", userId))
		return err
	}
	var count int64
	if ok {
		//cache命中
		count, err = strconv.ParseInt(cacheCountStr, 64, 10)
		if err != nil {
			utils.Logger.Error("parse follower count error", zap.Error(err), zap.Int64("userId", userId))
			return str.ErrRelationError
		}

	} else {
		//cache未命中
		//mysql里查询
		count, err = mysql.GetFollowCount(userId)
		if err != nil {
			utils.Logger.Error("get follower count error", zap.Error(err), zap.Int64("userId", userId))
			return str.ErrRelationError
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

func updateFansCountCache(ctx context.Context, userId int64, followOp bool) error {
	cacheKey := fmt.Sprintf("CountFans:%d", userId)
	cacheCountStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		utils.Logger.Error("get follow count from cache error", zap.Error(err), zap.Int64("userId", userId))
		return err
	}
	var count int64
	if ok {
		//cache命中
		count, err = strconv.ParseInt(cacheCountStr, 64, 10)
		if err != nil {
			utils.Logger.Error("parse fans count error", zap.Error(err), zap.Int64("userId", userId))
			return str.ErrRelationError
		}

	} else {
		//cache未命中
		//mysql里查询
		count, err = mysql.GetFollowCount(userId)
		if err != nil {
			utils.Logger.Error("get fans count error", zap.Error(err), zap.Int64("userId", userId))
			return str.ErrRelationError
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

func updateFollowerListCache(ctx context.Context, userId int64, followOp bool, actorId int64) (err error) {
	key := fmt.Sprintf("GetFollowerList:%d", userId)
	if followOp {
		err = redis.Client.SAdd(ctx, key, actorId).Err()
	} else {
		err = redis.Client.SRem(ctx, key, actorId).Err()
	}
	if err != nil {
		utils.Logger.Error("update user follower list error", zap.Error(err), zap.Int64("userId", userId))
		return err
	}
	return nil
}

func updateFansListCache(ctx context.Context, userId int64, followOp bool, actorId int64) (err error) {
	key := fmt.Sprintf("GetFansList:%d", userId)
	if followOp {
		err = redis.Client.SAdd(ctx, key, actorId).Err()
	} else {
		err = redis.Client.SRem(ctx, key, actorId).Err()
	}
	if err != nil {
		utils.Logger.Error("update user fans list error", zap.Error(err), zap.Int64("userId", userId))
		return err
	}
	return nil
}

func (r *RelationSrv) IsFollow(ctx context.Context, req *relationPb.IsFollowRequest, resp *relationPb.IsFollowResponse) error {
	result, err := isFollow(ctx, req.UserId, req.FollowId)
	if err != nil {
		utils.Logger.Error("get is Follow result error", zap.Error(err), zap.Int64("userId", req.UserId), zap.Int64("followId", req.FollowId))
		resp.Result = false
		return str.ErrRelationError
	}
	resp.Result = result
	return nil
}

func isFollow(ctx context.Context, userId, beFollowerId int64) (bool, error) {
	cacheKey := fmt.Sprintf("IsFollow_%d_%d", userId, beFollowerId)

	countStr, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		return mysql.IsFollow(userId, beFollowerId)
	})
	if err != nil {
		utils.Logger.Error("is follow error", zap.Error(err), zap.Int64("userId", userId))
		return false, err
	}
	return countStr != "0", nil
}
