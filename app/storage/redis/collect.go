package redis

import (
	"context"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"time"
)

func CollectPostAction(ctx context.Context, userId int64, postId int64) error {
	_, err := Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		pipe.IncrBy(ctx, fmt.Sprintf("post:%d:collected_count", postId), 1)
		pipe.ZAdd(ctx, fmt.Sprintf("user:%d:collect_posts", userId), redis2.Z{
			Member: postId,
			Score:  float64(time.Now().Unix()),
		})
		return nil
	})
	return err
}

func UnCollectPostAction(ctx context.Context, userId int64, postId int64) error {
	_, err := Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
		pipe.IncrBy(ctx, fmt.Sprintf("post:%d:collected_count", postId), -1)
		pipe.ZRem(ctx, fmt.Sprintf("user:%d:collect_posts", userId), postId)
		return nil
	})
	return err
}
