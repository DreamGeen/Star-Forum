package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"star/models"
	"star/utils"
)

// 更新帖子评论数
func updatePostComment(tx *sqlx.Tx, postId, count int64) error {
	sqlStr := "UPDATE post SET comment = comment + ? WHERE postId = ? "
	_, err := tx.Exec(sqlStr, count, postId)
	if err != nil {
		utils.Logger.Error("更新帖子评论数失败", zap.Error(err), zap.Int64("postId", postId))
		return err
	}
	return nil
}

// 更新评论回复数
func updateReplyComment(tx *sqlx.Tx, commentId, count int64) error {
	sqlStr := "UPDATE postComment SET reply = reply + ? WHERE commentId = ? "
	_, err := tx.Exec(sqlStr, count, commentId)
	if err != nil {
		utils.Logger.Error("更新评论回复数失败", zap.Error(err), zap.Int64("commentId", commentId))
		return err
	}
	return nil
}

// CreateComment 发布评论
func CreateComment(comment *models.Comment) error {
	utils.Logger.Info("开始发布评论")
	// 开始事务
	tx, err := db.Beginx()
	if err != nil {
		utils.Logger.Error("发布评论事务开启失败", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				utils.Logger.Error("发布评论回滚事务失败", zap.Error(err))
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				utils.Logger.Error("发布评论事务提交失败", zap.Error(err))
			}
		}
	}()

	// 检查一下帖子是否存在
	var exists bool
	sqlStr := "SELECT EXISTS(SELECT 1 FROM post WHERE postId = ? AND deletedAt IS NULL)"
	if err = tx.Get(&exists, sqlStr, comment.PostId); err != nil {
		utils.Logger.Error("评论发布:检查帖子是否存在失败", zap.Error(err))
		return err
	}
	if !exists {
		// 如果帖子不存在
		utils.Logger.Error("评论发布:帖子不存在", zap.Int64("postId", comment.PostId))
		return fmt.Errorf("帖子ID: %d 不存在或已被删除", comment.PostId)
	}

	// 判断是否有父评论
	if comment.BeCommentId != 0 {
		// 检查评论是否存在
		var exists bool
		sqlStr := "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
		err := tx.Get(&exists, sqlStr, comment.BeCommentId)
		if err != nil {
			utils.Logger.Error("评论发布:检查是否有父评论失败", zap.Error(err))
			return err
		}
		if !exists {
			// 如果评论不存在
			utils.Logger.Error("评论发布:尝试关联的父评论ID不存在或已被删除", zap.Int64("beCommentId", comment.BeCommentId))
			return fmt.Errorf("关联评论ID: %d 不存在或已被删除", comment.BeCommentId)
		} else {
			// 检查父评论的帖子ID是否与当前评论的帖子ID相同
			var parentPostId int64
			sqlStr := "SELECT postId FROM postComment WHERE commentId = ? AND deletedAt IS NULL"
			err := tx.Get(&parentPostId, sqlStr, comment.BeCommentId)
			if err != nil {
				utils.Logger.Error("评论发布:检查父评论的帖子ID失败", zap.Error(err))
				return err
			}
			if parentPostId != comment.PostId {
				utils.Logger.Error("评论发布:父评论的帖子ID与当前评论的帖子ID不一致", zap.Int64("parentPostId", parentPostId), zap.Int64("commentPostId", comment.PostId))
				return fmt.Errorf("父评论的帖子ID: %d 与当前评论的帖子ID: %d 不一致", parentPostId, comment.PostId)
			}
		}
	}

	// 雪花算法生成评论id
	commentId := utils.GetID()
	sqlStr = "INSERT INTO postComment (commentId, postId, userId, content, beCommentId) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, commentId, comment.PostId, comment.UserId, comment.Content, comment.BeCommentId)
	if err != nil {
		utils.Logger.Error("评论发布:评论插入数据库失败", zap.Error(err))
		return err
	}

	// 更新帖子评论数(+1)
	err = updatePostComment(tx, comment.PostId, 1)
	if err != nil {
		utils.Logger.Error("评论发布:更新帖子评论数失败", zap.Error(err))
		return err
	}

	// 如果存在父评论，则更新其回复数(+1)
	if comment.BeCommentId != 0 {
		err := updateReplyComment(tx, comment.BeCommentId, 1)
		if err != nil {
			utils.Logger.Error("更新父评论回复数失败", zap.Error(err))
			return err
		}
	}

	utils.Logger.Info("评论发布成功", zap.Int64("commentId", commentId))

	return err
}

