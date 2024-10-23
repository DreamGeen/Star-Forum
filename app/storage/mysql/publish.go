package mysql

import (
	"star/app/models"
	"strconv"
)

const (
	insertPostSQL = "insert into post(postId, userId,collection,star,content,isScan,communityId) values(?,?,?,?,?,?,?)"
	countPostSQL  = "select count(1) from post where userId=?"
	listPostSQL   = "select postId,userId,content,isScan,communityId from post where deletedAt is NULL and userId=? order by createdAt desc"
)

func InsertPost(post *models.Post) error {
	if _, err := Client.Exec(insertPostSQL, post.PostId, post.UserId, post.Collection, post.Star, post.Content, post.IsScan, post.CommunityId); err != nil {
		return err
	}
	return nil
}

func CountPost(userId int64) (string, error) {
	var count int64
	if err := Client.Get(&count, countPostSQL, userId); err != nil {
		return "", err
	}
	return strconv.FormatInt(count, 64), nil
}

func ListPost(userId int64) ([]*models.Post, error) {
	var posts []*models.Post
	if err := Client.Select(&posts, listPostSQL, userId); err != nil {
		return nil, err
	}
	return posts, nil
}
