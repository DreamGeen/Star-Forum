package mysql

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"star/constant/str"
	"star/models"
	"star/utils"
)

const (
	queryCommunityByNameSQL = "select communityId from community where communityName=?"
	insertCommunitySQL      = "insert into community(communityId,communityName,description,member,leaderId,img) values (?,?,?,?,?,?)"
	queryCommunityListSQL   = "select communityId,communityName,Img from community "
	getCommunityInfo        = "select communityId, description, communityName, member, leaderId, manageId,img from community  where communityId=?"
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
	if err := Client.Get(community, getCommunityInfo, communityId); err != nil {
		utils.Logger.Error("mysql query community info error", zap.Error(err), zap.Int64("communityId", communityId))
		return nil, str.ErrCommunityError
	}
	return community, nil
}
