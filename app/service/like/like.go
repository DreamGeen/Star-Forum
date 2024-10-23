package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/app/utils/rabbitmq"
	"star/proto/comment/commentPb"
	"star/proto/feed/feedPb"
	"star/proto/like/likePb"
	"star/proto/message/messagePb"
	"strconv"
)

const (
	post    uint32 = 1
	comment uint32 = 2
)

type LikeSrv struct {
}

var feedService feedPb.FeedService
var messageService messagePb.MessageService
var commentService commentPb.CommentService
var conn *amqp091.Connection
var channel *amqp091.Channel
var likeSrvIns *LikeSrv

func failOnError(err error, msg string) {
	if err != nil {
		logging.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		logging.Logger.Error("like service close rabbitmq conn error",
			zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		logging.Logger.Error("like service close rabbitmq channel error",
			zap.Error(err))
		panic(err)
	}
}

func (l *LikeSrv) New() {
	postMicroService := micro.NewService(micro.Name(str.FeedServiceClient))
	feedService = feedPb.NewFeedService(str.FeedService, postMicroService.Client())

	messageMicroService := micro.NewService(micro.Name(str.MessageServiceClient))
	messageService = messagePb.NewMessageService(str.MessageService, messageMicroService.Client())

	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

	var err error
	conn, err = amqp091.Dial(rabbitmq.ReturnRabbitmqUrl())
	failOnError(err, "like service failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "like service failed to open a channel")

	err = channel.ExchangeDeclare(str.FavorExchange, "topic", false, false, false, false, nil)
	failOnError(err, "like service failed to declare an exchange")

	_, err = channel.QueueDeclare(str.LikePost, false, false, false, false, nil)
	failOnError(err, "like service failed to declare a like queue")

	_, err = channel.QueueDeclare(str.LikeComment, false, false, false, false, nil)
	failOnError(err, "like service failed to declare a like queue")

	err = channel.QueueBind(str.LikePost, str.RoutPost, str.FavorExchange, false, nil)
	failOnError(err, "like service failed to bind a queue to like")

	err = channel.QueueBind(str.LikeComment, str.RoutComment, str.FavorExchange, false, nil)
	failOnError(err, "like service failed to bind a queue to like")
}

func produceLike(ctx context.Context, req *likePb.LikeActionRequest) {
	var err error
	var star int
	var msg []byte
	if req.ActionTye == 1 {
		star = 1
	} else {
		star = -1
	}
	switch req.SourceType {
	case post:
		message := models.Post{
			PostId: req.SourceId,
			Star:   star,
		}
		msg, err = json.Marshal(message)
		if err != nil {
			logging.Logger.Error("produce feed like error,json marshal error",
				zap.Error(err),
				zap.Int64("postId", req.SourceId),
				zap.Int64("actorId", req.UserId))
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
	case comment:
		message := models.Comment{
			CommentId: req.SourceId,
			Star:      star,
		}
		msg, err = json.Marshal(message)
		if err != nil {
			logging.Logger.Error("produce comment like error,json marshal error",
				zap.Error(err),
				zap.Int64("comment", req.SourceId),
				zap.Int64("actorId", req.UserId))
			return
		}
		err = channel.Publish(
			str.FavorExchange,
			str.RoutComment,
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			})

	}
	if err != nil {
		logging.Logger.Error("produce like error",
			zap.Error(err),
			zap.Int64("sourceId", req.SourceId),
			zap.Int64("actorId", req.UserId))
		return
	}

}

// LikeAction 点赞
func (l *LikeSrv) LikeAction(ctx context.Context, req *likePb.LikeActionRequest, resp *likePb.LikeActionResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "LikeActionService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "LikeService.LikeAction")

	var err error
	switch req.ActionTye {
	case post:
		err = likePost(ctx, req, span, logger)
	case comment:
		err = likeComment(ctx, req, span, logger)
	}
	if err != nil {
		logger.Error("likeAction service error",
			zap.Error(err),
			zap.Int64("sourceId", req.SourceId),
			zap.Uint32("sourceType", req.SourceType),
			zap.Int64("actorId", req.UserId),
			zap.Uint32("actionType", req.ActionTye))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	return nil
}

func likePost(ctx context.Context, req *likePb.LikeActionRequest, span trace.Span, logger *zap.Logger) error {
	postExistResp, err := feedService.QueryPostExist(ctx, &feedPb.QueryPostExistRequest{
		PostId: req.SourceId,
	})
	if err != nil {
		logger.Error("query feed exist error",
			zap.Error(err),
			zap.Int64("post_id", req.SourceId))
		logging.SetSpanError(span, err)
		return err
	}
	if !postExistResp.Exist {
		logger.Error("feed not exist",
			zap.Int64("post_id", req.SourceId))
		return str.ErrPostNotExists
	}
	user_like_id := fmt.Sprintf("user:%d:like_posts", req.UserId) //用户点赞的作品key
	//贴子信息
	postInfo, err := redis.GetPostInfo(ctx, req.SourceId)
	if err != nil {
		logger.Error("get post info error",
			zap.Error(err),
			zap.Int64("post_id", req.SourceId),
			zap.Int64("user_id", req.UserId))
		logging.SetSpanError(span, err)
		return err
	}
	postId := fmt.Sprintf("%d", req.SourceId)
	//先检查是否重复点赞
	value, err := redis.Client.ZScore(ctx, user_like_id, postId).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("like_post redis service error",
			zap.Error(err),
			zap.Int64("post_id", req.SourceId),
			zap.Int64("actorId", req.UserId))
		logging.SetSpanError(span, err)
		return err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if req.ActionTye == 1 {
		//点赞
		if value > 0 {
			//重复点赞
			logger.Warn("user duplicate like",
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			return nil
		} else {
			if err := redis.LikePostAction(ctx, req.UserId, postInfo.PostId, postInfo.UserId); err != nil {
				logger.Error("redis user like feed error",
					zap.Error(err),
					zap.Int64("post_id", req.SourceId),
					zap.Int64("userId", req.UserId))
				logging.SetSpanError(span, err)
				return err
			}
			go func() {
				_, err = messageService.SendRemindMessage(ctx, &messagePb.SendRemindMessageRequest{
					SenderId:    req.UserId,
					RecipientId: postInfo.UserId,
					SourceId:    req.SourceId,
					SourceType:  "feed",
					RemindType:  "like",
					Content:     postInfo.Content,
					Url:         req.Url,
					IsDeleted:   false,
				})
				if err != nil {
					logger.Error("send like feed remind message error",
						zap.Error(err),
						zap.Int64("post_id", req.SourceId),
						zap.Int64("actorId", req.UserId))
				}
				produceLike(ctx, req)
			}()
		}
	} else {
		//取消点赞
		if value == 0 {
			//用户未点赞
			logger.Warn("user did not like, cancel liking",
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			return nil
		} else {
			//正常取消点赞
			if err := redis.UnlikePostAction(ctx, req.UserId, postInfo.PostId, postInfo.UserId); err != nil {
				logger.Error("redis user cancel like feed error",
					zap.Error(err),
					zap.Int64("post_id", req.SourceId),
					zap.Int64("userId", req.UserId))
				logging.SetSpanError(span, err)
				return err
			}
			go func() {
				_, err = messageService.SendRemindMessage(ctx, &messagePb.SendRemindMessageRequest{
					SenderId:    req.UserId,
					RecipientId: postInfo.UserId,
					SourceId:    req.SourceId,
					SourceType:  "feed",
					RemindType:  "like",
					Content:     postInfo.Content,
					Url:         req.Url,
					IsDeleted:   true,
				})
				if err != nil {
					logger.Error("send unlike feed  remind message error",
						zap.Error(err),
						zap.Int64("post_id", req.SourceId),
						zap.Int64("actorId", req.UserId))
				}
				produceLike(ctx, req)
			}()
		}
	}
	return nil
}

func likeComment(ctx context.Context, req *likePb.LikeActionRequest, span trace.Span, logger *zap.Logger) error {
	//检查评论是否存在
	commentInfo, err := redis.GetCommentInfo(ctx, req.SourceId)
	if err != nil {
		if errors.Is(err, str.ErrCommentNotExists) {
			logger.Error("like comment service error,comment not exist",
				zap.Error(err),
				zap.Int64("commentId", req.SourceId),
				zap.Int64("actorId", req.UserId))
			return str.ErrCommentNotExists
		}
		logger.Error("like comment service error,get comment info error",
			zap.Error(err),
			zap.Int64("commentId", req.SourceId),
			zap.Int64("actorId", req.UserId))
		return err
	}
	if req.ActionTye == 1 {
		if err := redis.LikeCommentAction(ctx, commentInfo.CommentId, commentInfo.BeCommentId); err != nil {
			logger.Error("redis user like comment error",
				zap.Error(err),
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			logging.SetSpanError(span, err)
			return err
		}
		go func() {
			_, err = messageService.SendRemindMessage(ctx, &messagePb.SendRemindMessageRequest{
				SenderId:    req.UserId,
				RecipientId: commentInfo.UserId,
				SourceId:    req.SourceId,
				SourceType:  "feed",
				RemindType:  "like",
				Content:     commentInfo.Content,
				Url:         req.Url,
				IsDeleted:   false,
			})
			if err != nil {
				logger.Error("send like comment remind message error",
					zap.Error(err),
					zap.Int64("commentId", req.SourceId),
					zap.Int64("actorId", req.UserId))
			}
			produceLike(ctx, req)
		}()

	} else {
		if err := redis.UnLikeCommentAction(ctx, commentInfo.CommentId, commentInfo.UserId); err != nil {
			logger.Error("redis user cancel like feed error",
				zap.Error(err),
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			return err
		}
		go func() {
			_, err = messageService.SendRemindMessage(ctx, &messagePb.SendRemindMessageRequest{
				SenderId:    req.UserId,
				RecipientId: commentInfo.UserId,
				SourceId:    req.SourceId,
				SourceType:  "feed",
				RemindType:  "like",
				Content:     commentInfo.Content,
				Url:         req.Url,
				IsDeleted:   true,
			})
			if err != nil {
				logger.Error("send unlike comment remind message error",
					zap.Error(err),
					zap.Int64("commentId", req.SourceId))
			}
			produceLike(ctx, req)
		}()
	}
	return nil
}

// GetUserTotalLike 获取用户总的被赞数
func (l *LikeSrv) GetUserTotalLike(ctx context.Context, req *likePb.GetUserTotalLikeRequest, resp *likePb.GetUserTotalLikeResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserTotalLikeService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "LikeService.GetUserTotalLike")

	key := fmt.Sprintf("user:%d:liked_count", req.UserId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("redis get user total like count error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logger.Error("strconv user total like count error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	resp.Count = count
	return nil
}

// LikeList  用户点赞帖子列表
func (l *LikeSrv) LikeList(ctx context.Context, req *likePb.LikeListRequest, resp *likePb.LikeListResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "LikeListService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "LikeService.LikeList")

	key := fmt.Sprintf("user:%d:like_posts", req.UserId)
	postIdsStr, err := redis.Client.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		logger.Error("redis get user all like posts id error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
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
	queryPostsResp, err := feedService.QueryPosts(ctx, &feedPb.QueryPostsRequest{
		ActorId: req.UserId,
		PostIds: postIds,
	})
	if err != nil {
		logger.Error("query posts detail error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	resp.Posts = queryPostsResp.Posts
	return nil
}

// GetLikeCount 获取点赞数量
func (l *LikeSrv) GetLikeCount(ctx context.Context, req *likePb.GetLikeCountRequest, resp *likePb.GetLikeCountResponse) (err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetLikeCountService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "LikeService.GetLikeCount")

	var count int64
	switch req.SourceType {
	case post:
		count, err = countPostLike(ctx, req.SourceId, span, logger)
	case comment:
		count, err = countCommentLike(ctx, req.SourceId, span, logger)
	}
	if err != nil {
		logger.Error("GetLikeCount service error",
			zap.Error(err),
			zap.Int64("sourceId", req.SourceId))
		logging.SetSpanError(span, err)
		resp.Count = 0
		return str.ErrLikeError
	}
	resp.Count = count
	return nil
}

func countPostLike(ctx context.Context, postId int64, span trace.Span, logger *zap.Logger) (count int64, err error) {
	key := fmt.Sprintf("feed:%d:liked_count", postId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("redis get feed like count error",
			zap.Error(err),
			zap.Int64("postId", postId))
		logging.SetSpanError(span, err)
		return 0, err
	}
	if errors.Is(err, redis2.Nil) {
		return 0, nil
	}
	count, err = strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logger.Error("strconv feed like count error",
			zap.Error(err),
			zap.Int64("postId", postId),
			zap.String("countStr", countStr))
		logging.SetSpanError(span, err)
		return 0, err
	}
	return count, nil
}

func countCommentLike(ctx context.Context, commentId int64, span trace.Span, logger *zap.Logger) (count int64, err error) {
	key := fmt.Sprintf("comment:%d:liked_count", commentId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("redis get comment like count error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return 0, err
	}
	if errors.Is(err, redis2.Nil) {
		return 0, nil
	}
	count, err = strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logger.Error("strconv comment like count error",
			zap.Error(err),
			zap.String("countStr", countStr))
		logging.SetSpanError(span, err)
		return 0, err
	}
	return count, nil
}

// GetUserLikeCount 获取用户点赞帖子数量
func (l *LikeSrv) GetUserLikeCount(ctx context.Context, req *likePb.GetUserLikeCountRequest, resp *likePb.GetUserLikeCountResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserLikeCountService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "LikeService.GetUserLikeCount")

	key := fmt.Sprintf("user:%d:like_posts", req.UserId)
	count, err := redis.Client.ZCard(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("redis get user like count error",
			zap.Error(err),
			zap.Int64("user_id", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	resp.Count = count
	return nil
}

// IsLike 是否点赞
func (l *LikeSrv) IsLike(ctx context.Context, req *likePb.IsLikeRequest, resp *likePb.IsLikeResponse) (err error) {
	ctx, span := tracing.Tracer.Start(ctx, "IsLikeService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "LikeService.IsLike")

	switch req.SourceType {
	case post:
		resp.Result, err = isLikePost(ctx, req.ActorId, req.SourceId, span, logger)
	case comment:
		resp.Result, err = isLikeComment(ctx, req.ActorId, req.SourceId, span, logger)
	}
	if err != nil {
		logger.Error("IsLike service error",
			zap.Error(err),
			zap.Int64("sourceId", req.SourceId),
			zap.Uint32("sourceType", req.SourceType),
			zap.Int64("userId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrLikeError
	}
	return nil
}

func isLikePost(ctx context.Context, actorId, postId int64, span trace.Span, logger *zap.Logger) (bool, error) {
	key := fmt.Sprintf("user:%d:like_posts", actorId)
	postIdStr := fmt.Sprintf("%d", postId)
	ok, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		logger.Error("IsCollect redis service error",
			zap.Error(err),
			zap.Int64("post_id", postId),
			zap.Int64("userId", actorId))
		logging.SetSpanError(span, err)
		return false, err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if ok != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func isLikeComment(ctx context.Context, actorId, commentId int64, span trace.Span, logger *zap.Logger) (bool, error) {
	cacheKey := fmt.Sprintf("IsLike_actorId:%d_commentId:%d", actorId, commentId)
	countStr, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		return mysql.IsLikeComment(actorId, commentId)
	})
	if err != nil {
		logger.Error("is Like comment error",
			zap.Error(err),
			zap.Int64("actorId", actorId),
			zap.Int64("commentId", commentId))
		logging.SetSpanError(span, err)
		return false, err
	}
	return countStr != "0", nil
}
