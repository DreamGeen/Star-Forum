package models

type CreatePost struct {
	UserId      int64  `json:"user_id"`
	CommunityId int64  `json:"community_id" binding:"required"`
	Content     string `json:"content"  binding:"required"`
	IsScan      bool   `json:"is_scan"  binding:"required"`
}

type GetCommunityPost struct {
	CommunityId int64 `json:"community_id" binding:"required"`
	Page        int64 `json:"page"  binding:"required"`
	ActorId     int64  `json:"actor_id" `
	LastPostId  int64  `json:"last_post_id"`
}

type GetCommunityPostByNewReply struct {
	CommunityId int64 `json:"community_id" binding:"required"`
	Page        int64 `json:"page"  binding:"required"`
	ActorId     int64  `json:"actor_id" `
	LastRelyTime string `json:"last_reply_time"`
}

