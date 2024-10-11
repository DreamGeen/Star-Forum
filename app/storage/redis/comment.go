package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand/v2"
	"star/app/models"
	"star/app/storage/mysql"
	"time"
)

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