// DeleteComment 删除评论入口
func DeleteComment(commentId int64) error {
	utils.Logger.Info("开始删除评论")
	// 开始事务
	tx, err := db.Beginx()
	if err != nil {
		utils.Logger.Error("开启事务失败", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				utils.Logger.Error("删除评论事务回滚失败", zap.Error(err))
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				utils.Logger.Error("删除评论事务提交失败", zap.Error(err))
			}
		}
	}()

	// 查询对应的postId
	var postId int64
	sqlStr := "SELECT postId FROM postComment WHERE commentId = ? AND deletedAt IS NULL"
	err = db.Get(&postId, sqlStr, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 评论不存在或已被删除
			utils.Logger.Error("尝试删除的评论ID不存在或已被删除", zap.Int64("commentId", commentId))
			return fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
		}
		utils.Logger.Error("评论删除:查询评论关联的postId失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("准备递归删除评论")

	// 递归删除评论
	if err := deleteComments(tx, commentId, postId); err != nil {
		utils.Logger.Error("递归删除评论失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("评论ID: %d 删除成功", zap.Int64("commentId", commentId))
	return nil
}

// 使用递归删除评论，删除评论逻辑
func deleteComments(tx *sqlx.Tx, commentId int64, postId int64) error {
	// 检查评论是否存在（即未被软删除）
	var exists bool
	sqlStr := "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
	err := tx.Get(&exists, sqlStr, commentId)
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
	sqlStr = "SELECT beCommentId FROM postComment WHERE commentId = ?"
	err = tx.Get(&beCommentId, sqlStr, commentId)
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
	sqlStr = "SELECT commentId FROM postComment WHERE beCommentId = ? AND deletedAt IS NULL"
	rows, err := tx.Query(sqlStr, commentId)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			utils.Logger.Error("递归删除评论:关闭rows失败", zap.Error(err))
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
	utils.Logger.Info("开始获取评论")
	// 检查帖子是否存在（即未被软删除）
	var exists bool
	sqlStr := "SELECT EXISTS(SELECT 1 FROM post WHERE postId = ? AND deletedAt IS NULL)"
	err := db.Get(&exists, sqlStr, postId)
	if err != nil {
		utils.Logger.Error("评论获取:帖子查询失败", zap.Error(err))
		return nil, fmt.Errorf("帖子查询失败")
	}
	if !exists {
		// 如果帖子不存在（可能已经被软删除）
		utils.Logger.Error("评论获取:帖子ID不存在或已被删除", zap.Int64("postId", postId))
		return nil, fmt.Errorf("帖子ID: %d 不存在或已被删除", postId)
	}

	// 构造SQL查询语句，包含分页和软删除检查
	sqlStr = "SELECT commentId, postId, userId, content, star, reply, beCommentId, createdAt FROM postComment WHERE postId = ? AND deletedAt IS NULL ORDER BY star DESC, createdAt DESC LIMIT ?, ?"

	// 执行SQL查询
	rows, err := db.Query(sqlStr, postId, (page-1)*pageSize, pageSize)
	if err != nil {
		// 如果查询失败，返回错误。
		utils.Logger.Error("评论获取:查询数据库错误", zap.Error(err))
		return nil, fmt.Errorf("查询数据库错误，err:%v", err)
	}

	// 延迟关闭rows，确保在函数退出前释放数据库连接
	defer func() {
		if err := rows.Close(); err != nil {
			utils.Logger.Error("评论获取:关闭rows时出错", zap.Error(err))
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
			// 如果扫描失败，返回错误
			utils.Logger.Error("评论获取:调用Scan错误", zap.Error(err))
			return nil, fmt.Errorf("调用Scan错误，err:%v", err)
		}

		// 将comment实例的指针添加到comments切片中
		comments = append(comments, &comment)
	}

	// 检查是否有在迭代过程中发生的错误
	if err := rows.Err(); err != nil {
		// 如果有错误，返回错误
		utils.Logger.Error("评论获取:迭代出错", zap.Error(err))
		return nil, fmt.Errorf("迭代出错，err:%v", err)
	}

	//if len(comments) == 0 {
	//	// 如果没有找到任何有效评论，返回错误
	//	return nil, fmt.Errorf("暂无评论")
	//}

	// 返回查询到的评论列表，如果为空，返回空评论
	utils.Logger.Info("评论获取成功", zap.Int64("postId", postId))
	return comments, nil
}

// UpdateStar 点赞评论
func UpdateStar(commentId int64, increment int8) error {
	utils.Logger.Info("开始点赞评论")
	// 检查评论是否存在（即未被软删除）
	var exists bool
	sqlStr := "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
	err := db.Get(&exists, sqlStr, commentId)
	if err != nil {
		utils.Logger.Error("评论点赞:检查评论存在时出错", zap.Error(err))
		return err
	}
	if !exists {
		// 如果评论不存在（可能已经被软删除）
		utils.Logger.Error("评论点赞:评论不存在或已被删除", zap.Int64("commentId", commentId))
		return fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
	} else {
		// 如果评论存在且没被删除，则执行点赞
		sqlStr := "UPDATE postComment SET star = star + ? WHERE commentId = ?"
		_, err = db.Exec(sqlStr, increment, commentId)
		if err != nil {
			utils.Logger.Error("更新点赞数时出错", zap.Error(err))
			return err
		}
	}

	utils.Logger.Info("评论点赞成功", zap.Int64("commentId", commentId))
	return nil
}

// GetStar 获取点赞数
func GetStar(commentId int64) (int64, error) {
	// 创建查询语句
	sqlStr := "SELECT star FROM postComment WHERE commentId = ? AND deletedAt IS NULL"
	// 执行查询
	var starCount int64
	err := db.Get(&starCount, sqlStr, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Error("获取点赞数:评论不存在或已被删除", zap.Int64("commentId", commentId))
			return 0, fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
		}
		utils.Logger.Error("获取点赞数时出错", zap.Error(err))
		return 0, fmt.Errorf("点赞数更新失败，err:%v", err)
	}

	// 返回点赞数
	return starCount, nil
}
