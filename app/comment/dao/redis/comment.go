package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"star/models"
	"star/utils"
	"strconv"
	"time"
)

// IncrementCommentStar 点赞评论
func IncrementCommentStar(commentId int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	utils.Logger.Info("Redis点赞评论成功")
	return Client.Incr(Ctx, key).Err()
}

// SetCommentStar 设置评论的点赞数到Redis，并设置过期时间为24小时
func SetCommentStar(commentId int64, starCount int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	// 24小时转换为秒
	overtime := 24 * time.Hour
	// 使用SetEX设置键值对和过期时间
	utils.Logger.Info("设置评论的点赞数到Redis，并设置过期时间为24小时")
	return Client.Set(Ctx, key, starCount, overtime).Err()
}

// GetCommentStar 获取点赞数
func GetCommentStar(commentId int64) (int64, error) {
	key := fmt.Sprintf("comment:star:%d", commentId)
	val, err := Client.Get(Ctx, key).Result()
	// 缓存未命中
	if err == redis.Nil {
		utils.Logger.Error("缓存redis未命中", zap.Error(err))
		return 0, fmt.Errorf("缓存redis未命中")
	}
	// Redis客户端错误
	if err != nil {
		utils.Logger.Error("从redis获取点赞数失败", zap.Error(err))
		return 0, fmt.Errorf("从redis获取点赞数失败：%v", err)
	}

	// 将字符串值转换为 int64
	starCount, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		// 转换错误
		utils.Logger.Error("将点赞数转换为整数失败", zap.Error(err))
		return 0, fmt.Errorf("将点赞数转换为整数失败：%v", err)
	}

	// 返回点赞数
	utils.Logger.Info("Redis获取点赞成功")
	return starCount, nil
}

// CreateComment 存储新发布的评论到Redis
func CreateComment(comment *models.Comment) error {
	// 将评论对象序列化为JSON字符串
	commentJSON, err := json.Marshal(comment)
	if err != nil {
		// 如果序列化失败，记录错误并返回
		utils.Logger.Error("序列化评论失败", zap.Error(err))
		return err
	}

	// 定义键名
	key := fmt.Sprintf("comment:%d", comment.CommentId)

	// 将评论存储到Redis
	err = Client.Set(context.Background(), key, commentJSON, 0).Err()
	if err != nil {
		// 如果存储到Redis失败，记录错误并返回
		utils.Logger.Error("存储评论到Redis失败", zap.Error(err))
		return err
	}

	// 记录成功存储的事件
	utils.Logger.Info("评论存储到Redis成功", zap.String("key", key))
	return nil
}
