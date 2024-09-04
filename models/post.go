package models

import "time"

type Post struct {
	PostId      int64     `db:"postId"`
	UserId      int64     `db:"userId"`
	CommunityId int64     `db:"communityId"`
	Star        int       `db:"star"`
	Collection  int       `db:"collection"`
	Title       string    `db:"title"`
	Content     string    `db:"content"`
	IsScan      bool      `db:"isScan"`
	CreateTime  time.Time `db:"createTime"`
	DeleteTime  time.Time `db:"deleteTime"`
}
