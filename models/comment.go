package models

type Comment struct {
	CreatedAt   string `db:"createdAt"`   // 创建时间
	DeletedAt   string `db:"deletedAt"`   // 删除时间
	CommentId   int64  `db:"commentId"`   // 评论id
	PostId      int64  `db:"postId"`      // 帖子id
	UserId      int64  `db:"userId"`      // 用户id
	Content     string `db:"content"`     // 评论内容
	Star        int64  `db:"star"`        // 评论点赞数
	Reply       int64  `db:"reply"`       // 评论回复数
	BeCommentId int64  `db:"beCommentId"` // 关联评论id
}
