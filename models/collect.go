package models

type Collect struct {
	PostId     int64 `json:"post_id" db:"postId"`
	Collection int64 `json:"collection" db:"collection"`
	UserId     int64 `json:"user_id" db:"UserId"`
}
