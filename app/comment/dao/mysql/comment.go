package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"star/models"
	"star/utils"
)

// 更新帖子评论数
func updatePostComment(tx *sql.Tx, postId, count int64) error {
	sqlStr := "UPDATE post SET comment = comment + ? WHERE postId = ? "
	_, err := tx.Exec(sqlStr, count, postId)
	if err != nil {
		return err
	}
	return nil
}

// 更新评论回复数
func updateReplyComment(tx *sql.Tx, commentId, count int64) error {
	updateStr := "UPDATE postComment SET reply = reply + ? WHERE commentId = ? "
	_, err := tx.Exec(updateStr, count, commentId)
	if err != nil {
		return err
	}
	return nil
}

// CreateComment 发布评论
func CreateComment(comment *models.Comment) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				log.Println("发布评论回滚事务失败，err:", err)
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				log.Println("发布评论提交事务失败，err:", err)
			}
		}
	}()

	// 判断是否有父评论
	if comment.BeCommentId != 0 {
		// 检查评论是否存在（即未被软删除）
		var exists bool
		checkStr := "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
		err := db.QueryRow(checkStr, comment.BeCommentId).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			// 如果评论不存在（可能已经被软删除）
			return fmt.Errorf("关联评论ID: %d 不存在或已被删除", comment.BeCommentId)
		}
	}

	// 雪花算法生成评论id
	commentId := utils.GetID()
	sqlStr := "INSERT INTO postComment (commentId, postId, userId, content, beCommentId) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, commentId, comment.PostId, comment.UserId, comment.Content, comment.BeCommentId)

	// 更新帖子评论数(+1)
	err = updatePostComment(tx, comment.PostId, 1)
	if err != nil {
		return err
	}

	// 如果存在父评论，则更新其回复数(+1)
	if comment.BeCommentId != 0 {
		err := updateReplyComment(tx, comment.BeCommentId, 1)
		if err != nil {
			return err
		}
	}

	return err
}

// DeleteComment 删除评论入口
func DeleteComment(commentId int64) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				log.Println("删除评论事务回滚失败，err:", err)
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				log.Println("删除评论事务提交失败，err:", err)
			}
		}
	}()

	// 查询对应的postId
	var postId int64
	err = tx.QueryRow("SELECT postId FROM postComment WHERE commentId = ? AND deletedAt IS NULL", commentId).Scan(&postId)
	if err != nil {
		if err == sql.ErrNoRows {
			// 评论不存在或已被删除
			return fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
		}
		return err
	}

	// 递归删除评论
	if err := deleteComments(tx, commentId, postId); err != nil {
		return err
	}

	return nil
}

// 使用递归删除评论，删除评论逻辑
func deleteComments(tx *sql.Tx, commentId int64, postId int64) error {
	// 检查评论是否存在（即未被软删除）
	var exists bool
	checkStr := "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
	err := tx.QueryRow(checkStr, commentId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		// 如果评论不存在（可能已经被软删除）
		return fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
	}

	// 软删除原始评论
	if _, err := tx.Exec("UPDATE postComment SET deletedAt = CURRENT_TIMESTAMP WHERE commentId = ?", commentId); err != nil {
		return err
	}

	// 更新帖子评论数(-1)
	err = updatePostComment(tx, postId, -1)
	if err != nil {
		return err
	}

	// 检查当前评论是否有父评论（即是否是子评论）
	var beCommentId int64
	err = tx.QueryRow("SELECT beCommentId FROM postComment WHERE commentId = ? ", commentId).Scan(&beCommentId)
	if err != nil {
		return err
	}

	// 如果有父评论，则更新父评论的回复数
	if beCommentId != 0 {
		err := updateReplyComment(tx, beCommentId, -1)
		if err != nil {
			return err
		}
	}

	// 查找所有以当前commentId为beCommentId的评论（即直接回复）
	rows, err := tx.Query("SELECT commentId FROM postComment WHERE beCommentId = ? AND deletedAt IS NULL", commentId)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("删除评论rows.Close() err:", err)
		}
	}()

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

	// 递归删除所有子评论
	for _, childId := range childIds {
		if err := deleteComments(tx, childId, postId); err != nil {
			return err
		}
	}

	return nil
}

// GetCommentsStar 获取评论（按照点赞数排序）
func GetCommentsStar(postId int64, page int64, pageSize int64) ([]*models.Comment, error) {
	// 检查帖子是否存在（即未被软删除）
	var exists bool
	checkStr := "SELECT EXISTS(SELECT 1 FROM post WHERE postId = ? AND deletedAt IS NULL)"
	err := db.QueryRow(checkStr, postId).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("帖子查询失败")
	}
	if !exists {
		// 如果帖子不存在（可能已经被软删除）
		return nil, fmt.Errorf("帖子ID: %d 不存在或已被删除", postId)
	}

	// 构造 SQL 查询语句，包含分页和软删除检查
	sqlStr := "SELECT commentId, postId, userId, content, star, reply, beCommentId,createdAt FROM postComment WHERE postId = ? AND deletedAt IS NULL ORDER BY star DESC LIMIT ?, ?"

	// 执行 SQL 查询。
	rows, err := db.Query(sqlStr, postId, (page-1)*pageSize, pageSize)
	if err != nil {
		// 如果查询失败，返回错误。
		return nil, fmt.Errorf("查询数据库错误，err:%v", err)
	}

	// 延迟关闭 rows，确保在函数退出前释放数据库连接
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("rows.Close() err:", err)
		}
	}()

	// 初始化一个空的切片来存储查询到的评论
	var comments []*models.Comment

	// 遍历查询结果。
	for rows.Next() {
		// 创建一个Comment结构体实例来存储当前行的数据
		var comment models.Comment

		// 将当前行的数据扫描到comment实例中
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Content, &comment.Star, &comment.Reply, &comment.BeCommentId, &comment.CreatedAt); err != nil {
			// 如果扫描失败，返回错误。
			return nil, fmt.Errorf("调用Scan错误，err:%v", err)
		}

		// 将comment实例的指针添加到comments切片中
		comments = append(comments, &comment)
	}

	// 检查是否有在迭代过程中发生的错误
	if err := rows.Err(); err != nil {
		// 如果有错误，返回错误。
		return nil, fmt.Errorf("迭代出错，err:%v", err)
	}

	//if len(comments) == 0 {
	//	// 如果没有找到任何有效评论，返回错误
	//	return nil, fmt.Errorf("暂无评论")
	//}

	// 返回查询到的评论列表，如果为空，返回空评论
	return comments, nil
}

// UpdateStar 点赞评论
func UpdateStar(commentId int64, increment int8) error {
	// 检查评论是否存在（即未被软删除）
	var exists bool
	checkStr := "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
	err := db.QueryRow(checkStr, commentId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		// 如果评论不存在（可能已经被软删除）
		return fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
	} else {
		// 如果评论存在且没被删除，则执行点赞
		sqlStr := "UPDATE postComment SET star = star + ? WHERE commentId = ?"
		_, err = db.Exec(sqlStr, increment, commentId)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetStar 获取点赞数
func GetStar(commentId int64) (int64, error) {
	// 创建查询语句
	sqlStr := "SELECT star FROM postComment WHERE commentId = ?"
	// 执行查询
	var starCount int64
	err := db.QueryRow(sqlStr, commentId).Scan(&starCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
		}
		return 0, fmt.Errorf("点赞数更新失败，err:%v", err)
	}

	// 返回点赞数
	return starCount, nil
}
