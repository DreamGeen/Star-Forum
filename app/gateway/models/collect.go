package models

type CollectAction struct {
	ActionId   int64  `json:"action_id"`
	PostId     int64  `json:"post_id"`
	ActionType uint32 `json:"action_type" binding:"required"`
}
