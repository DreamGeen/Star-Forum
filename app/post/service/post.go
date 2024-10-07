package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"log"
	"star/app/comment/dao/redis"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/proto/comment/commentPb"
	"star/proto/community/communityPb"
	"star/proto/like/likePb"
	"star/proto/post/postPb"
	"star/proto/user/userPb"
	"star/utils"
	"sync"
	"time"
)

type PostSrv struct {
}

var userService userPb.UserService
var communityService communityPb.CommunityService
var likeService likePb.LikeService
var commentService commentPb.CommentService

func (p *PostSrv) New() {

	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	communityMicroService := micro.NewService(micro.Name(str.CommunityServiceClient))
	communityService = communityPb.NewCommunityService(str.CommunityService, communityMicroService.Client())

	likeMicroService := micro.NewService(micro.Name(str.LikeServiceClient))
	likeService = likePb.NewLikeService(str.LikeService, likeMicroService.Client())

	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

	cronRunner := cron.New()
	cronRunner.AddFunc("@every 10m", updatePopularPost)
	cronRunner.AddFunc("@hourly", cleanPost)
	cronRunner.Start()

}

func (p *PostSrv) QueryPostExist(ctx context.Context, req *postPb.QueryPostExistRequest, resp *postPb.QueryPostExistResponse) error {
	key := fmt.Sprintf("QueryPostExist:%d", req.PostId)
	_, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		return mysql.QueryPostExist(req.PostId)
	})
	if err != nil {
		if errors.Is(err, str.ErrPostNotExists) {
			resp.Exist = false
			return nil
		}
		log.Println("query post exist err:", err)
		return str.ErrPostError
	}
	resp.Exist = true
	return nil
}

func (p *PostSrv) CreatePost(ctx context.Context, req *postPb.CreatePostRequest, resp *postPb.CreatePostResponse) error {
	post := &models.Post{
		PostId:      utils.GetID(),
		UserId:      req.UserId,
		Star:        0,
		Collection:  0,
		Content:     req.Content,
		IsScan:      req.IsScan,
		CommunityId: req.CommunityId,
	}
	key := fmt.Sprintf("GetPostByTime:%d", req.CommunityId)
	if err := redis.Client.LPush(ctx, key, post).Err(); err != nil {
		utils.Logger.Error("create post error", zap.Int64("user_id", req.UserId), zap.Error(err))
		return str.ErrPostError
	}
	if err := mysql.InsertPost(post); err != nil {
		utils.Logger.Error("create post error", zap.Int64("user_id", req.UserId), zap.Error(err))
		return str.ErrPostError
	}
	return nil
}

func (p *PostSrv) GetPostByPopularity(ctx context.Context, req *postPb.GetPostListByPopularityRequest, resp *postPb.GetPostListByPopularityResponse) error {
	key := fmt.Sprintf("GetPostByPopularity:%d", req.CommunityId)
	val, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("get post by popularity error", zap.Error(err))
		return str.ErrPostError
	}
	if errors.Is(err, redis2.Nil) {
		posts, err := mysql.GetPostByPopularity(str.DefaultLoadPostNumber, req.CommunityId)
		if err != nil {
			utils.Logger.Error("get post by popularity error", zap.Error(err))
			return str.ErrPostError
		}
		pposts := convertGetPostToPB(ctx, posts)
		resp.Posts = pposts
		ppostJson, err := json.Marshal(pposts)
		if err != nil {
			utils.Logger.Error("json marshal error", zap.Error(err))
			return nil
		}
		if err = redis.Client.Set(ctx, key, ppostJson, time.Hour).Err(); err != nil {
			utils.Logger.Error("redis set error", zap.Error(err))
			return nil
		}
		return nil
	}
	var posts []*postPb.Post
	if err = json.Unmarshal([]byte(val), &posts); err != nil {
		utils.Logger.Error("json unmarshal error", zap.Error(err))
		return str.ErrPostError
	}
	resp.Posts = posts
	return nil
}

