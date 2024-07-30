package mysql

import (
	"star/models"
)

func CreateComment(comment *models.Comment) error {
	sqlStr := "INSERT INTO postComment (postId, userId, content, beCommentId) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, comment.PostId, comment.UserId, comment.Content, comment.BeCommentId)
	return err
}

func DeleteComment(commentId int64) error {
	sqlStr := "DELETE FROM postComment WHERE commentId = ? OR beCommentId = ?"
	_, err := db.Exec(sqlStr, commentId, commentId)
	return err
}

func GetComments(postId int64, page int64, pageSize int64) ([]*models.Comment, error) {
	sqlStr := "SELECT commentId, postId, userId, content, star, beCommentId FROM postComment WHERE postId = ? LIMIT ?, ?"
	rows, err := db.Query(sqlStr, postId, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Content, &comment.Star, &comment.BeCommentId); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func UpdateStar(commentId int64, increment int) error {
	sqlStr := "UPDATE postComment SET star = star + ? WHERE commentId = ?"
	_, err := db.Exec(sqlStr, increment, commentId)
	return err
}
