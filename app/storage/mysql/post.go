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
	queryPostExistSQL               = "select postId from post where postId=? and deletedAt is null;"
	insertPostSQL                   = "insert into post(postId, userId,collection,star,content,title,isScan,communityId) values(?,?,?,?,?,?,?,?)"
	getPostByPopularitySQL          = "select postId, userId,collection,star,content,isScan,communityId from post where isScan=true order by star desc,collection desc limit ?"
	getCommunityPostByPopularitySQL = "select postId, userId,collection,star,content,isScan,communityId from post where isScan=true and  commnutyId=?  order by star desc,collection desc limit ?"
	getPostByTimeSQL                = "select postId, userId,collection,star,content,isScan,communityId from post where isScan=true order by createdAt desc limit ?"
	queryPostsSQL                   = "select postId, userId,collection,star,content,isScan,communityId from post where postId in (?)"
)

func QueryPostExist(postId int64) (string, error) {
	post := new(models.Post)
	if err := Client.Get(post, queryPostExistSQL, postId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.False, str.ErrPostNotExists
		}
		return str.False, str.ErrPostError
	}
	return str.True, nil
}

func InsertPost(post *models.Post) error {
	if _, err := Client.Exec(insertPostSQL, post.PostId, post.UserId, post.Collection, post.Star, post.Content, post.Title, post.IsScan, post.CommunityId); err != nil {
		utils.Logger.Error("insert post error", zap.Error(err), zap.Any("post", post), zap.Int64("userId", post.UserId))
		return err
	}
	return nil
}

func GetPostByPopularity(limit int, communityId int64) ([]*models.Post, error) {
	var posts []*models.Post
	if communityId == 0 {
		if err := Client.Select(&posts, getPostByPopularitySQL, limit); err != nil {
			utils.Logger.Error("GetPostByPopularity error", zap.Error(err))
			return nil, err
		}
		return posts, nil
	}
	if err := Client.Get(&posts, getCommunityPostByPopularitySQL, communityId, limit); err != nil {
		utils.Logger.Error("GetPostByPopularity error", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

func GetPostByTime(limit int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, getPostByTimeSQL, limit); err != nil {
		utils.Logger.Error("GetPostByTime error", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

func QueryPosts(postIds []int64) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, queryPostsSQL, postIds); err != nil {
		utils.Logger.Error("mysql QueryPosts error", zap.Error(err), zap.Any("postIds", postIds))
		return nil, err
	}
	return posts, nil
}
