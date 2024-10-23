package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis_rate/v10"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/app/utils/snowflake"
	"star/proto/feed/feedPb"
	"star/proto/publish/publishPb"
	"strconv"
	"time"
)

type PublishSrv struct {
}

const (
	redisPublishQPS = 3
)

var feedService feedPb.FeedService
var publishSrvIns *PublishSrv

func (p *PublishSrv) New() {
	feedMicroService := micro.NewService(micro.Name(str.FeedServiceClient))
	feedService = feedPb.NewFeedService(str.FeedService, feedMicroService.Client())
}

func publishLimitKey(userId int64) string {
	return fmt.Sprintf("redis_post_limiter:%d", userId)
}

// CreatePost 创建帖子
func (p *PublishSrv) CreatePost(ctx context.Context, req *publishPb.CreatePostRequest, resp *publishPb.CreatePostResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CreatePostService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.CreatePost")

	//redis limit
	limiter := redis_rate.NewLimiter(redis.Client)
	limiterKey := publishLimitKey(req.UserId)
	limitRes, err := limiter.Allow(ctx, limiterKey, redis_rate.PerSecond(redisPublishQPS))
	if err != nil {
		logger.Error("feed limiter error",
			zap.Error(err),
			zap.Int64("actorId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	if limitRes.Allowed == 0 {
		logger.Error("user create feed too frequently",
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrRequestTooFrequently
	}

	post := &models.Post{
		PostId:      snowflake.GetID(),
		UserId:      req.UserId,
		Star:        0,
		Collection:  0,
		Content:     req.Content,
		IsScan:      req.IsScan,
		CommunityId: req.CommunityId,
	}
	if err := mysql.InsertPost(post); err != nil {
		logger.Error("mysql insert feed error",
			zap.Int64("user_id", req.UserId),
			zap.Error(err),
			zap.Any("post", post))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}

	_, err = redis.Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		getCommunityPostByTimeKey := fmt.Sprintf("GetCommunityPostByTime:%d", req.CommunityId)
		length, err := pipe.LLen(ctx, getCommunityPostByTimeKey).Result()
		if err != nil {
			logger.Error("GetCommunityPostByTime redis error",
				zap.Error(err),
				zap.Int64("actorId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrPublishError
		}
		if length >= 300 {
			//获取锁
			lockKey := fmt.Sprintf("Lock_GetCommunityPostByTime:%d", req.CommunityId)
			ok, err := redis.Client.SetNX(ctx, lockKey, 1, 5*time.Second).Result()
			if err != nil {
				logger.Error("get lock error",
					zap.Error(err))
				return str.ErrPublishError
			}
			if ok {
				defer func() {
					_, err := redis.Client.Del(ctx, lockKey).Result()
					if err != nil {
						logger.Error("failed to release lock",
							zap.Error(err))
					}
				}()
				err = pipe.RPopCount(ctx, getCommunityPostByTimeKey, int(length)-120).Err()
				if err != nil {
					logger.Error("clean redis GetCommunityPostByTime error",
						zap.Error(err))
					logging.SetSpanError(span, err)
					return str.ErrPublishError
				}
			}
		}
		postJson, err := json.Marshal(post)
		if err != nil {
			logger.Error("marshal post error",
				zap.Error(err),
				zap.Any("post", post))
			logging.SetSpanError(span, err)
			return str.ErrPublishError
		}
		_, err = pipe.LPush(ctx, getCommunityPostByTimeKey, postJson).Result()
		if err != nil {
			logger.Error(" GetCommunityPostByTime redis lpush error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return str.ErrPublishError
		}
		return nil
	})
	if err != nil {
		logger.Error("update GetCommunityPostByTime redis error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	listPublishKey := fmt.Sprintf("ListPost:%d", req.UserId)
	err = redis.Client.LPush(ctx, listPublishKey, post.PostId).Err()
	if err != nil {
		logger.Error("update user list post redis error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	return nil
}

func (p *PublishSrv) CountPost(ctx context.Context, req *publishPb.CountPostRequest, resp *publishPb.CountPostResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CountPostService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "PublishService.CountPost")

	key := fmt.Sprintf("CountPost:%d", req.UserId)
	countStr, err := cached.GetWithFunc(ctx, key, func(s string) (string, error) {
		return mysql.CountPost(req.UserId)
	})
	if err != nil {
		logger.Error("count user post error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		cached.Delete(ctx, key)
		return str.ErrPublishError
	}
	resp.Count, err = strconv.ParseInt(countStr, 64, 10)
	if err != nil {
		logger.Error("parse countStr error",
			zap.Error(err),
			zap.String("countStr", countStr),
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	return nil
}
func (p *PublishSrv) ListPost(ctx context.Context, req *publishPb.ListPostRequest, resp *publishPb.ListPostResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "ListPostService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "PublishService.ListPost")

	key := fmt.Sprintf("ListPost:%d", req.UserId)
	postsIdStr, err := redis.Client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		logger.Error("redis list user publish postId error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("actorId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	if len(postsIdStr) == 0 {
		return nil
	}
	postsId := make([]int64, len(postsIdStr))
	for i, postIdStr := range postsIdStr {
		postsId[i], _ = strconv.ParseInt(postIdStr, 10, 64)
	}
	qResp, err := feedService.QueryPosts(ctx, &feedPb.QueryPostsRequest{
		ActorId: req.ActorId,
		PostIds: postsId,
	})
	if err != nil {
		logger.Error("query post detail error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.Int64("actorId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	resp.Posts = qResp.Posts
	return nil
}
