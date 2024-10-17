package redis

import (
	"context"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"time"
)

func LikePostAction(ctx context.Context, userId int64, postId int64, userLikedId int64) error {
	_, err := Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		pipe.IncrBy(ctx, fmt.Sprintf("feed:%d:liked_count", postId), 1)
		pipe.IncrBy(ctx, fmt.Sprintf("user:%d:liked_count", userLikedId), 1)
		pipe.ZAdd(ctx, fmt.Sprintf("user:%d:like_posts", userId), redis2.Z{
			Member: float64(postId),
			Score:  float64(time.Now().Unix()),
		})
		return nil
	})
	return err
}

func UnlikePostAction(ctx context.Context, userId int64, postId int64, userLikedId int64) error {
	_, err := Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		pipe.IncrBy(ctx, fmt.Sprintf("feed:%d:liked_count", postId), -1)
		pipe.IncrBy(ctx, fmt.Sprintf("user:%d:liked_count", userLikedId), -1)
		pipe.ZRem(ctx, fmt.Sprintf("user:%d:like_posts", userId), postId)
		return nil
	})
	return err
}

func LikeCommentAction(ctx context.Context, commentId int64, userLikedId int64) error {
	_, err := Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		pipe.IncrBy(ctx, fmt.Sprintf("comment:%d:liked_count", commentId), 1)
		pipe.IncrBy(ctx, fmt.Sprintf("user:%d:liked_count", userLikedId), 1)
		return nil
	})
	return err
}
func UnLikeCommentAction(ctx context.Context, commentId int64, userLikedId int64) error {
	_, err := Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		pipe.IncrBy(ctx, fmt.Sprintf("comment:%d:liked_count", commentId), -1)
		pipe.IncrBy(ctx, fmt.Sprintf("user:%d:liked_count", userLikedId), -1)
		return nil
	})
	return err
}
