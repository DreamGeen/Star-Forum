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
	"star/app/storage/cached"
	"star/app/storage/mq"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/constant/str"
	"star/models"
	"star/proto/comment/commentPb"
	"star/proto/like/likePb"
	"star/proto/message/messagePb"
	"star/proto/post/postPb"
	"star/utils"
	"strconv"
)

const (
	post    uint32 = 1
	comment uint32 = 2
)

type LikeSrv struct {
}

var postService postPb.PostService
var messageService messagePb.MessageService
var commentService commentPb.CommentService
var conn *amqp091.Connection
var channel *amqp091.Channel

func failOnError(err error, msg string) {
	if err != nil {
		utils.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		utils.Logger.Error("close rabbitmq conn error", zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		utils.Logger.Error("close rabbitmq channel error", zap.Error(err))
		panic(err)
	}
}
func New() {
	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
	postService = postPb.NewPostService(str.PostService, postMicroService.Client())

	messageMicroService := micro.NewService(micro.Name(str.MessageServiceClient))
	messageService = messagePb.NewMessageService(str.MessageService, messageMicroService.Client())

	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

	var err error
	conn, err = amqp091.Dial(mq.ReturnRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = channel.ExchangeDeclare(str.LikeExchange, "topic", false, false, false, false, nil)
	failOnError(err, "Failed to declare an exchange")

	_, err = channel.QueueDeclare(str.LikePost, false, false, false, false, nil)
	failOnError(err, "Failed to declare a like queue")

	_, err = channel.QueueDeclare(str.LikeComment, false, false, false, false, nil)
	failOnError(err, "Failed to declare a like queue")

	err = channel.QueueBind(str.LikePost, str.RoutPost, str.LikeExchange, false, nil)
	failOnError(err, "Failed to bind a queue to like")

	err = channel.QueueBind(str.LikeComment, str.RoutComment, str.LikeExchange, false, nil)
	failOnError(err, "Failed to bind a queue to like")
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
			utils.Logger.Error("json marshal error", zap.Error(err))
			return
		}
		err = channel.Publish(
			str.LikeExchange,
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
			utils.Logger.Error("json marshal error", zap.Error(err))
			return
		}
		err = channel.Publish(
			str.LikeExchange,
			str.RoutComment,
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			})

	}
	if err != nil {
		utils.Logger.Error("produce like error", zap.Error(err))
		return
	}

}
func (l *LikeSrv) LikeAction(ctx context.Context, req *likePb.LikeActionRequest, resp *likePb.LikeActionResponse) error {
	var err error
	switch req.ActionTye {
	case post:
		err = likePost(ctx, req)
	case comment:
		err = likeComment(ctx, req)
	}
	if err != nil {
		utils.Logger.Error("likeAction service error", zap.Error(err),
			zap.Int64("sourceId", req.SourceId), zap.Uint32("sourceType", req.SourceType),
			zap.Int64("actorId", req.UserId), zap.Uint32("actionType", req.ActionTye))
		return str.ErrLikeError
	}
	return nil
}

