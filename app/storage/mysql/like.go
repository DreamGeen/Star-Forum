package mysql

import (
	"database/sql"
	"errors"
	"strconv"
)

const (
	isLikeCommentSQL     = "select count(1) from like_remind where source_id=? and source_type='comment'and sender_id=?"
	updatePostLikeSQL    = "update  post set  star=star+? where postId=? "
	updateCommentLikeSQL = "update  postComment set  star=star+? where commentId=? "
)

func IsLikeComment(actorId, commentId int64) (string, error) {
	var count int64
	if err := Client.Get(&count, isLikeCommentSQL, commentId, actorId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "0", nil
		}
		return "0", err
	}
	countStr := strconv.FormatInt(count, 64)
	return countStr, nil
}

func UpdatePostLike(postId int64, count int) error {
	if _, err := Client.Exec(updatePostLikeSQL, count, postId); err != nil {
		return err
	}
	return nil
}
func UpdateCommentLike(commentId int64, count int) error {
	if _, err := Client.Exec(updateCommentLikeSQL, count, commentId); err != nil {
		return err
	}
	return nil
}
