package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand/v2"
	redis2 "star/app/comment/dao/redis"
	"star/app/storage/mysql"
	"star/models"

	"log"
	"strconv"
	"time"
)

// IncrementCommentStar 点赞评论
func IncrementCommentStar(commentId int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	log.Println("Redis点赞评论成功")
	return redis2.Client.Incr(redis2.Ctx, key).Err()
}

// SetCommentStar 设置评论的点赞数到Redis，并设置过期时间为24小时
func SetCommentStar(commentId int64, starCount int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	// 过期时间24小时
	overtime := 24 * time.Hour
	// 设置键值对和过期时间
	log.Println("设置评论的点赞数到Redis，并设置过期时间为24小时")
	return redis2.Client.Set(redis2.Ctx, key, starCount, overtime).Err()
}

// GetCommentStar 获取点赞数
func GetCommentStar(commentId int64) (int64, error) {
	key := fmt.Sprintf("comment:star:%d", commentId)
	val, err := redis2.Client.Get(redis2.Ctx, key).Result()
	// 缓存未命中
	if err == redis.Nil {
		log.Println("缓存redis未命中", err)
		return 0, fmt.Errorf("缓存redis未命中")
	}
	// Redis客户端错误
	if err != nil {
		log.Println("从redis获取点赞数失败", err)
		return 0, fmt.Errorf("从redis获取点赞数失败：%v", err)
	}

	// 将字符串值转换为 int64
	starCount, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		// 转换错误
		log.Println("将点赞数转换为整数失败", err)
		return 0, fmt.Errorf("将点赞数转换为整数失败：%v", err)
	}

	// 返回点赞数
	log.Println("Redis获取点赞成功")
	return starCount, nil
}

func GetCommentInfo(ctx context.Context, commentId int64) (*models.Comment, error) {
	key := fmt.Sprintf("GetCommentInfo:%d", commentId)
	commentInfoStr, err := Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if errors.Is(err, redis.Nil) {
		comment, err := mysql.GetCommentInfo(commentId)
		if err != nil {
			return nil, err
		}
		commentJson, err := json.Marshal(comment)
		if err != nil {
			return comment, err
		}
		Client.Set(ctx, key, string(commentJson), 10*time.Minute+time.Duration(rand.IntN(60))*time.Second)
		return comment, nil
	}
	comment := &models.Comment{}
	if err := json.Unmarshal([]byte(commentInfoStr), &comment); err != nil {
		return nil, err
	}
	return comment, nil

}

//// CreateComment 存储新发布的评论到Redis
//func CreateComment(comment *models.Comment) error {
//	// 将评论对象序列化为JSON字符串
//	commentJSON, err := json.Marshal(comment)
//	if err != nil {
//		// 如果序列化失败，记录错误并返回
//		logger.CommentLogger.Error("序列化评论失败", zap.Error(err))
//		return err
//	}
//
//	// 定义键名
//	key := fmt.Sprintf("comment:%d", comment.CommentId)
//
//	// 将评论存储到Redis
//	err = Client.Set(context.Background(), key, commentJSON, 0).Err()
//	if err != nil {
//		// 如果存储到Redis失败，记录错误并返回
//		logger.CommentLogger.Error("存储评论到Redis失败", zap.Error(err))
//		return err
//	}
//
//	// 记录成功存储的事件
//	logger.CommentLogger.Info("评论存储到Redis成功", zap.String("key", key))
//	return nil
//}