func likePost(ctx context.Context, req *likePb.LikeActionRequest) error {
	postExistResp, err := postService.QueryPostExist(ctx, &postPb.QueryPostExistRequest{
		PostId: req.SourceId,
	})
	if err != nil {
		utils.Logger.Error("query post exist error", zap.Error(err),
			zap.Int64("post_id", req.SourceId))
		return err
	}
	if !postExistResp.Exist {
		utils.Logger.Error("post not exist", zap.Int64("post_id", req.SourceId))
		return str.ErrPostNotExists
	}
	user_like_id := fmt.Sprintf("user:%d:like_posts", req.UserId) //用户点赞的作品key
	//贴子信息
	postInfo, err := redis.GetPostInfo(ctx, req.SourceId)
	if err != nil {
		utils.Logger.Error("get post info error", zap.Error(err),
			zap.Int64("post_id", req.SourceId),
			zap.Int64("user_id", req.UserId))
		return str.ErrLikeError
	}
	postId := fmt.Sprintf("%d", req.SourceId)
	//先检查是否重复点赞
	value, err := redis.Client.ZScore(ctx, user_like_id, postId).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("like_post redis service error", zap.Error(err), zap.Int64("post_id", req.SourceId))
		return err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if req.ActionTye == 1 {
		//点赞
		if value > 0 {
			//重复点赞
			utils.Logger.Warn("user duplicate like",
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			return nil
		} else {
			if err := redis.LikePostAction(ctx, req.UserId, postInfo.PostId, postInfo.UserId); err != nil {
				utils.Logger.Error("redis user like post error", zap.Error(err),
					zap.Int64("post_id", req.SourceId),
					zap.Int64("userId", req.UserId))
				return err
			}
			go func() {
				_, err = messageService.SendRemindMessage(ctx, &messagePb.SendRemindMessageRequest{
					SenderId:    req.UserId,
					RecipientId: postInfo.UserId,
					SourceId:    req.SourceId,
					SourceType:  "post",
					RemindType:  "like",
					Content:     postInfo.Content,
					Url:         req.Url,
				})
				if err != nil {
					utils.Logger.Error("send like  remind message error", zap.Error(err),
						zap.Int64("post_id", req.SourceId))
					err = nil
				}
				produceLike(ctx, req)
			}()
		}
	} else {
		//取消点赞
		if value == 0 {
			//用户未点赞
			utils.Logger.Warn("user did not like, cancel liking",
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			return nil
		} else {
			//正常取消点赞
			if err := redis.UnlikePostAction(ctx, req.UserId, postInfo.PostId, postInfo.UserId); err != nil {
				utils.Logger.Error("redis user cancel like post error", zap.Error(err),
					zap.Int64("post_id", req.SourceId),
					zap.Int64("userId", req.UserId))
				return err
			}

		}
	}

	return nil
}

func likeComment(ctx context.Context, req *likePb.LikeActionRequest) error {
	//检查评论是否存在
	commentInfo, err := redis.GetCommentInfo(ctx, req.SourceId)
	if err != nil {
		if errors.Is(err, str.ErrCommentNotExists) {
			utils.Logger.Error("comment not exist", zap.Error(err),
				zap.Int64("commentId", req.SourceId))
			return str.ErrCommentNotExists
		}
		utils.Logger.Error("get comment info error", zap.Error(err),
			zap.Int64("commentId", req.SourceId))
		return str.ErrLikeError
	}
	if req.ActionTye == 1 {
		if err := redis.LikeCommentAction(ctx, commentInfo.CommentId, commentInfo.BeCommentId); err != nil {
			utils.Logger.Error("redis user like post error", zap.Error(err),
				zap.Int64("post_id", req.SourceId),
				zap.Int64("userId", req.UserId))
			return err
		}
		go func() {
			_, err = messageService.SendRemindMessage(ctx, &messagePb.SendRemindMessageRequest{
				SenderId:    req.UserId,
				RecipientId: commentInfo.UserId,
				SourceId:    req.SourceId,
				SourceType:  "post",
				RemindType:  "like",
				Content:     commentInfo.Content,
				Url:         req.Url,
			})
			if err != nil {
				utils.Logger.Error("send like  remind message error", zap.Error(err),
					zap.Int64("post_id", req.SourceId))
				err = nil
			}
			produceLike(ctx, req)
		}()

	} else {
		if err := redis.UnLikeCommentAction(ctx, commentInfo.CommentId, commentInfo.UserId); err != nil {
			utils.Logger.Error("redis user cancel like post error", zap.Error(err),
				zap.Int64("post_id", req.SourceId), zap.Int64("userId", req.UserId))
			return err
		}

	}
	return nil
}

func (l *LikeSrv) GetUserTotalLike(ctx context.Context, req *likePb.GetUserTotalLikeRequest, resp *likePb.GetUserTotalLikeResponse) error {

	key := fmt.Sprintf("user:%d:liked_count", req.UserId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get user total like count error", zap.Error(err),
			zap.Int64("user_id", req.UserId))
		return str.ErrLikeError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		utils.Logger.Error("strconv user total like count error", zap.Error(err),
			zap.Int64("user_id", req.UserId))
		return str.ErrLikeError
	}
	resp.Count = count
	return nil
}

func (l *LikeSrv) LikeList(ctx context.Context, req *likePb.LikeListRequest, resp *likePb.LikeListResponse) error {
	key := fmt.Sprintf("user:%d:like_posts", req.UserId)
	postIdsStr, err := redis.Client.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		utils.Logger.Error("redis get user all like posts id error", zap.Error(err),
			zap.Int64("user_id", req.UserId))
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
		ActorId: req.UserId,
		PostIds: postIds,
	})
	if err != nil {
		utils.Logger.Error("query posts detail error", zap.Error(err),
			zap.Int64("user_id", req.UserId))
		return str.ErrLikeError
	}
	resp.Posts = queryPostsResp.Posts
	return nil
}

