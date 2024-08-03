package mysql

import (
	"database/sql"
	"fmt"
	"star/models"
	"star/utils"
)

// CreateComment 发布评论
func CreateComment(comment *models.Comment) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // 如果发生错误，则回滚事务
		} else {
			tx.Commit() // 如果没有错误，则提交事务
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

	// 如果存在父评论，则更新其回复数
	if comment.BeCommentId != 0 {
		updateStr := "UPDATE postComment SET comment = comment + 1 WHERE commentId = ? AND deletedAt IS NULL"
		_, err = tx.Exec(updateStr, comment.BeCommentId)
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
			tx.Rollback() // 如果发生错误，则回滚事务
		} else {
			tx.Commit() // 如果没有错误，则提交事务
		}
	}()

	// 递归删除评论
	if err := deleteComments(tx, commentId); err != nil {
		return err
	}

	return nil
}

// 使用递归删除评论，删除评论逻辑
func deleteComments(tx *sql.Tx, commentId int64) error {
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

	// 递归删除所有子评论
	for _, childId := range childIds {
		if err := deleteComments(tx, childId); err != nil {
			return err
		}
	}

	return nil
}

// GetComments 获取评论
func GetComments(postId int64, page int64, pageSize int64) ([]*models.Comment, error) {
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
	sqlStr := "SELECT commentId, postId, userId, content, star,comment, beCommentId,createdAt FROM postComment WHERE postId = ? AND deletedAt IS NULL LIMIT ?, ?"

	// 执行 SQL 查询。
	rows, err := db.Query(sqlStr, postId, (page-1)*pageSize, pageSize)
	if err != nil {
		// 如果查询失败，返回错误。
		return nil, fmt.Errorf("查询数据库错误，err:%v", err)
	}

	// 延迟关闭 rows，确保在函数退出前释放数据库连接。
	defer rows.Close()

	// 初始化一个空的切片来存储查询到的评论。
	var comments []*models.Comment

	// 遍历查询结果。
	for rows.Next() {
		// 创建一个 Comment 结构体实例来存储当前行的数据。
		var comment models.Comment

		// 将当前行的数据扫描到 comment 实例中。
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Content, &comment.Star, &comment.Comment, &comment.BeCommentId, &comment.CreatedAt); err != nil {
			// 如果扫描失败，返回错误。
			return nil, fmt.Errorf("调用Scan错误，err:%v", err)
		}

		// 将 comment 实例的指针添加到 comments 切片中。
		comments = append(comments, &comment)
	}

	// 检查是否有在迭代过程中发生的错误
	if err := rows.Err(); err != nil {
		// 如果有错误，返回错误。
		return nil, fmt.Errorf("迭代出错，err:%v", err)
	}

	if len(comments) == 0 {
		// 如果没有找到任何有效评论，返回错误
		return nil, fmt.Errorf("暂无评论")
	}

	// 返回查询到的评论列表和 nil 错误（表示没有错误发生）。
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
