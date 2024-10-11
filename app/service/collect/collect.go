package collect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/app/utils/rabbitmq"
	"star/proto/collect/collectPb"
	"star/proto/post/postPb"
	"strconv"
)

type CollectSrv struct {
}

var postService postPb.PostService
var conn *amqp091.Connection
var channel *amqp091.Channel

func failOnError(err error, msg string) {
	if err != nil {
		logging.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		logging.Logger.Error("collect service close rabbitmq conn error",
			zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		logging.Logger.Error("collect service close rabbitmq channel error",
			zap.Error(err))
		panic(err)
	}
}

func New() {
	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
	postService = postPb.NewPostService(str.PostService, postMicroService.Client())

	var err error
	conn, err = amqp091.Dial(rabbitmq.ReturnRabbitmqUrl())
	failOnError(err, "collect service failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "collect service failed to open a channel")

	err = channel.ExchangeDeclare(str.FavorExchange, "topic", false, false, false, false, nil)
	failOnError(err, "collect service failed to declare an exchange")

	_, err = channel.QueueDeclare(str.CollectPost, false, false, false, false, nil)
	failOnError(err, "collect service failed to declare a collect queue")

	err = channel.QueueBind(str.CollectPost, str.RoutCollectPost, str.FavorExchange, false, nil)
	failOnError(err, "collect service failed to bind a queue to favor")

}

func produceCollect(ctx context.Context, req *collectPb.CollectActionRequest) {
	var collection int64
	if req.ActionType == 1 {
		collection = 1
	} else {
		collection = -1
	}
	message := models.Collect{
		PostId:     req.PostId,
		UserId:     req.ActorId,
		Collection: collection,
	}
	msg, err := json.Marshal(message)
	if err != nil {
		logging.Logger.Error("produce collect json marshal error",
			zap.Error(err),
			zap.Int64("postId", req.PostId),
			zap.Int64("actorId", req.ActorId),
			zap.Uint32("actionType", req.ActionType))
		return
	}
	err = channel.Publish(
		str.FavorExchange,
		str.RoutPost,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		logging.Logger.Error("produce  collect message error",
			zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("postId", req.PostId))
	}
	return
}

// IsCollect 是否收藏帖子
func (c *CollectSrv) IsCollect(ctx context.Context, req *collectPb.IsCollectRequest, resp *collectPb.IsCollectResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "IsCollectService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CollectService.IsCollect")

	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdStr := fmt.Sprintf("%d", req.PostId)
	ok, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("IsCollect redis service error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId),
			zap.Int64("userId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrCollectError
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if ok != 0 {
		resp.Result = true
	} else {
		resp.Result = false
	}
	return nil
}

// CollectList 用户收藏帖子列表
func (c *CollectSrv) CollectList(ctx context.Context, req *collectPb.CollectListRequest, resp *collectPb.CollectListResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CollectListService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CollectService.CollectList")

	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdsStr, err := redis.Client.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		logger.Error("redis get user all like posts id error",
			zap.Error(err),
			zap.Int64("user_id", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	if len(postIdsStr) == 0 {
		resp.Posts = nil
		return nil
	}
	postIds := make([]int64, len(postIdsStr))
	for i, postIdStr := range postIdsStr {
		postId, _ := strconv.ParseInt(postIdStr, 10, 64)
		postIds[i] = postId
	}
	queryPostsResp, err := postService.QueryPosts(ctx, &postPb.QueryPostsRequest{
		ActorId: req.ActorId,
		PostIds: postIds,
	})
	if err != nil {
		logger.Error("query posts detail error",
			zap.Error(err),
			zap.Int64("user_id", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	resp.Posts = queryPostsResp.Posts
	return nil
}

// CollectAction 收藏或取消收藏
func (c *CollectSrv) CollectAction(ctx context.Context, req *collectPb.CollectActionRequest, resp *collectPb.CollectActionResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "CollectActionService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CollectService.CollectAction")

	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdStr := fmt.Sprintf("%d", req.PostId)
	value, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("collect_post redis service error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId),
			zap.Int64("userId", req.ActorId))
		logging.SetSpanError(span, err)
		return err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if req.ActionType == 1 {
		//收藏
		if value > 0 {
			//重复收藏
			logger.Warn("user duplicate collect",
				zap.Int64("post_id", req.PostId),
				zap.Int64("userId", req.ActorId))
			return nil
		} else {
			//正常收藏
			err = redis.CollectPostAction(ctx, req.ActorId, req.PostId)
			if err != nil {
				logger.Error("collect_post redis service error",
					zap.Error(err),
					zap.Int64("post_id", req.PostId),
					zap.Int64("userId", req.ActorId))
				logging.SetSpanError(span, err)
				return str.ErrCollectError
			}
			go func() {
				produceCollect(ctx, req)
			}()
		}
	} else {
		//取消收藏
		if value == 0 {
			//用户未点赞
			logger.Warn("user did not collect, cancel collecting",
				zap.Int64("post_id", req.PostId),
				zap.Int64("userId", req.ActorId))
			return nil
		} else {
			err = redis.UnCollectPostAction(ctx, req.ActorId, req.PostId)
			if err != nil {
				logger.Error("collect_post redis service error",
					zap.Error(err),
					zap.Int64("post_id", req.PostId),
					zap.Int64("userId", req.ActorId))
				logging.SetSpanError(span, err)
				return str.ErrCollectError
			}
			go func() {
				produceCollect(ctx, req)
			}()
		}
	}
	return nil
}

// GetCollectCount 获取帖子收藏数
func (c *CollectSrv) GetCollectCount(ctx context.Context, req *collectPb.GetCollectCountRequest, resp *collectPb.GetCollectCountResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetCollectCountService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CollectService.GetCollectCount")

	key := fmt.Sprintf("post:%d:collected_count", req.PostId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("redis get post collect count error",
			zap.Error(err),
			zap.Int64("postId", req.PostId))
		logging.SetSpanError(span, err)
		return str.ErrCollectError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logger.Error("strconv post collect count error",
			zap.Error(err), zap.Int64("postId", req.PostId),
			zap.String("countStr", countStr))
		logging.SetSpanError(span, err)
		return str.ErrCollectError
	}
	resp.Count = count
	return nil
}

// GetUserCollectCount 获取用户收藏帖子数量
func (c *CollectSrv) GetUserCollectCount(ctx context.Context, req *collectPb.GetUserCollectCountRequest, resp *collectPb.GetUserCollectCountResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserCollectCountService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CollectService.GetUserCollectCount")

	key := fmt.Sprintf("user:%d:collect_posts", req.UserId)
	count, err := redis.Client.ZCard(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("redis get user collect count error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrCollectError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	resp.Count = count
	return nil
}