func convertGetPostToPB(ctx context.Context, posts []*models.Post) []*postPb.Post {
	pposts := make([]*postPb.Post, len(posts))
	var wg sync.WaitGroup
	postRusultChan := make(chan struct {
		index  int
		pposts *postPb.Post
	}, len(posts))
	goroutineLimiter := make(chan struct{}, 15)
	for i, post := range posts {
		wg.Add(1)
		goroutineLimiter <- struct{}{}
		go func(i int, post *models.Post) {
			defer func() {
				<-goroutineLimiter
				wg.Done()
			}()
			userResp, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
				UserId: post.UserId,
			})
			if err != nil {
				utils.Logger.Error("get user info error", zap.Error(err))
				return
			}
			communityResp, err := communityService.GetCommunityInfo(ctx, &communityPb.GetCommunityInfoRequest{
				CommunityId: post.CommunityId,
			})
			if err != nil {
				utils.Logger.Error("get community info error", zap.Error(err))
				return
			}
			ppost := &postPb.Post{
				PostId:    post.PostId,
				Author:    userResp.User,
				Community: communityResp.Community,
				Content:   post.Content,
			}
			postRusultChan <- struct {
				index  int
				pposts *postPb.Post
			}{index: i, pposts: ppost}
		}(i, post)
	}
	go func() {
		wg.Wait()
		close(postRusultChan)
		close(goroutineLimiter)
	}()
	for post := range postRusultChan {
		pposts[post.index] = post.pposts
	}
	return pposts
}

func updatePopularPost() {
	key := "GetPostByPopularity"
	posts, err := mysql.GetPostByPopularity(str.DefaultLoadPostNumber, 0)
	if err != nil {
		utils.Logger.Error("get post by popularity error", zap.Error(err))
		return
	}
	pposts := convertGetPostToPB(context.Background(), posts)
	ppostJson, err := json.Marshal(pposts)
	if err != nil {
		utils.Logger.Error("json marshal error", zap.Error(err))
		return
	}
	if err = redis.Client.Set(context.Background(), key, ppostJson, time.Hour).Err(); err != nil {
		utils.Logger.Error("redis set error", zap.Error(err))
		return
	}
}
func (p *PostSrv) GetPostByTime(ctx context.Context, req *postPb.GetPostListByTimeRequest, resp *postPb.GetPostListByTimeResponse) error {
	limit := req.Limit
	page := req.Page
	offset := limit * (page - 1)

	key := fmt.Sprintf("GetPostByTime:%d", req.CommunityId)
	var post []*models.Post
	if err := redis.Client.LRange(ctx, key, offset, offset+limit).ScanSlice(&post); err != nil {
		utils.Logger.Error("get redis list post error", zap.Error(err), zap.Int64("communityId", req.CommunityId))
		return str.ErrPostError
	}
	if len(post) != 0 {
		resp.Posts = convertGetPostToPB(ctx, post)
		return nil
	}
	return nil
}

func (p *PostSrv) QueryPosts(ctx context.Context, req *postPb.QueryPostsRequest, resp *postPb.QueryPostsResponse) error {
	var err error
	resp.Posts, err = query(ctx, req.PostIds, req.ActorId)
	if err != nil {
		utils.Logger.Error("query posts error", zap.Error(err), zap.Any("postIds", req.PostIds))
		return str.ErrPostError
	}
	return nil
}

func query(ctx context.Context, postIds []int64, actorId int64) ([]*postPb.Post, error) {
	posts, err := mysql.QueryPosts(postIds)
	if err != nil {
		return nil, err
	}
	return queryDetailed(ctx, posts, actorId)
}

