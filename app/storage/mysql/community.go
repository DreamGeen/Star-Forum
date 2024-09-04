package mysql

import (
	"database/sql"
	"errors"
	"star/constant/str"
	"star/models"
)

const (
	queryCommunityByNameSQL = "select communityId from community where communityName=?"
	insertCommunitySQL      = "insert into community(communityId,communityName,description,member,leaderId,img) values (?,?,?,?,?,?)"
	queryCommunityListSQL   = "select communityId,communityName,Img from community "
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
