package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis_rate/v10"
	redis2 "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"slices"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/app/utils/snowflake"
	"star/proto/publish/publishPb"
)

type PublishSrv struct {
}

const (
	redisPublishQPS = 3
)

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
			zap.Any("feed", post))
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
		if length >= 10 {
			err = pipe.RPopCount(ctx, getCommunityPostByTimeKey, int(length)-6).Err()
			if err != nil {
				logger.Error("clean redis GetCommunityPostByTime error",
					zap.Error(err))
				logging.SetSpanError(span, err)
				return str.ErrPublishError
			}
		}
		FPostsJson, err := pipe.LPop(ctx, getCommunityPostByTimeKey).Result()
		if err != nil {
			logger.Error(" GetCommunityPostByTime redis lpop error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return str.ErrPublishError
		}
		var FPosts []*models.Post
		err = json.Unmarshal([]byte(FPostsJson), &FPosts)
		if err != nil {
			logger.Error("unmarshal feed error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return str.ErrPublishError
		}
		x := (len(FPosts)+1)/20 + 1
		start := 19
		end := 39
		for i := 0; i < x; i++ {
			var segPost []*models.Post
			if i == 0 {
				segPost = slices.Concat([]*models.Post{post}, FPosts[:19])
			} else {
				segPost = FPosts[start:end]
			}
			segPostJson, err := json.Marshal(segPost)
			if err != nil {
				logger.Error("marshal feed error",
					zap.Error(err))
				logging.SetSpanError(span, err)
				continue
			}
			pipe.LPush(ctx, getCommunityPostByTimeKey, segPostJson)
			start += 20
			end += 20
		}
		return nil
	})
	if err != nil {
		logger.Error("update GetCommunityPostByTime redis error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrPublishError
	}
	return nil
}

func (p *PublishSrv) CountPost(ctx context.Context, req *publishPb.CountPostRequest, resp *publishPb.CountPostResponse) error {

	return nil
}
