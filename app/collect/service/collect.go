package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/storage/mq"
	"star/app/storage/redis"
	"star/constant/str"
	"star/models"
	"star/proto/collect/collectPb"
	"star/proto/post/postPb"
	"star/utils"
	"strconv"
)

type CollectSrv struct {
}

var postService postPb.PostService
var conn *amqp091.Connection
var channel *amqp091.Channel

func failOnError(err error, msg string) {
	if err != nil {
		utils.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		utils.Logger.Error("collect service close rabbitmq conn error",
			zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		utils.Logger.Error("collect service close rabbitmq channel error",
			zap.Error(err))
		panic(err)
	}
}

func New() {
	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
	postService = postPb.NewPostService(str.PostService, postMicroService.Client())

	var err error
	conn, err = amqp091.Dial(mq.ReturnRabbitmqUrl())
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
		utils.Logger.Error("produce collect json marshal error",
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
		utils.Logger.Error("produce  collect message error",
			zap.Error(err),
			zap.Int64("userId", req.ActorId),
			zap.Int64("postId", req.PostId))
	}
	return
}
func (c *CollectSrv) IsCollect(ctx context.Context, req *collectPb.IsCollectRequest, resp *collectPb.IsCollectResponse) error {
	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdStr := fmt.Sprintf("%d", req.PostId)
	ok, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("IsCollect redis service error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId),
			zap.Int64("userId", req.ActorId))
		return err
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

func (c *CollectSrv) CollectList(ctx context.Context, req *collectPb.CollectListRequest, resp *collectPb.CollectListResponse) error {
	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdsStr, err := redis.Client.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		utils.Logger.Error("redis get user all like posts id error",
			zap.Error(err),
			zap.Int64("user_id", req.ActorId))
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
		utils.Logger.Error("query posts detail error",
			zap.Error(err),
			zap.Int64("user_id", req.ActorId))
		return str.ErrLikeError
	}
	resp.Posts = queryPostsResp.Posts
	return nil
}

func (c *CollectSrv) CollectAction(ctx context.Context, req *collectPb.CollectActionRequest, resp *collectPb.CollectActionResponse) error {

	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdStr := fmt.Sprintf("%d", req.PostId)
	value, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("collect_post redis service error",
			zap.Error(err),
			zap.Int64("post_id", req.PostId),
			zap.Int64("userId", req.ActorId))
		return err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if req.ActionType == 1 {
		//收藏
		if value > 0 {
			//重复收藏
			utils.Logger.Warn("user duplicate collect",
				zap.Int64("post_id", req.PostId),
				zap.Int64("userId", req.ActorId))
			return nil
		} else {
			//正常收藏
			err = redis.CollectPostAction(ctx, req.ActorId, req.PostId)
			if err != nil {
				utils.Logger.Error("collect_post redis service error",
					zap.Error(err),
					zap.Int64("post_id", req.PostId),
					zap.Int64("userId", req.ActorId))
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
			utils.Logger.Warn("user did not collect, cancel collecting",
				zap.Int64("post_id", req.PostId),
				zap.Int64("userId", req.ActorId))
			return nil
		} else {
			err = redis.UnCollectPostAction(ctx, req.ActorId, req.PostId)
			if err != nil {
				utils.Logger.Error("collect_post redis service error",
					zap.Error(err),
					zap.Int64("post_id", req.PostId),
					zap.Int64("userId", req.ActorId))
				return str.ErrCollectError
			}
			go func() {
				produceCollect(ctx, req)
			}()
		}

	}
	return nil
}

func (c *CollectSrv) GetCollectCount(ctx context.Context, req *collectPb.GetCollectCountRequest, resp *collectPb.GetCollectCountResponse) error {
	key := fmt.Sprintf("post:%d:collected_count", req.PostId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get post collect count error",
			zap.Error(err),
			zap.Int64("postId", req.PostId))
		return str.ErrCollectError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		utils.Logger.Error("strconv post collect count error",
			zap.Error(err), zap.Int64("postId", req.PostId),
			zap.String("countStr", countStr))
		return str.ErrCollectError
	}
	resp.Count = count
	return nil
}
func (c *CollectSrv) GetUserCollectCount(ctx context.Context, req *collectPb.GetUserCollectCountRequest, resp *collectPb.GetUserCollectCountResponse) error {
	key := fmt.Sprintf("user:%d:collect_posts", req.UserId)
	count, err := redis.Client.ZCard(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get user collect count error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
		return str.ErrCollectError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	resp.Count = count
	return nil
}
