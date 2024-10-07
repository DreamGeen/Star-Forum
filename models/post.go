package models

import (
	"encoding/json"
	"time"
)

type Post struct {
	PostId      int64     `db:"postId"`
	UserId      int64     `db:"userId"`
	CommunityId int64     `db:"communityId"`
	Star        int       `db:"star"`
	Collection  int       `db:"collection"`
	Content     string    `db:"content"`
	IsScan      bool      `db:"isScan"`
	CreateTime  time.Time `db:"createTime"`
	DeleteTime  time.Time `db:"deleteTime"`
}

func (p *Post) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

func (p *Post) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