func (l *LikeSrv) GetLikeCount(ctx context.Context, req *likePb.GetLikeCountRequest, resp *likePb.GetLikeCountResponse) (err error) {
	var count int64
	switch req.SourceType {
	case post:
		count, err = countPostLike(ctx, req.SourceId)
	case comment:
		count, err = countCommentLike(ctx, req.SourceId)
	}
	if err != nil {
		utils.Logger.Error("GetLikeCount service error", zap.Error(err),
			zap.Int64("sourceId", req.SourceId))
		resp.Count = 0
		return str.ErrLikeError
	}
	resp.Count = count
	return nil
}

func countPostLike(ctx context.Context, postId int64) (count int64, err error) {
	key := fmt.Sprintf("post:%d:liked_count", postId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get post like count error", zap.Error(err),
			zap.Int64("postId", postId))
		return 0, err
	}
	if errors.Is(err, redis2.Nil) {
		return 0, nil
	}
	count, err = strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		utils.Logger.Error("strconv post like count error", zap.Error(err),
			zap.Int64("postId", postId),
			zap.String("countStr", countStr))
		return 0, err
	}
	return count, nil
}

func countCommentLike(ctx context.Context, commentId int64) (count int64, err error) {
	key := fmt.Sprintf("comment:%d:liked_count", commentId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get comment like count error", zap.Error(err))
		return 0, err
	}
	if errors.Is(err, redis2.Nil) {
		return 0, nil
	}
	count, err = strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		utils.Logger.Error("strconv comment like count error", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (l *LikeSrv) GetUserLikeCount(ctx context.Context, req *likePb.GetUserLikeCountRequest, resp *likePb.GetUserLikeCountResponse) error {
	key := fmt.Sprintf("user:%d:like_posts", req.UserId)
	count, err := redis.Client.ZCard(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get user like count error", zap.Error(err),
			zap.Int64("user_id", req.UserId))
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	resp.Count = count
	return nil
}
func (l *LikeSrv) IsLike(ctx context.Context, req *likePb.IsLikeRequest, resp *likePb.IsLikeResponse) (err error) {
	switch req.SourceType {
	case post:
		resp.Result, err = isLikePost(ctx, req.ActorId, req.SourceId)
	case comment:
		resp.Result, err = isLikeComment(ctx, req.ActorId, req.SourceId)
	}
	if err != nil {
		utils.Logger.Error("isLike error", zap.Error(err))
		return str.ErrLikeError
	}
	return nil
}

func isLikePost(ctx context.Context, actorId, postId int64) (bool, error) {
	key := fmt.Sprintf("user:%d:like_posts", actorId)
	postIdStr := fmt.Sprintf("%d", postId)
	ok, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("IsCollect redis service error", zap.Error(err),
			zap.Int64("post_id", postId),
			zap.Int64("userId", actorId))
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

func isLikeComment(ctx context.Context, actorId, commentId int64) (bool, error) {
	cacheKey := fmt.Sprintf("IsLike_actorId:%d_commentId:%d", actorId, commentId)
	countStr, err := cached.GetWithFunc(ctx, cacheKey, func(key string) (string, error) {
		return mysql.IsLikeComment(actorId, commentId)
	})
	if err != nil {
		utils.Logger.Error("is Like comment error", zap.Error(err),
			zap.Int64("actorId", actorId),
			zap.Int64("commentId", commentId))
		return false, err
	}
	return countStr != "0", nil
}
