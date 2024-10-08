package models

type LikeAction struct {
	UserId     int64  `json:"user_id"`
	SourceId   int64  `json:"source_id" binding:"required"`
	SourceType uint32 `json:"source_type" binding:"required"`
	ActionType uint32 `json:"action_type" binding:"required"`
	Url        string `json:"url"`
}
