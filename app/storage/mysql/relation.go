package mysql

import (
	"database/sql"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/utils/logging"
	"strconv"
	"time"
)

const (
	getFollowerIdListSQL = "select be_followed_id from user_follows  where user_id=? and deletedAt IS NULL"
	checkFollowExistSQL  = "select count(1) from user_follows where user_id=? and be_followed_id=?"
	followExistSQL       = "update user_follows set  deletedAt=NULL, status=? where user_id=? and be_followed_id=?"
	fansExistSQL         = "update  user_fans set deletedAt=NULL,status=? where user_id=? and  fans_id=? "
	followUnExistSQL     = "insert into user_follows(user_id,be_followed_id,status)values (?,?,?)"
	fansUnExist          = "insert into user_follows(user_id,fans_id,status)values (?,?,?)"
	unFollowSQL          = "update user_follows set  deletedAt=?, status=false where user_id=? and be_followed_id=?"
	unFansSQL            = "update  user_fans set deletedAt=?,status=false where user_id=? and fans_id=?"
	unBothFollowSQL      = "update user_follows set status=false where user_id=? and be_followed_id=?"
	unBothFansSQL        = "update  user_fans set status=false where user_id=? and fans_id=?"
	getFansIdListSQL     = "select fans_id from user_fans where user_id=? and deletedAt IS NULL limit 100"
	getFollowerCountSQL  = "select count(1) from user_follows where user_id=? and deletedAt IS NOT NULL"
	getFansCountSQL      = "select count(1) from user_fans where user_id=? and deletedAt IS NULL"
	isFollowSQL          = "select count(1) from user_follows where user_id=? and be_followed_id=? and  deletedAt IS NULL"
)

func GetFollowIdList(userId int64) ([]int64, error) {
	var followerIdList []int64
	if err := Client.Select(&followerIdList, getFollowerIdListSQL, userId); err != nil {
		return nil, err
	}
	return followerIdList, nil
}

func GetFansIdList(userId int64) ([]int64, error) {
	var fansIdList []int64
	if err := Client.Select(&fansIdList, getFansIdListSQL, userId); err != nil {
		return nil, err
	}
	return fansIdList, nil
}

func GetFollowCount(userId int64) (int64, error) {
	var count int64
	if err := Client.Get(&count, getFollowerCountSQL, userId); err != nil {
		return 0, err
	}
	return count, nil
}

func GetFansCount(userId int64) (int64, error) {
	var count int64
	if err := Client.Get(&count, getFansCountSQL, userId); err != nil {
		return 0, err
	}
	return count, nil
}

func Follow(userId, beFollowerId int64, beFollowerStatus bool, span trace.Span, logger *zap.Logger) error {
	tx, err := Client.Beginx()
	if err != nil {
		logger.Error("start follow transaction error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	defer func() {
		if p := recover(); p != nil {
			logger.Error(" follow recovered from panic, transaction rolled back:",
				zap.Any("panic", p))
			logging.SetSpanError(span, err)
			tx.Rollback()
			err = str.ErrRelationError
		} else if err != nil {
			logger.Error(" follow transaction rolled back due to error:",
				zap.Error(err))
			logging.SetSpanError(span, err)
			tx.Rollback()
		}
	}()
	//检查follow记录是否存在
	var count int
	if err := tx.Get(&count, checkFollowExistSQL, userId, beFollowerId); err != nil {
		logger.Error("check user follow exist error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("beFollowerId", beFollowerId))
		logging.SetSpanError(span, err)
		return err
	}
	if count > 0 {
		if _, err = tx.Exec(followExistSQL, beFollowerStatus, userId, beFollowerId); err != nil {
			logger.Error("update  user_follows deleteTime to null error",
				zap.Error(err),
				zap.Int64("userId", userId),
				zap.Int64("beFollowerId", beFollowerId))
			logging.SetSpanError(span, err)
			return err
		}
		if _, err = tx.Exec(fansExistSQL, beFollowerStatus, beFollowerId, userId); err != nil {
			logger.Error("update user_fans deleteTime to null error",
				zap.Error(err),
				zap.Int64("userId", userId),
				zap.Int64("beFollowerId", beFollowerId))
			logging.SetSpanError(span, err)
			return err
		}
	} else {
		if _, err = tx.Exec(followUnExistSQL, userId, beFollowerId, beFollowerStatus); err != nil {
			logger.Error("insert follow error",
				zap.Error(err),
				zap.Int64("userId", userId),
				zap.Int64("beFollowerId", beFollowerId))
			logging.SetSpanError(span, err)
			return err
		}
		if _, err = tx.Exec(fansUnExist, beFollowerId, userId, beFollowerStatus); err != nil {
			logger.Error("insert fans error ",
				zap.Error(err),
				zap.Int64("userId", beFollowerId),
				zap.Int64("fansId", userId))
			logging.SetSpanError(span, err)
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		logger.Error("commit follow transaction error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return err
	}
	return nil
}
func IsFollow(userId, beFollowerId int64) (string, error) {
	var count int64
	if err := Client.Get(&count, isFollowSQL, userId, beFollowerId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "0", nil
		}
		return "0", err
	}
	countStr := strconv.FormatInt(count, 10)
	return countStr, nil
}

func Unfollow(userId, unFollowerId int64, unFollowerStatus bool, span trace.Span, logger *zap.Logger) error {
	tx, err := Client.Beginx()
	if err != nil {
		logger.Error("start unfollow transaction error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrRelationError
	}
	defer func() {
		if p := recover(); p != nil {
			logger.Error(" unfollow recovered from panic, transaction rolled back:",
				zap.Any("panic", p))
			logging.SetSpanError(span, err)
			tx.Rollback()
			err = str.ErrRelationError
		} else if err != nil {
			logger.Error(" unfollow transaction rolled back due to error:",
				zap.Error(err))
			logging.SetSpanError(span, err)
			tx.Rollback()
		}
	}()
	if _, err = tx.Exec(unFollowSQL, time.Now().UTC(), userId, unFollowerId); err != nil {
		logger.Error("update user_follows deleteTime to now error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("unFollowerId", unFollowerId))
		logging.SetSpanError(span, err)
		return err
	}
	if _, err = tx.Exec(unFansSQL, time.Now().UTC(), unFollowerId, userId); err != nil {
		logger.Error("update user_fans deleteTime to now error",
			zap.Error(err), zap.Int64("userId", userId),
			zap.Int64("unFollowerId", unFollowerId))
		logging.SetSpanError(span, err)
		return err
	}

	if unFollowerStatus {
		//将取关的人的互关状态设为false
		if _, err := tx.Exec(unBothFollowSQL, unFollowerId, userId); err != nil {
			logger.Error("update user_follows unFollowerStatus to false error",
				zap.Error(err),
				zap.Int64("userId", unFollowerId),
				zap.Int64("followerId", userId))
			logging.SetSpanError(span, err)
			return err
		}
		if _, err = tx.Exec(unBothFansSQL, unFollowerId, userId); err != nil {
			logger.Error("update user_fans unFollowerStatus to false error",
				zap.Error(err),
				zap.Int64("userId", unFollowerId),
				zap.Int64("fansId", userId))
			logging.SetSpanError(span, err)
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		logger.Error("commit unfollow transaction error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return err
	}
	return nil
}
