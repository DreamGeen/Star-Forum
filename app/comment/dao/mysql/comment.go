package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
	"star/constant/str"
	"star/models"
	"star/utils"
)

const (
	checkComment       = "SELECT EXISTS(SELECT 1 FROM postComment WHERE commentId = ? AND deletedAt IS NULL)"
	checkPost          = "SELECT EXISTS(SELECT 1 FROM post WHERE postId = ? AND deletedAt IS NULL)"
	checkUser          = "SELECT EXISTS(SELECT 1 FROM user WHERE userId = ? AND deletedAt IS NULL)"
	updatePostComment  = "UPDATE post SET comment = comment + ? WHERE postId = ? "
	updateReplyComment = "UPDATE postComment SET reply = reply + ? WHERE commentId = ? "
	queryPostId        = "SELECT postId FROM postComment WHERE commentId = ? AND deletedAt IS NULL"
	insertComment      = "INSERT INTO postComment (commentId, postId, userId, content, beCommentId) VALUES (?, ?, ?, ?, ?)"
	deleteComment      = "UPDATE postComment SET deletedAt = CURRENT_TIMESTAMP WHERE commentId = ?"
	queryBeCommentId   = "SELECT beCommentId FROM postComment WHERE commentId = ?"
	queryReply         = "SELECT commentId FROM postComment WHERE beCommentId = ? AND deletedAt IS NULL"
	queryCommentsStar  = "SELECT commentId, postId, userId, content, star, reply, beCommentId, createdAt FROM postComment WHERE postId = ? AND deletedAt IS NULL ORDER BY star DESC, createdAt DESC"
	queryCommentsTime  = "SELECT commentId, postId, userId, content, star, reply, beCommentId, createdAt FROM postComment WHERE postId = ? AND deletedAt IS NULL ORDER BY createdAt DESC, star DESC"
	starComment        = "UPDATE postComment SET star = star + ? WHERE commentId = ?"
	queryStar          = "SELECT star FROM postComment WHERE commentId = ? AND deletedAt IS NULL"
)

// CheckComment 检查评论是否存在
func CheckComment(commentId int64) error {
	var exists bool
	err := db.Get(&exists, checkComment, commentId)
	if err != nil {
		log.Println("检查评论是否存在失败", err)
		return str.ErrCommentError
	}
	if !exists {
		// 如果评论不存在
		log.Println("评论不存在或已被删除", "commentId", commentId)
		return str.ErrCommentError
	}
	return nil
}

// CheckPost 检查帖子是否存在
func CheckPost(postId int64) error {
	// 检查一下帖子是否存在
	var exists bool
	if err := db.Get(&exists, checkPost, postId); err != nil {
		log.Println("检查帖子是否存在失败", err)
		return str.ErrCommentError
	}
	if !exists {
		// 如果帖子不存在
		log.Println("帖子不存在或已被删除", "postId", postId)
		return str.ErrCommentNotExists
	}
	return nil
}

// CheckUser 检查用户是否存在
func CheckUser(userId int64) error {
	// 检查一下用户是否存在
	var exists bool
	if err := db.Get(&exists, checkUser, userId); err != nil {
		log.Println("检查用户是否存在失败", err)
		return err
	}
	if !exists {
		// 如果用户不存在
		log.Println("用户不存在或已被删除", "userId", userId)
		return fmt.Errorf("用户ID: %d 不存在或已被删除", userId)
	}
	return nil
}

// UpdatePostComment 更新帖子评论数
func UpdatePostComment(tx *sqlx.Tx, postId, count int64) error {
	_, err := tx.Exec(updatePostComment, count, postId)
	if err != nil {
		log.Println("更新帖子评论数失败", err, "postId", postId)
		return err
	}
	return nil
}

// UpdateReplyComment 更新评论回复数
func UpdateReplyComment(tx *sqlx.Tx, commentId, count int64) error {
	_, err := tx.Exec(updateReplyComment, count, commentId)
	if err != nil {
		log.Println("更新评论回复数失败", err, "commentId", commentId)
		return err
	}
	return nil
}

// CreateComment 发布评论
func CreateComment(comment *models.Comment) error {
	log.Println("开始发布评论")
	// 开始事务
	tx, err := db.Beginx()
	if err != nil {
		log.Println("发布评论事务开启失败", err)
		return str.ErrCommentError
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				log.Println("发布评论回滚事务失败", err)
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				log.Println("发布评论事务提交失败", err)
			}
		}
	}()

	// 检查帖子是否存在
	if err := CheckPost(comment.PostId); err != nil {
		return err
	}

	// 检查用户是否存在
	if err := CheckUser(comment.UserId); err != nil {
		return err
	}

	// 判断是否有父评论
	if comment.BeCommentId != 0 {
		// 检查父评论是否存在
		if err := CheckComment(comment.BeCommentId); err != nil {
			return err
		} else {
			// 检查父评论的帖子ID是否与当前评论的帖子ID相同
			var parentPostId int64
			err := tx.Get(&parentPostId, queryPostId, comment.BeCommentId)
			if err != nil {
				log.Println("评论发布:检查父评论的帖子ID失败", err)
				return err
			}
			if parentPostId != comment.PostId {
				log.Println("评论发布:父评论的帖子ID与当前评论的帖子ID不一致", "parentPostId", parentPostId, "commentPostId", comment.PostId)
				return fmt.Errorf("父评论的帖子ID: %d 与当前评论的帖子ID: %d 不一致", parentPostId, comment.PostId)
			}
		}
	}

	// 雪花算法生成评论id
	commentId := utils.GetID()
	_, err = db.Exec(insertComment, commentId, comment.PostId, comment.UserId, comment.Content, comment.BeCommentId)
	if err != nil {
		log.Println("评论发布:评论插入数据库失败", err)
		return err
	}

	// 更新帖子评论数(+1)
	err = UpdatePostComment(tx, comment.PostId, 1)
	if err != nil {
		log.Println("评论发布:更新帖子评论数失败", err)
		return err
	}

	// 如果存在父评论，则更新其回复数(+1)
	if comment.BeCommentId != 0 {
		err := UpdateReplyComment(tx, comment.BeCommentId, 1)
		if err != nil {
			log.Println("更新父评论回复数失败", err)
			return err
		}
	}

	log.Println("评论发布成功", "commentId", commentId)

	return err
}

