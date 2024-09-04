package mysql

import (
	"database/sql"
	"errors"
	"star/constant/str"
	"star/models"
)

const (
	queryPostExistSQL = "select postId from post where postId=? ;"
)

func QueryPostExist(postId int64) (string, error) {
	post := new(models.Post)
	if err := Client.Get(post, queryPostExistSQL, postId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.False, str.ErrPostNotExists
		}
		return str.False, str.ErrPostError
	}
	return str.True, nil
}
