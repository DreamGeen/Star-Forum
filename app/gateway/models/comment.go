package models

// PostComment 校验评论发布结构体
type PostComment struct {
	PostId      int64  `json:"postId" binding:"required"`
	UserId      int64  `json:"userId" binding:"required"`
	Content     string `json:"content" binding:"required"`
	BeCommentId int64  `json:"beCommentId"`
}
