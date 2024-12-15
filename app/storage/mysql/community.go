package mysql

import (
	"database/sql"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/models"
	"star/app/utils/logging"
	"strconv"
	"time"
)

const (
	queryCommunityByNameSQL      = "select communityId from community where communityName=?"
	insertCommunitySQL           = "insert into community(communityId,communityName,description,member,leaderId,img) values (?,?,?,?,?,?)"
	queryCommunityListSQL        = "select communityId,communityName,Avatar from community "
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
			return str.ErrCommunityNotExists
		}
		return err
	}
	return nil
}

func QueryCommunityList() error {
	var communityList []*models.Community
	if err := Client.Select(&communityList, queryCommunityListSQL); err != nil {
		return err
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
		return err
	}
	return nil
}

func GetCommunityInfo(communityId int64) (*models.Community, error) {
	community := new(models.Community)
	if err := Client.Get(community, getCommunityInfoSQL, communityId); err != nil {
		return nil, err
	}
	return community, nil
}

func GetAllCommunityId() ([]int64, error) {
	var commnutyIds []int64
	if err := Client.Select(&commnutyIds, getAllCommunityIdSQL); err != nil {
		return nil, err
	}
	return commnutyIds, nil
}

func CountCommunityFollow(userId int64) (int64, error) {
	var count int64
	if err := Client.Get(&count, countCommunityFollowSQL, userId); err != nil {
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
		return "0", err
	}
	countStr := strconv.FormatInt(count, 10)
	return countStr, nil
}

func GetCommunityFollowId(userId int64) ([]int64, error) {
	var communityIds []int64
	if err := Client.Select(&communityIds, getCommunityFollowIdSQL, userId); err != nil {
		return nil, err
	}
	return communityIds, nil
}

func FollowCommunity(userId, communityId int64, span trace.Span, logger *zap.Logger) error {
	tx, err := Client.Beginx()
	if err != nil {
		logger.Error("start community follow transaction error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			logger.Error(" community follow recovered from panic, transaction rolled back:",
				zap.Any("panic", p))
			logging.SetSpanError(span, err)
			tx.Rollback()
			err = str.ErrCommunityError
		} else if err != nil {
			logger.Error(" community transaction rolled back due to error:",
				zap.Error(err))
			logging.SetSpanError(span, err)
			tx.Rollback()
		}
	}()
	//检查follow记录是否存在
	var count int
	if err := tx.Get(&count, checkCommunityFollowExistSQL, userId, communityId); err != nil {
		logger.Error("check user community follow exist error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("beFollowerId", communityId))
		logging.SetSpanError(span, err)
		return err
	}
	if count > 0 {
		if _, err = tx.Exec(followCommunityExistSQL, userId, communityId); err != nil {
			logger.Error("update  community_follows deleteTime to null error",
				zap.Error(err),
				zap.Int64("userId", userId),
				zap.Int64("community", communityId))
			logging.SetSpanError(span, err)
			return err
		}
	} else {
		if _, err = tx.Exec(followCommunityUnExistSQL, userId, communityId); err != nil {
			logger.Error("insert community_follows error",
				zap.Error(err),
				zap.Int64("userId", userId),
				zap.Int64("communityId", communityId))
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

func UnFollowCommunity(userId, communityId int64) error {
	if _, err := Client.Exec(unFollowCommunitySQL, time.Now().UTC(), userId, communityId); err != nil {
		return err
	}
	return nil
}
