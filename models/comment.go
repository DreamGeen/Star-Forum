package models

type Comment struct {
	CreatedAt   string `db:"createdAt"`
	DeletedAt   string `db:"deletedAt"`
	CommentId   int64  `db:"commentId"`
	PostId      int64  `db:"postId"`
	UserId      int64  `db:"userId"`
	Content     string `db:"content"`
	Star        int64  `db:"star"`
	Comment     int64  `db:"comment"`
	BeCommentId *int64 `db:"beCommentId"`
}