// DeleteComment 删除评论入口
func DeleteComment(commentId int64) error {
	log.Println("开始删除评论")
	// 开始事务
	tx, err := db.Beginx()
	if err != nil {
		log.Println("开启事务失败", err)
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				log.Println("删除评论事务回滚失败", zap.Error(err))
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				log.Println("删除评论事务提交失败", err)
			}
		}
	}()

	// 查询对应的postId
	var postId int64
	err = db.Get(&postId, queryPostId, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 评论不存在或已被删除
			log.Println("尝试删除的评论ID不存在或已被删除", "commentId", commentId)
			return fmt.Errorf("评论ID: %d 不存在或已被删除", commentId)
		}
		log.Println("评论删除:查询评论关联的postId失败", err)
		return err
	}

	log.Println("准备递归删除评论")

	// 递归删除评论
	if err := deleteComments(tx, commentId, postId); err != nil {
		log.Println("递归删除评论失败", err)
		return err
	}

	log.Println("评论ID:  删除成功", "commentId", commentId)
	return nil
}

// 使用递归删除评论
func deleteComments(tx *sqlx.Tx, commentId int64, postId int64) error {
	// 检查评论是否存在
	if err := CheckComment(commentId); err != nil {
		return err
	}

	// 软删除原始评论
	if _, err := tx.Exec(deleteComment, commentId); err != nil {
		return err
	}

	// 更新帖子评论数(-1)
	err := UpdatePostComment(tx, postId, -1)
	if err != nil {
		return err
	}

	// 检查当前评论是否有父评论
	var beCommentId int64
	err = tx.Get(&beCommentId, queryBeCommentId, commentId)
	if err != nil {
		return err
	}

	// 如果有父评论，则更新父评论的回复数
	if beCommentId != 0 {
		err := UpdateReplyComment(tx, beCommentId, -1)
		if err != nil {
			return err
		}
	}

	// 查找所有以当前commentId为beCommentId的评论（即直接回复）
	rows, err := tx.Query(queryReply, commentId)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("递归删除评论:关闭rows失败", err)
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
func GetCommentsStar(postId int64) ([]*models.Comment, error) {
	log.Println("开始获取评论")

	// 执行SQL查询，包含分页和软删除检查
	rows, err := db.Query(queryCommentsStar, postId)
	if err != nil {
		// 如果查询失败，返回错误。
		log.Println("评论获取:查询数据库错误", err)
		return nil, fmt.Errorf("查询数据库错误，err:%v", err)
	}

	// 延迟关闭rows，确保在函数退出前释放数据库连接
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("评论获取:关闭rows时出错", err)
		}
	}()

	// 初始化一个空的切片来存储查询到的评论
	var comments []*models.Comment

	// 遍历查询结果
	for rows.Next() {
		// 创建一个Comment结构体实例来存储当前行的数据
		var comment models.Comment

		// 将当前行的数据扫描到comment实例中
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Content, &comment.Star, &comment.Reply, &comment.BeCommentId, &comment.CreatedAt); err != nil {
			// 如果扫描失败，返回错误
			log.Println("评论获取:调用Scan错误", err)
			return nil, fmt.Errorf("调用Scan错误，err:%v", err)
		}

		// 将comment实例的指针添加到comments切片中
		comments = append(comments, &comment)
	}

	// 检查是否有在迭代过程中发生的错误
	if err := rows.Err(); err != nil {
		// 如果有错误，返回错误
		log.Println("评论获取:迭代出错", err)
		return nil, fmt.Errorf("迭代出错，err:%v", err)
	}

	//if len(comments) == 0 {
	//	// 如果没有找到任何有效评论，返回错误
	//	return nil, fmt.Errorf("暂无评论")
	//}

	// 返回查询到的评论列表，如果为空，返回空评论
	log.Println("评论获取成功", "postId", postId)
	return comments, nil
}

// UpdateStar 点赞评论
func UpdateStar(commentId int64, increment int8) error {
	log.Println("开始点赞评论")
	_, err := db.Exec(starComment, increment, commentId)
	if err != nil {
		log.Println("更新点赞数时出错", err)
		return err
	}

	log.Println("评论点赞成功", "commentId", commentId)
	return nil
}

// GetStar 获取点赞数
func GetStar(commentId int64) (int64, error) {
	// 执行查询
	var starCount int64
	err := db.Get(&starCount, queryStar, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("获取点赞数:评论不存在或已被删除", "commentId", commentId)
			return 0, str.ErrCommentNotExists
		}
		log.Println("获取点赞数时出错", err)
		return 0, str.ErrCommentError
	}

	// 返回点赞数
	return starCount, nil
}
