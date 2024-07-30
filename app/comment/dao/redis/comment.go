package redis

import (
	"fmt"
)

func IncrementCommentStar(commentId int64) error {
	key := fmt.Sprintf("comment:star:%d", commentId)
	return Client.Incr(Ctx, key).Err()
}

func GetCommentStar(commentId int64) (int64, error) {
	key := fmt.Sprintf("comment:star:%d", commentId)
	return Client.Get(Ctx, key).Int64()
}

func IncrementCommentReplyCount(commentId int64) error {
	key := fmt.Sprintf("comment:reply:%d", commentId)
	return Client.Incr(Ctx, key).Err()
}

func GetCommentReplyCount(commentId int64) (int64, error) {
	key := fmt.Sprintf("comment:reply:%d", commentId)
	return Client.Get(Ctx, key).Int64()
}
