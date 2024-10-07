package mysql

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"star/constant/str"
	"star/utils"
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
	if err := Client.Select(&followerIdList, getFollowerIdListSQL); err != nil {
		utils.Logger.Error("get user followerId list error", zap.Error(err), zap.Int64("userId", userId))
		return nil, str.ErrRelationError
	}
	return followerIdList, nil
}

func GetFansIdList(userId int64) ([]int64, error) {
	var fansIdList []int64
	if err := Client.Select(&fansIdList, getFansIdListSQL, userId); err != nil {
		utils.Logger.Error("get user fansId list error", zap.Error(err), zap.Int64("userId", userId))
		return nil, str.ErrRelationError
	}
	return fansIdList, nil
}

func GetFollowCount(userId int64) (int64, error) {
	var count int64
	if err := Client.Get(&count, getFollowerCountSQL, userId); err != nil {
		utils.Logger.Error("mysql get user follower count", zap.Error(err), zap.Int64("userId", userId))
		return 0, str.ErrRelationError
	}
	return count, nil
}
func GetFansCount(userId int64) (int64, error) {
	var count int64
	if err := Client.Get(&count, getFansCountSQL, userId); err != nil {
		utils.Logger.Error("mysql get user fans count", zap.Error(err), zap.Int64("userId", userId))
		return 0, str.ErrRelationError
	}
	return count, nil
}

func Follow(userId, beFollowerId int64, beFollowerStatus bool) error {
	tx, err := Client.Beginx()
	if err != nil {
		utils.Logger.Error("start follow transaction error", zap.Error(err))
		return str.ErrRelationError
	}
	defer func() {
		if p := recover(); p != nil {
			utils.Logger.Error(" follow recovered from panic, transaction rolled back:", zap.Any("panic", p))
			tx.Rollback()
			err = str.ErrRelationError
		} else if err != nil {
			utils.Logger.Error(" follow transaction rolled back due to error:", zap.Error(err))
			tx.Rollback()
		}
	}()
	//检查follow记录是否存在
	var count int
	if err := tx.Get(&count, checkFollowExistSQL, userId, beFollowerId); err != nil {
		utils.Logger.Error("check user follow exist error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("beFollowerId", beFollowerId))
		return err
	}
	if count > 0 {
		if _, err = tx.Exec(followExistSQL, beFollowerStatus, userId, beFollowerId); err != nil {
			utils.Logger.Error("update  user_follows deleteTime to null error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("beFollowerId", beFollowerId))
			return err
		}
		if _, err = tx.Exec(fansExistSQL, beFollowerStatus, beFollowerId, userId); err != nil {
			utils.Logger.Error("update user_fans deleteTime to null error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("beFollowerId", beFollowerId))
			return err
		}
	} else {
		if _, err = tx.Exec(followUnExistSQL, userId, beFollowerId, beFollowerStatus); err != nil {
			utils.Logger.Error("insert follow error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("beFollowerId", beFollowerId))
			return err
		}
		if _, err = tx.Exec(fansUnExist, beFollowerId, userId, beFollowerStatus); err != nil {
			utils.Logger.Error("insert fans error ", zap.Error(err), zap.Int64("userId", beFollowerId), zap.Int64("fansId", userId))
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		utils.Logger.Error("commit follow transaction error", zap.Error(err))
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
		utils.Logger.Error("check  be_follower follow exist error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("unFollowerId", unFollowerId))
		return "0", err
	}
	countStr := strconv.FormatInt(count, 10)
	return countStr, nil
}

func Unfollow(userId, unFollowerId int64, unFollowerStatus bool) error {
	tx, err := Client.Beginx()
	if err != nil {
		utils.Logger.Error("start unfollow transaction error", zap.Error(err))
		return str.ErrRelationError
	}
	defer func() {
		if p := recover(); p != nil {
			utils.Logger.Error(" unfollow recovered from panic, transaction rolled back:", zap.Any("panic", p))
			tx.Rollback()
			err = str.ErrRelationError
		} else if err != nil {
			utils.Logger.Error(" unfollow transaction rolled back due to error:", zap.Error(err))
			tx.Rollback()
		}
	}()
	if _, err = tx.Exec(unFollowSQL, time.Now().UTC(), userId, unFollowerId); err != nil {
		utils.Logger.Error("update user_follows deleteTime to now error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("unFollowerId", unFollowerId))
		return err
	}
	if _, err = tx.Exec(unFansSQL, time.Now().UTC(), unFollowerId, userId); err != nil {
		utils.Logger.Error("update user_fans deleteTime to now error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("unFollowerId", unFollowerId))
		return err
	}

	if unFollowerStatus {
		//将取关的人的互关状态设为false
		if _, err := tx.Exec(unBothFollowSQL, unFollowerId, userId); err != nil {
			utils.Logger.Error("update user_follows unFollowerStatus to false error", zap.Error(err), zap.Int64("userId", unFollowerId), zap.Int64("followerId", userId))
			return err
		}
		if _, err = tx.Exec(unBothFansSQL, unFollowerId, userId); err != nil {
			utils.Logger.Error("update user_fans unFollowerStatus to false error", zap.Error(err), zap.Int64("userId", unFollowerId), zap.Int64("fansId", userId))
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		utils.Logger.Error("commit unfollow transaction error", zap.Error(err))
		return err
	}
	return nil
}
