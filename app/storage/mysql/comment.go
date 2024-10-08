package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"star/constant/str"
	"star/models"
	"star/utils"
	"strconv"
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
	countCommentSQL    = "SELECT comment_count FROM post WHERE postId = ?"
	getCommentInfoSQL  = "SELECT  createdAt, commentId, postId, userId, content, star, reply, beCommentId   FROM  postcomment  WHERE commentId=? AND  deletedAt IS NULL;"
)

// CheckComment 检查评论是否存在
func CheckComment(commentId int64) error {
	var exists bool
	err := Client.Get(&exists, checkComment, commentId)
	if err != nil {
		return err
	}
	if !exists {
		// 如果评论不存在
		return str.ErrCommentNotExists
	}
	return nil
}

// CheckPost 检查帖子是否存在
func CheckPost(postId int64) error {
	// 检查一下帖子是否存在
	var exists bool
	if err := Client.Get(&exists, checkPost, postId); err != nil {
		return err
	}
	if !exists {
		// 如果帖子不存在
		return str.ErrCommentNotExists
	}
	return nil
}

// CheckUser 检查用户是否存在
func CheckUser(userId int64) error {
	// 检查一下用户是否存在
	var exists bool
	if err := Client.Get(&exists, checkUser, userId); err != nil {
		return err
	}
	if !exists {
		// 如果用户不存在
		return fmt.Errorf("用户ID: %d 不存在或已被删除", userId)
	}
	return nil
}

// UpdatePostComment 更新帖子评论数
func UpdatePostComment(tx *sqlx.Tx, postId, count int64) error {
	_, err := tx.Exec(updatePostComment, count, postId)
	if err != nil {
		return err
	}
	return nil
}

// UpdateReplyComment 更新评论回复数
func UpdateReplyComment(tx *sqlx.Tx, commentId, count int64) error {
	_, err := tx.Exec(updateReplyComment, count, commentId)
	if err != nil {
		return err
	}
	return nil
}

// CreateComment 发布评论
func CreateComment(comment *models.Comment) error {
	// 开始事务
	tx, err := Client.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				return
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				return
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
				return err
			}
			if parentPostId != comment.PostId {
				return fmt.Errorf("父评论的帖子ID: %d 与当前评论的帖子ID: %d 不一致", parentPostId, comment.PostId)
			}
		}
	}

	// 雪花算法生成评论id
	commentId := utils.GetID()
	_, err = Client.Exec(insertComment, commentId, comment.PostId, comment.UserId, comment.Content, comment.BeCommentId)
	if err != nil {
		return err
	}

	// 更新帖子评论数(+1)
	err = UpdatePostComment(tx, comment.PostId, 1)
	if err != nil {
		return err
	}

	// 如果存在父评论，则更新其回复数(+1)
	if comment.BeCommentId != 0 {
		err := UpdateReplyComment(tx, comment.BeCommentId, 1)
		if err != nil {
			return err
		}
	}

	return err
}

// DeleteComment 删除评论入口
func DeleteComment(commentId int64) error {
	// 开始事务
	tx, err := Client.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// 如果发生错误，则回滚事务
			if err = tx.Rollback(); err != nil {
				return
			}
		} else {
			// 如果没有错误，则提交事务
			if err = tx.Commit(); err != nil {
				return
			}
		}
	}()

	// 查询对应的postId
	var postId int64
	err = Client.Get(&postId, queryPostId, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
			return
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

	// 执行SQL查询，包含分页和软删除检查
	rows, err := Client.Query(queryCommentsStar, postId)
	if err != nil {
		// 如果查询失败，返回错误。
		return nil, fmt.Errorf("查询数据库错误，err:%v", err)
	}

	// 延迟关闭rows，确保在函数退出前释放数据库连接
	defer func() {
		if err := rows.Close(); err != nil {
			return
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
			return nil, fmt.Errorf("调用Scan错误，err:%v", err)
		}

		// 将comment实例的指针添加到comments切片中
		comments = append(comments, &comment)
	}

	// 检查是否有在迭代过程中发生的错误
	if err := rows.Err(); err != nil {
		// 如果有错误，返回错误
		return nil, fmt.Errorf("迭代出错，err:%v", err)
	}

	//if len(comments) == 0 {
	//	// 如果没有找到任何有效评论，返回错误
	//	return nil, fmt.Errorf("暂无评论")
	//}

	// 返回查询到的评论列表，如果为空，返回空评论
	return comments, nil
}

func CountComment(postId int64) (string, error) {
	var count int
	if err := Client.Get(&count, countCommentSQL, postId); err != nil {
		return "", err
	}
	return strconv.Itoa(count), nil
}

func GetCommentInfo(commentId int64) (*models.Comment, error) {
	comment := new(models.Comment)
	if err := Client.Get(comment, getCommentInfoSQL, commentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, str.ErrCommentNotExists
		}
		return nil, err
	}
	return comment, nil
}
