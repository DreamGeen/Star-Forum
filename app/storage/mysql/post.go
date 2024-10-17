package mysql

import (
	"database/sql"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/models"
	"star/app/utils/logging"
)

const (
	queryPostExistSQL               = "select postId from feed where postId=? and deletedAt is null;"
	insertPostSQL                   = "insert into feed(postId, userId,collection,star,content,title,isScan,communityId) values(?,?,?,?,?,?,?,?)"
	getPostByPopularitySQL          = "select postId, userId,collection,star,content,isScan,communityId from feed where isScan=true order by star desc,collection desc limit ?"
	getCommunityPostByPopularitySQL = "select postId, userId,collection,star,content,isScan,communityId from feed where isScan=true and  commnutyId=?  order by star desc,collection desc limit ?"
	getPostByTimeSQL                = "select postId, userId,collection,star,content,isScan,communityId from feed where isScan=true and postId<? limit ?"
	queryPostsSQL                   = "select postId, userId,collection,star,content,isScan,communityId from feed where postId in (?)"
	getCommunityPostByTimeSQL       = "select postId, userId,collection,star,content,isScan,communityId  from feed where isScan=true and communityId=?and postId<? limit ?"
)

func QueryPostExist(postId int64) (string, error) {
	post := new(models.Post)
	if err := Client.Get(post, queryPostExistSQL, postId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str2.False, str2.ErrPostNotExists
		}
		return str2.False, str2.ErrPostError
	}
	return str2.True, nil
}

func InsertPost(post *models.Post) error {
	if _, err := Client.Exec(insertPostSQL, post.PostId, post.UserId, post.Collection, post.Star, post.Content, post.IsScan, post.CommunityId); err != nil {
		return err
	}
	return nil
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
	if err := Client.Select(&posts, queryPostsSQL, postIds); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetCommunityPostByTime(communityId int64, lastPostId int64, limit int64) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, getCommunityPostByTimeSQL, communityId, lastPostId, limit); err != nil {
		return nil, err
	}
	return posts, nil
}
