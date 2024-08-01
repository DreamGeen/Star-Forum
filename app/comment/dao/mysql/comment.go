package mysql

import (
	"log"
	"star/models"
	"star/utils"
)

func CreateComment(comment *models.Comment) error {
	commentId := utils.GetID()
	sqlStr := "INSERT INTO postComment (commentId, postId, userId, content, beCommentId) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, commentId, comment.PostId, comment.UserId, comment.Content, comment.BeCommentId)
	return err
}

func DeleteComment(commentId int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println("删除评论事务开始失败", err)
		return err
	}
	defer func() {
		if err != nil {
			log.Println("删除评论事务开始回滚", err)
			if err = tx.Rollback(); err != nil {
				log.Println("删除评论事务回滚失败", err)
			}
		} else {
			log.Println("删除评论事务开始提交")
			if err = tx.Commit(); err != nil {
				log.Println("删除评论事务提交失败", err)
			}
		}
	}()

	// 迭代删除所有相关评论
	for {
		// 查找所有以当前commentId为beCommentId的评论（即直接回复）
		rows, err := tx.Query("SELECT commentId FROM postComment WHERE beCommentId = ?", commentId)
		if err != nil {
			return err
		}
		defer rows.Close()

		var childIds []int64
		for rows.Next() {
			var childId int64
			if err := rows.Scan(&childId); err != nil {
				return err
			}
			childIds = append(childIds, childId)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		// 如果没有找到任何直接回复，则跳出循环
		if len(childIds) == 0 {
			break
		}

		// 删除所有找到的直接回复
		for _, childId := range childIds {
			if _, err := tx.Exec("DELETE FROM postComment WHERE commentId = ?", childId); err != nil {
				return err
			}
		}

		// 更新commentId为最后一个找到的子评论ID，以便在下一次迭代中查找它的回复
		if len(childIds) > 0 {
			commentId = childIds[len(childIds)-1]
		} else {
			// 如果没有子评论了，但还需要删除原始评论本身
			break
		}
	}

	// 删除原始评论
	if _, err := tx.Exec("DELETE FROM postComment WHERE commentId = ?", commentId); err != nil {
		return err
	}

	return nil
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
