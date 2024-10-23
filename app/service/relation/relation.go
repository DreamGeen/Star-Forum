package relation

import (
	"context"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
	"strconv"
)

type RelationSrv struct {
}

var userService userPb.UserService
var relationSrvIns  *RelationSrv

func (r *RelationSrv)New() {
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())
}

// Follow 关注
func (r *RelationSrv) Follow(ctx context.Context, req *relationPb.FollowRequest, resp *relationPb.FollowResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "FollowService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FollowService.Follow")

	if err := updateFollowCountCache(ctx, req.UserId, true, span, logger); err != nil {
		logger.Error("update follow count error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := updateFollowerListCache(ctx, req.UserId, false, req.BeFollowerId, span, logger); err != nil {
		logger.Error("update follow list error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := updateFansCountCache(ctx, req.UserId, false, span, logger); err != nil {
		logger.Error("update fans count error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := updateFansListCache(ctx, req.UserId, false, req.BeFollowerId, span, logger); err != nil {
		logger.Error("update fans count error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	beFollowerStatus, err := isFollow(ctx, req.BeFollowerId, req.UserId, span, logger)
	if err != nil {
		logger.Error("get is follow error",
			zap.Error(err),
			zap.Int64("user_id", req.BeFollowerId),
			zap.Int64("be_follower_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := mysql.Follow(req.UserId, req.BeFollowerId, beFollowerStatus, span, logger); err != nil {
		logger.Error("Follow error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("follower_id", req.BeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	cached.Delete(ctx, fmt.Sprintf("IsFollow_%d_%d", req.UserId, req.BeFollowerId))
	return nil
}

// UnFollow 取消关注
func (r *RelationSrv) UnFollow(ctx context.Context, req *relationPb.UnFollowRequest, resp *relationPb.UnFollowResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "UnFollowService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FollowService.UnFollow")

	if err := updateFollowCountCache(ctx, req.UserId, true, span, logger); err != nil {
		logger.Error("update follow count error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := updateFollowerListCache(ctx, req.UserId, false, req.UnBeFollowerId, span, logger); err != nil {
		logger.Error("update follow list error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := updateFansCountCache(ctx, req.UserId, false, span, logger); err != nil {
		logger.Error("update fans count error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := updateFansListCache(ctx, req.UserId, false, req.UnBeFollowerId, span, logger); err != nil {
		logger.Error("update fans count error",
			zap.Error(err), zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		logging.SetSpanError(span, err)
	}
	unFollowerStatus, err := isFollow(ctx, req.UnBeFollowerId, req.UserId, span, logger)
	if err != nil {
		logger.Error("get is follow error",
			zap.Error(err),
			zap.Int64("user_id", req.UnBeFollowerId),
			zap.Int64("be_follower_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if err := mysql.Unfollow(req.UserId, req.UnBeFollowerId, unFollowerStatus, span, logger); err != nil {
		logger.Error("UnFollow error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("un_follower_id", req.UnBeFollowerId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	cached.Delete(ctx, fmt.Sprintf("IsFollow_%d_%d", req.UserId, req.UnBeFollowerId))
	return nil
}

// GetFollowList 获取关注列表
func (r *RelationSrv) GetFollowList(ctx context.Context, req *relationPb.GetFollowListRequest, resp *relationPb.GetFollowListResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetFollowListService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FollowService.GetFollowList")

	key := fmt.Sprintf("GetFollowList:%d", req.UserId)
	var followIdList []int64
	followIdStrList, err := redis.Client.SMembers(ctx, key).Result()
	var db bool
	if err != nil {
		logger.Error("redis get followIdStrList error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		db = true
	} else {
		for _, followIdStr := range followIdStrList {
			followerId, err := strconv.ParseInt(followIdStr, 10, 64)
			if err != nil {
				logger.Warn("parse follower userId error",
					zap.Error(err),
					zap.Int64("userId", req.UserId),
					zap.String("followerId", followIdStr))
				//如果有follower的id不合法，则删除缓存
				if _, err := redis.Client.Del(ctx, key).Result(); err != nil {
					logger.Error("delete redis followerList error",
						zap.Error(err),
						zap.Int64("userId", req.UserId),
						zap.String("key", key))
				}
				//去mysql里查询
				db = true
				break
			}
			followIdList = append(followIdList, followerId)
		}
	}
	if db {
		followIdList, err = mysql.GetFollowIdList(req.UserId)
		if err != nil {
			logger.Error("get followerIdList error",
				zap.Error(err),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}
		go func() {
			//将follow存入redis中去
			_, err = redis.Client.Pipelined(ctx, func(pipe redis2.Pipeliner) error {
				for _, followId := range followIdList {
					pipe.SAdd(ctx, key, followId)
				}
				return nil
			})
			if err != nil {
				logger.Warn("redis add followIdList error",
					zap.Error(err),
					zap.Int64("userId", req.UserId))
			}
		}()
	}
	var followList []*userPb.User
	for _, followId := range followIdList {
		userResponse, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
			UserId: followId,
		})
		if err != nil {
			logger.Error("get follower Info error",
				zap.Error(err),
				zap.Int64("userId", req.UserId),
				zap.Int64("followerId", followId))
			logging.SetSpanError(span, err)
		} else {
			followList = append(followList, userResponse.User)
		}
	}
	resp.FollowList = followList
	return nil
}

// GetFansList 获取粉丝列表
func (r *RelationSrv) GetFansList(ctx context.Context, req *relationPb.GetFansListRequest, resp *relationPb.GetFansListResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetFansListService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FollowService.GetFansList")

	key := fmt.Sprintf("GetFansList:%d", req.UserId)
	var fansIdList []int64
	fansIdStrList, err := redis.Client.SMembers(ctx, key).Result()
	var db bool
	if err != nil {
		logger.Error("redis get fansStrList error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		db = true
	} else {
		for _, fansIdStr := range fansIdStrList {
			fansId, err := strconv.ParseInt(fansIdStr, 10, 64)
			if err != nil {
				logging.Logger.Warn("parse fans userId error",
					zap.Error(err),
					zap.Int64("userId", req.UserId),
					zap.String("fansIdStr", fansIdStr))
				//如果有fans的id不合法，则删除缓存
				if _, err := redis.Client.Del(ctx, key).Result(); err != nil {
					logging.Logger.Error("delete redis fansList error",
						zap.Error(err),
						zap.Int64("userId", req.UserId),
						zap.String("key", key))
				}
				//去mysql里查询
				db = true
				break
			}
			fansIdList = append(fansIdList, fansId)
		}
	}
	if db {
		fansIdList, err = mysql.GetFansIdList(req.UserId)
		if err != nil {
			logging.Logger.Error("get fansIdList error",
				zap.Error(err),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}
		go func() {
			//将follow存入redis中去
			_, err = redis.Client.Pipelined(ctx, func(pipe redis2.Pipeliner) error {
				for _, fansId := range fansIdList {
					pipe.SAdd(ctx, key, fansId)
				}
				return nil
			})
			if err != nil {
				logger.Warn("redis add followIdList error",
					zap.Error(err),
					zap.Int64("userId", req.UserId))
			}
		}()

	}
	var fansList []*userPb.User
	for _, fansId := range fansIdList {
		userResponse, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
			UserId: fansId,
		})
		if err != nil {
			logger.Error("get fans Info error",
				zap.Error(err),
				zap.Int64("userId", req.UserId),
				zap.Int64("fansId", fansId))
			logging.SetSpanError(span, err)
		}
		fansList = append(fansList, userResponse.User)
	}
	resp.FansList = fansList
	return nil
}

// CountFollow 获取用户关注数量
func (r *RelationSrv) CountFollow(ctx context.Context, req *relationPb.CountFollowRequest, resp *relationPb.CountFollowResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CountFollowService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FollowService.CountFollow")

	cacheKey := fmt.Sprintf("CountFollower:%d", req.UserId)
	countStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("get follower count from cache error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if ok {
		//cache命中
		count, err := strconv.ParseInt(countStr, 64, 10)
		if err != nil {
			logger.Error("parse follower count error",
				zap.Error(err),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}
		resp.Count = count
		return nil
	} else {
		//cache未命中
		//mysql里查询
		count, err := mysql.GetFollowCount(req.UserId)
		if err != nil {
			logger.Error("get follower count error",
				zap.Error(err),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}
		resp.Count = count
		countStr = strconv.FormatInt(count, 10)
		cached.Write(ctx, cacheKey, countStr, true)
	}
	return nil
}

// CountFans 获取用户粉丝数量
func (r *RelationSrv) CountFans(ctx context.Context, req *relationPb.CountFansRequest, resp *relationPb.CountFansResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CountFansService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FollowService.CountFans")

	cacheKey := fmt.Sprintf("CountFans:%d", req.UserId)
	countStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("get fans count from cache error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	if ok {
		//cache命中
		count, err := strconv.ParseInt(countStr, 64, 10)
		if err != nil {
			logger.Error("parse fans count error",
				zap.Error(err),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}
		resp.Count = count
		return nil
	} else {
		//cache未命中
		//mysql里查询
		count, err := mysql.GetFansCount(req.UserId)
		if err != nil {
			logger.Error("get fans count error",
				zap.Error(err),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
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
// 更新关注数量缓存
func updateFollowCountCache(ctx context.Context, userId int64, followOp bool, span trace.Span, logger *zap.Logger) error {
	cacheKey := fmt.Sprintf("CountFollower:%d", userId)
	cacheCountStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("get follow count from cache error",
			zap.Error(err),
			zap.Int64("userId", userId))
		logging.SetSpanError(span, err)
		return err
	}
	var count int64
	if ok {
		//cache命中
		count, err = strconv.ParseInt(cacheCountStr, 64, 10)
		if err != nil {
			logger.Error("parse follower count error",
				zap.Error(err),
				zap.Int64("userId", userId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}

	} else {
		//cache未命中
		//mysql里查询
		count, err = mysql.GetFollowCount(userId)
		if err != nil {
			logger.Error("get follower count error",
				zap.Error(err),
				zap.Int64("userId", userId))
			logging.SetSpanError(span, err)
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

// 更新粉丝数量缓存
func updateFansCountCache(ctx context.Context, userId int64, followOp bool, span trace.Span, logger *zap.Logger) error {
	cacheKey := fmt.Sprintf("CountFans:%d", userId)
	cacheCountStr, ok, err := cached.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("get follow count from cache error",
			zap.Error(err),
			zap.Int64("userId", userId))

		logging.SetSpanError(span, err)
		return err
	}
	var count int64
	if ok {
		//cache命中
		count, err = strconv.ParseInt(cacheCountStr, 64, 10)
		if err != nil {
			logger.Error("parse fans count error",
				zap.Error(err),
				zap.Int64("userId", userId))
			logging.SetSpanError(span, err)
			return str.ErrRelationError
		}

	} else {
		//cache未命中
		//mysql里查询
		count, err = mysql.GetFollowCount(userId)
		if err != nil {
			logger.Error("get fans count error",
				zap.Error(err),
				zap.Int64("userId", userId))
			logging.SetSpanError(span, err)
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

// 更新关注列表缓存
func updateFollowerListCache(ctx context.Context, userId int64, followOp bool, actorId int64, span trace.Span, logger *zap.Logger) (err error) {
	key := fmt.Sprintf("GetFollowerList:%d", userId)
	if followOp {
		err = redis.Client.SAdd(ctx, key, actorId).Err()
	} else {
		err = redis.Client.SRem(ctx, key, actorId).Err()
	}
	if err != nil {
		logger.Error("update user follower list error",
			zap.Error(err),
			zap.Int64("userId", userId))
		logging.SetSpanError(span, err)
		return err
	}
	return nil
}

// 更新粉丝列表缓存
func updateFansListCache(ctx context.Context, userId int64, followOp bool, actorId int64, span trace.Span, logger *zap.Logger) (err error) {
	key := fmt.Sprintf("GetFansList:%d", userId)
	if followOp {
		err = redis.Client.SAdd(ctx, key, actorId).Err()
	} else {
		err = redis.Client.SRem(ctx, key, actorId).Err()
	}
	if err != nil {
		logger.Error("update user fans list error",
			zap.Error(err),
			zap.Int64("userId", userId))
		logging.SetSpanError(span, err)
		return err
	}
	return nil
}

// IsFollow 判断是否关注
func (r *RelationSrv) IsFollow(ctx context.Context, req *relationPb.IsFollowRequest, resp *relationPb.IsFollowResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "IsFollow")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "IsFollow")

	result, err := isFollow(ctx, req.UserId, req.FollowId, span, logger)
	if err != nil {
		logger.Error("get is Follow result error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("followId", req.FollowId))
		logging.SetSpanError(span, err)
		resp.Result = false
		return str.ErrRelationError
	}
	resp.Result = result
	return nil
}

func isFollow(ctx context.Context, userId, beFollowerId int64, span trace.Span, logger *zap.Logger) (bool, error) {
	cacheKey := fmt.Sprintf("IsFollow_%d_%d", userId, beFollowerId)

	countStr, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		return mysql.IsFollow(userId, beFollowerId)
	})
	if err != nil {
		logger.Error("is follow error",
			zap.Error(err),
			zap.Int64("userId", userId))
		logging.SetSpanError(span, err)
		return false, err
	}
	return countStr != "0", nil
}