func queryDetailed(ctx context.Context, posts []*models.Post, actorId int64) ([]*postPb.Post, error) {
	respPosts := make([]*postPb.Post, len(posts))
	userMap := make(map[int64]*userPb.User)
	communityMap := make(map[int64]*communityPb.Community)
	for i, post := range posts {
		respPosts[i] = &postPb.Post{
			PostId: post.PostId,
		}
		if _, exist := userMap[post.UserId]; !exist {
			userMap[post.UserId] = &userPb.User{}
		}
		if _, exist := communityMap[post.CommunityId]; !exist {
			communityMap[post.CommunityId] = &communityPb.Community{}
		}
	}
	wg := &sync.WaitGroup{}
	goroutineLimiter := make(chan struct{}, 30)
	defer close(goroutineLimiter)
	for userId := range userMap {
		wg.Add(1)
		goroutineLimiter <- struct{}{}
		go func(userId int64) {
			defer func() {
				<-goroutineLimiter
				wg.Done()
			}()
			userResp, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
				UserId: userId,
			})
			if err != nil {
				utils.Logger.Error("get user info error", zap.Error(err), zap.Int64("user_id", userId))
			}
			userMap[userId] = userResp.User
		}(userId)
	}
	for communityId := range communityMap {
		goroutineLimiter <- struct{}{}
		wg.Add(1)
		go func(communityId int64) {
			defer func() {
				<-goroutineLimiter
				wg.Done()
			}()
			communityResp, err := communityService.GetCommunityInfo(ctx, &communityPb.GetCommunityInfoRequest{
				CommunityId: communityId,
			})
			if err != nil {
				utils.Logger.Error("get community info error", zap.Error(err), zap.Int64("community_id", communityId))
			}
			communityMap[communityId] = communityResp.Community
		}(communityId)
	}
	wg.Wait()
	for i, post := range posts {
		wg.Add(2)
		//like count
		go func(i int, post *models.Post) {
			defer wg.Done()
			likeCountResp, err := likeService.GetLikeCount(ctx, &likePb.GetLikeCountRequest{
				SourceId:   post.PostId,
				SourceType: 1,
			})
			if err != nil {
				utils.Logger.Error("get like count error", zap.Error(err), zap.Int64("post_id", post.PostId))
				return
			}
			respPosts[i].LikeCount = likeCountResp.Count
		}(i, post)

		//comment count
		go func(i int, post *models.Post) {
			defer wg.Done()
			commentCountResp, err := commentService.CountComment(ctx, &commentPb.CountCommentRequest{
				ActorId: actorId,
				PostId:  post.PostId,
			})
			if err != nil {
				utils.Logger.Error("get comment count error", zap.Error(err), zap.Int64("post_id", post.PostId), zap.Int64("actor_id", actorId))
				return
			}
			respPosts[i].CommentCount = commentCountResp.Count
		}(i, post)
		if actorId != 0 {
			wg.Add(1)
			//IsLike
			go func(i int, post *models.Post) {
				defer wg.Done()
				isLikeResp, err := likeService.IsLike(ctx, &likePb.IsLikeRequest{
					ActorId:    actorId,
					SourceId:   post.PostId,
					SourceType: 1,
				})
				if err != nil {
					utils.Logger.Error("get post isLike error", zap.Error(err), zap.Int64("post_id", post.PostId), zap.Int64("actor_id", actorId))
					return
				}
				respPosts[i].IsLike = isLikeResp.Result
			}(i, post)
			//IsCollect
		}
	}
	wg.Wait()
	return respPosts, nil
}

func cleanPost() {
	communityIds, err := mysql.GetAllCommunityId()
	if err != nil {
		utils.Logger.Error("get all community id error", zap.Error(err))
		return
	}
	pipe := redis.Client.Pipeline()
	for _, communityId := range communityIds {
		key := fmt.Sprintf("GetPostByTime:%d", communityId)
		pipe.LLen(context.Background(), key)
	}
	cmder, err := pipe.Exec(context.Background())
	if err != nil {
		utils.Logger.Error("get all list length error", zap.Error(err))
		return
	}
	for i, cmd := range cmder {
		communityId := communityIds[i]
		key := fmt.Sprintf("GetPostByTime:%d", communityId)

		length := cmd.(*redis2.IntCmd).Val()
		if length < 200 {
			continue
		}
		if err := redis.Client.LTrim(context.Background(), key, 0, 199).Err(); err != nil {
			utils.Logger.Error("delete redis post error", zap.Error(err), zap.Int64("communityId", communityId))
			continue
		}

	}
}
