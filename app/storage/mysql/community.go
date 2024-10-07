package mysql

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"star/constant/str"
	"star/models"
	"star/utils"
	"strconv"
	"time"
)

const (
	queryCommunityByNameSQL      = "select communityId from community where communityName=?"
	insertCommunitySQL           = "insert into community(communityId,communityName,description,member,leaderId,img) values (?,?,?,?,?,?)"
	queryCommunityListSQL        = "select communityId,communityName,Img from community "
	getCommunityInfoSQL          = "select communityId, description, communityName, member, leaderId, manageId,img from community  where communityId=?"
	getAllCommunityIdSQL         = "select communityId from community"
	countCommunityFollowSQL      = "select count(1) from community_follows where userId=?"
	isFollowCommunitySQL         = "select count(1) from community_follows where userId=? and communityId=? and  deletedAt IS NULL"
	getCommunityFollowIdSQL      = "select  communityId from community_follows where userId=? and deletedAt IS NULL"
	checkCommunityFollowExistSQL = "select count(1) from community_follows where user_id=? and communityId=? and  deletedAt IS NOT NULL "
	followCommunityExistSQL      = "update community_follows set deletedAt=null where userId=? and communityId=?"
	followCommunityUnExistSQL    = "insert  into community_follows(userId,communityId) values (?,?)"
	unFollowCommunitySQL         = "update community_follows set deletedAt=? where userId=? and communityId=?"
)

func CheckCommunity(communityName string) error {
	community := new(models.Community)
	if err := Client.Get(community, queryCommunityByNameSQL, communityName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.ErrCommunityNameExists
		}
		return str.ErrCommunityError
	}
	return nil
}

func QueryCommunityList() error {
	var communityList []*models.Community
	if err := Client.Select(&communityList, queryCommunityListSQL); err != nil {
		return str.ErrCommunityError
	}
	return nil
}

func InsertCommunity(community *models.Community) error {
	if _, err := Client.Exec(insertCommunitySQL,
		community.CommunityId,
		community.CommunityName,
		community.Description,
		community.Member,
		community.LeaderId,
		community.Img,
	); err != nil {
		return str.ErrCommentError
	}
	return nil
}

func GetCommunityInfo(communityId int64) (*models.Community, error) {
	community := new(models.Community)
	if err := Client.Get(community, getCommunityInfoSQL, communityId); err != nil {
		utils.Logger.Error("mysql query community info error", zap.Error(err), zap.Int64("communityId", communityId))
		return nil, err
	}
	return community, nil
}

func GetAllCommunityId() ([]int64, error) {
	var commnutyIds []int64
	if err := Client.Select(&commnutyIds, getAllCommunityIdSQL); err != nil {
		utils.Logger.Error("mysql get all community id error", zap.Error(err))
		return nil, err
	}
	return commnutyIds, nil
}

func CountCommunityFollow(userId int64) (int64, error) {
	var count int64
	if err := Client.Get(&count, countCommunityFollowSQL, userId); err != nil {
		utils.Logger.Error("mysql count user community follow error", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func IsFollowCommunity(userId, communityId int64) (string, error) {
	var count int64
	if err := Client.Get(&count, isFollowCommunitySQL, userId, communityId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "0", nil
		}
		utils.Logger.Error("mysql check user is follow community error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("communityId", communityId))
		return "0", err
	}
	countStr := strconv.FormatInt(count, 10)
	return countStr, nil
}

func GetCommunityFollowId(userId int64) ([]int64, error) {
	var communityIds []int64
	if err := Client.Select(&communityIds, getCommunityFollowIdSQL, userId); err != nil {
		utils.Logger.Error("mysql get all community follow id error", zap.Error(err))
		return nil, err
	}
	return communityIds, nil
}

func FollowCommunity(userId, communityId int64) error {
	tx, err := Client.Beginx()
	if err != nil {
		utils.Logger.Error("start community follow transaction error", zap.Error(err))
		return str.ErrRelationError
	}
	defer func() {
		if p := recover(); p != nil {
			utils.Logger.Error(" community follow recovered from panic, transaction rolled back:", zap.Any("panic", p))
			tx.Rollback()
			err = str.ErrRelationError
		} else if err != nil {
			utils.Logger.Error(" community transaction rolled back due to error:", zap.Error(err))
			tx.Rollback()
		}
	}()
	//检查follow记录是否存在
	var count int
	if err := tx.Get(&count, checkCommunityFollowExistSQL, userId, communityId); err != nil {
		utils.Logger.Error("check user community follow exist error", zap.Error(err), zap.Int64("userId", userId), zap.Int64("beFollowerId", beFollowerId))
		return err
	}
	if count > 0 {
		if _, err = tx.Exec(followCommunityExistSQL, userId, communityId); err != nil {
			utils.Logger.Error("update  community_follows deleteTime to null error",
				zap.Error(err), zap.Int64("userId", userId),
				zap.Int64("community", communityId))
			return err
		}
	} else {
		if _, err = tx.Exec(followCommunityUnExistSQL, userId, communityId); err != nil {
			utils.Logger.Error("insert community_follows error", zap.Error(err),
				zap.Int64("userId", userId),
				zap.Int64("communityId", communityId))
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		utils.Logger.Error("commit follow transaction error", zap.Error(err))
		return err
	}
	return nil
}

func UnFollowCommunity(userId, communityId int64) error {
	if _, err := Client.Exec(unFollowCommunitySQL, time.Now().UTC(), userId, communityId); err != nil {
		utils.Logger.Error("update  community_follows deleteTime to null error",
			zap.Error(err), zap.Int64("userId", userId),
			zap.Int64("community", communityId))
		return err
	}
	return nil
}
