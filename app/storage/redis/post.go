package redis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand/v2"
	"star/app/models"
	"star/app/storage/mysql"
	"time"
)

func GetPostInfo(ctx context.Context, postId int64) (*models.Post, error) {
	key := fmt.Sprintf("GetPostAuthorIdAndContent:%d", postId)
	postInfoStr, err := Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if errors.Is(err, redis.Nil) {
		posts, err := mysql.QueryPosts([]int64{postId})
		if err != nil {
			return nil, err
		}
		if len(posts) == 0 {
			return nil, sql.ErrNoRows
		}
		postJson, err := json.Marshal(posts[0])
		if err != nil {
			return posts[0], err
		}
		Client.Set(ctx, key, string(postJson), 72*time.Hour+time.Duration(rand.IntN(60))*time.Minute)
		return posts[0], nil
	}
	post := &models.Post{}
	if err := json.Unmarshal([]byte(postInfoStr), &post); err != nil {
		return nil, err
	}
	return post, nil
}
