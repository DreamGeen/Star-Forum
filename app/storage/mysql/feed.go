package mysql

import (
	"database/sql"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/models"
	"star/app/utils/logging"
)

const (
	queryPostExistSQL               = "select postId from post where postId=? and deletedAt is null;"
	getPostByPopularitySQL          = "select postId, userId,collection,star,content,isScan,communityId from post where isScan=true order by star desc,collection desc limit ?"
	getCommunityPostByPopularitySQL = "select postId, userId,collection,star,content,isScan,communityId from post where isScan=true and  communityId=?  order by star desc,collection desc limit ?"
	getPostByTimeSQL                = "select postId, userId,collection,star,content,isScan,communityId from post where isScan=true and postId<? limit ?"
	queryPostsSQL                   = "select postId, userId,collection,star,content,isScan,communityId from post where postId in (?)"
	getCommunityPostByTimeSQL       = "select postId, userId,collection,star,content,isScan,communityId  from post where isScan=true and communityId=?and postId<? limit ?"
	getCommunityPostByNewReplySQL   = "select postId,userId,communityId,content from  post where isScan=true and communityId =?  and lastReplyTime <? order by lastReplyTime desc limit ? "
)

func QueryPostExist(postId int64) (string, error) {
	post := new(models.Post)
	if err := Client.Get(post, queryPostExistSQL, postId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.False, str.ErrPostNotExists
		}
		return str.False, str.ErrFeedError
	}
	return str.True, nil
}

func GetPostByPopularity(limit int, communityId int64, span trace.Span, logger *zap.Logger) ([]*models.Post, error) {
	var posts []*models.Post
	if communityId == 0 {
		if err := Client.Select(&posts, getPostByPopularitySQL, limit); err != nil {
			logger.Error("GetPostByPopularity error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return nil, err
		}
		return posts, nil
	}
	if err := Client.Get(&posts, getCommunityPostByPopularitySQL, communityId, limit); err != nil {
		logging.Logger.Error("GetPostByPopularity error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return nil, err
	}
	return posts, nil
}

func GetPostByTime(postId int64, limit int64) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, getPostByTimeSQL, postId, limit); err != nil {
		return nil, err
	}
	return posts, nil
}

func QueryPosts(postIds []int64) ([]*models.Post, error) {
	var posts []*models.Post
	cpostIds := make([]interface{}, len(postIds))
	for i, postId := range postIds {
		cpostIds[i] = postId
	}
	if err := Client.Select(&posts, queryPostsSQL, cpostIds...); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetCommunityPostByTime(communityId int64, lastPostId int64, limit int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, getCommunityPostByTimeSQL, communityId, lastPostId, limit); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetCommunityPostByNewReply(communityId int64, lastReplyTime string, limit int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, getCommunityPostByNewReplySQL, communityId, lastReplyTime, limit); err != nil {
		return nil, err
	}
	return posts, nil
}
