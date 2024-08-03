package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// IncrementCommentStar 点赞评论
func IncrementCommentStar(commentId int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	return Client.Incr(Ctx, key).Err()
}

// SetCommentStar 设置评论的点赞数到Redis，并设置过期时间为24小时
func SetCommentStar(commentId int64, starCount int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	// 24小时转换为秒
	overtime := 24 * time.Hour
	// 使用SetEX设置键值对和过期时间
	return Client.Set(Ctx, key, starCount, overtime).Err()
}

// GetCommentStar 获取点赞数
func GetCommentStar(commentId int64) (int64, error) {
	key := fmt.Sprintf("comment:star:%d", commentId)
	val, err := Client.Get(Ctx, key).Result()
	// 缓存未命中
	if err == redis.Nil {
		return 0, fmt.Errorf("缓存redis未命中")
	}
	// Redis客户端错误
	if err != nil {
		return 0, fmt.Errorf("从redis获取点赞数失败：%v", err)
	}

	// 将字符串值转换为 int64
	starCount, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		// 转换错误
		return 0, fmt.Errorf("将点赞数转换为整数失败：%v", err)
	}

	// 返回点赞数
	return starCount, nil
}

func IncrementCommentReplyCount(commentId int64) error {
	key := fmt.Sprintf("comment:reply:%d", commentId)
	return Client.Incr(Ctx, key).Err()
}

func GetCommentReplyCount(commentId int64) (int64, error) {
	key := fmt.Sprintf("comment:reply:%d", commentId)
	return Client.Get(Ctx, key).Int64()
}
