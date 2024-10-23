package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/proto/collect/collectPb"
	"star/proto/comment/commentPb"
	"star/proto/community/communityPb"
	"star/proto/feed/feedPb"
	"star/proto/like/likePb"
	"star/proto/publish/publishPb"
	"star/proto/relation/relationPb"
	"star/proto/user/userPb"
	"sync"
	"time"

	redis2 "github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go-micro.dev/v4"
	"go.uber.org/zap"
)

// 每页帖子数
const (
	maxNewestPostLength = 400 //GetPostByTime List 的最大长度
	defaultPostCount    = 20
)

type FeedSrv struct {
}

var userService userPb.UserService
var communityService communityPb.CommunityService
var likeService likePb.LikeService
var commentService commentPb.CommentService
var collectService collectPb.CollectService
var relationService relationPb.RelationService
var publishService publishPb.PublishService
var feedSrvIns *FeedSrv

func (p *FeedSrv) New() {

	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	communityMicroService := micro.NewService(micro.Name(str.CommunityServiceClient))
	communityService = communityPb.NewCommunityService(str.CommunityService, communityMicroService.Client())

	likeMicroService := micro.NewService(micro.Name(str.LikeServiceClient))
	likeService = likePb.NewLikeService(str.LikeService, likeMicroService.Client())

	commentMicroService := micro.NewService(micro.Name(str.CommentServiceClient))
	commentService = commentPb.NewCommentService(str.CommentService, commentMicroService.Client())

	collectMicroService := micro.NewService(micro.Name(str.CollectServiceClient))
	collectService = collectPb.NewCollectService(str.CollectService, collectMicroService.Client())

	relationMicroService := micro.NewService(micro.Name(str.RelationServiceClient))
	relationService = relationPb.NewRelationService(str.RelationService, relationMicroService.Client())

	publishMicroService := micro.NewService(micro.Name(str.PublishServiceClient))
	publishService = publishPb.NewPublishService(str.PublishService, publishMicroService.Client())

	cronRunner := cron.New()
	cronRunner.AddFunc("@every 10m", updatePopularPost)
	cronRunner.AddFunc("@hourly", cleanGetPostByTime)
	cronRunner.Start()

}

// QueryPostExist 查询帖子是否存在
func (p *FeedSrv) QueryPostExist(ctx context.Context, req *feedPb.QueryPostExistRequest, resp *feedPb.QueryPostExistResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "QueryPostExistService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.QueryPostExist")

	key := fmt.Sprintf("QueryPostExist:%d", req.PostId)
	_, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		return mysql.QueryPostExist(req.PostId)
	})
	if err != nil {
		if errors.Is(err, str.ErrPostNotExists) {
			resp.Exist = false
			return nil
		}
		logger.Error("query feed exist err:",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	resp.Exist = true
	return nil
}

// GetCommunityPostByNewReply 获取社区最新回复的帖子
func (p *FeedSrv) GetCommunityPostByNewReply(ctx context.Context, req *feedPb.GetCommunityPostByNewReplyRequest, resp *feedPb.GetCommunityPostByNewReplyResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetPostByPopularityService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.GetPostByPopularity")

	key := fmt.Sprintf("GetCommunityPostByNewReply:%d", req.CommunityId)
	length, err := redis.Client.ZCard(ctx, key).Result()
	if err != nil {
		logger.Error("get GetCommunityPostByNewReply list length error ",
			zap.Error(err),
			zap.Int64("actorId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	var posts []*models.Post
	offset := (req.Page - 1) * defaultPostCount
	end := offset + defaultPostCount - 1
	if length == 0 {
		posts, err = mysql.GetCommunityPostByNewReply(req.CommunityId, req.LastReplyTime, defaultPostCount)
		if err != nil {
			logger.Error("mysql GetCommunityPostByNewReply error ",
				zap.Error(err),
				zap.Int64("actorId", req.ActorId))
			logging.SetSpanError(span, err)
			return str.ErrFeedError
		}
		if len(posts) == 0 {
			return nil
		}
		lockKey := fmt.Sprintf("Lock_GetCommunityPostByNewReply:%d", req.CommunityId)
		ok, err := redis.Client.SetNX(ctx, lockKey, 1, 5*time.Second).Result()
		if err != nil {
			logger.Error("get lock error",
				zap.Error(err))
			return str.ErrFeedError
		}
		if ok {
			go func() {
				defer redis.Client.Del(ctx, lockKey)
				_, err := redis.Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
					for _, post := range posts {
						postJson, err := json.Marshal(post)
						if err != nil {
							logger.Error("marshal post error",
								zap.Error(err),
								zap.Any("post", post))
							continue
						}
						pipe.LPush(ctx, key, postJson)
					}
					pipe.Expire(ctx, key, 2*time.Hour)
					return nil
				})
				if err != nil {
					logger.Error("redis save CommunityPostByTime error",
						zap.Error(err))
				}
			}()
		}
	} else if offset < length {
		xPostsJson, err := redis.Client.ZRevRange(ctx, key, offset, end).Result()
		if err != nil {
			logger.Error("redis service error",
				zap.Error(err),
				zap.Int64("actorId", req.ActorId))
			logging.SetSpanError(span, err)
			return str.ErrFeedError
		}
		redis.Client.Expire(ctx, key, 2*time.Hour)
		xPosts := make([]*models.Post, len(xPostsJson))
		for i, xPostJson := range xPostsJson {
			err := json.Unmarshal([]byte(xPostJson), &xPosts[i])
			if err != nil {
				logger.Error("unmarshal  xpostJosn error",
					zap.Error(err),
					zap.Any("xpostJson", xPostJson))
			}
		}
		if len(xPosts) < defaultPostCount {
			yPost, err := mysql.GetCommunityPostByNewReply(req.CommunityId, xPosts[len(xPosts)-1].LastRelyTime, defaultPostCount-len(xPosts))
			if err != nil {
				logger.Error("mysql GetCommunityPostByTime server error",
					zap.Error(err),
					zap.Int64("actorId", req.ActorId),
					zap.Int64("communityId", req.CommunityId))
				logging.SetSpanError(span, err)
				return str.ErrFeedError
			}
			posts = slices.Concat(xPosts, yPost)
		} else {
			posts = xPosts
		}

	} else {
		posts, err = mysql.GetCommunityPostByNewReply(req.CommunityId, req.LastReplyTime, 20)
		if err != nil {
			logger.Error("mysql GetCommunityPostByTime server error",
				zap.Error(err),
				zap.Int64("actorId", req.ActorId),
				zap.Int64("communityId", req.CommunityId))
			logging.SetSpanError(span, err)
			return str.ErrFeedError
		}
	}
	resp.Posts, err = queryDetailed(ctx, posts, req.ActorId, logger)
	if err != nil {
		logger.Error("get feed detail error",
			zap.Error(err),
			zap.Int64("actorId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	resp.NewReplyTime = resp.Posts[len(resp.Posts)-1].LastReplyTime
	return nil
}

func convertGetPostToPB(ctx context.Context, posts []*models.Post, logger *zap.Logger) []*feedPb.Post {
	pposts := make([]*feedPb.Post, len(posts))
	var wg sync.WaitGroup
	postRusultChan := make(chan struct {
		index  int
		pposts *feedPb.Post
	}, len(posts))
	goroutineLimiter := make(chan struct{}, min(15, len(posts)))
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
				logger.Error("get user info error",
					zap.Error(err),
					zap.Int64("userId", post.UserId))
				return
			}
			communityResp, err := communityService.GetCommunityInfo(ctx, &communityPb.GetCommunityInfoRequest{
				CommunityId: post.CommunityId,
			})
			if err != nil {
				logger.Error("get community info error",
					zap.Error(err),
					zap.Int64("communityId", post.CommunityId))
				return
			}
			ppost := &feedPb.Post{
				PostId:    post.PostId,
				Author:    userResp.User,
				Community: communityResp.Community,
				Content:   post.Content,
			}
			postRusultChan <- struct {
				index  int
				pposts *feedPb.Post
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
	ctx, span := tracing.Tracer.Start(context.Background(), "updatePopularPostService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.updatePopularPost")

	key := "GetPostByPopularity"
	posts, err := mysql.GetPostByPopularity(str.DefaultLoadPostNumber, 0, span, logger)
	if err != nil {
		logger.Error("get feed by popularity error",
			zap.Error(err))
		return
	}
	pposts := convertGetPostToPB(context.Background(), posts, logger)
	ppostJson, err := json.Marshal(pposts)
	if err != nil {
		logger.Error("json marshal error",
			zap.Error(err))
		return
	}
	if err = redis.Client.Set(ctx, key, ppostJson, time.Hour).Err(); err != nil {
		logger.Error("redis set popular feed error",
			zap.Error(err))
		return
	}
}

// GetCommunityPostByTime 获取社区最新帖子
func (p *FeedSrv) GetCommunityPostByTime(ctx context.Context, req *feedPb.GetCommunityPostByTimeRequest, resp *feedPb.GetCommunityPostByTimeResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetCommunityPostByTimeService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "CommunityService.GetCommunityPostByTime")

	var posts []*models.Post
	key := fmt.Sprintf("GetCommunityPostByTime:%d", req.CommunityId)
	length, err := redis.Client.LLen(ctx, key).Result()
	if err != nil {
		logger.Error("GetCommunityPostByTime redis error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	offset := (req.Page - 1) * defaultPostCount
	end := offset + defaultPostCount - 1
	if length == 0 {
		posts, err = mysql.GetCommunityPostByTime(req.CommunityId, math.MaxInt64, 6*defaultPostCount)
		if err != nil {
			logger.Error("mysql GetCommunityPostByTime server error",
				zap.Error(err),
				zap.Int64("actorId", req.ActorId),
				zap.Int64("communityId", req.CommunityId))
			logging.SetSpanError(span, err)
			return str.ErrFeedError
		}
		if len(posts) == 0 {
			return nil
		}

		lockKey := fmt.Sprintf("Lock_GetCommunityPostByTime:%d", req.CommunityId)
		ok, err := redis.Client.SetNX(ctx, lockKey, 1, 5*time.Second).Result()
		if err != nil {
			logger.Error("get lock error",
				zap.Error(err))
			return str.ErrFeedError
		}
		if ok {
			go func() {
				defer redis.Client.Del(ctx, lockKey)
				_, err := redis.Client.TxPipelined(ctx, func(pipe redis2.Pipeliner) error {
					for _, post := range posts {
						postJson, err := json.Marshal(post)
						if err != nil {
							logger.Error("marshal post error",
								zap.Error(err),
								zap.Any("post", post))
							continue
						}
						pipe.LPush(ctx, key, postJson)
					}
					pipe.Expire(ctx, key, 2*time.Hour)
					return nil
				})
				if err != nil {
					logger.Error("redis save CommunityPostByTime error",
						zap.Error(err))
				}
			}()
		}
	} else if req.Page < length {

		postsJson, err := redis.Client.LRange(ctx, key, offset, end).Result()
		if err != nil {
			logger.Error("GetCommunityPostByTime redis error",
				zap.Error(err),
				zap.Int64("communityId", req.CommunityId))
			logging.SetSpanError(span, err)
			return str.ErrFeedError
		}
		redis.Client.Expire(ctx, key, 2*time.Hour)
		xPost := make([]*models.Post, len(postsJson))
		for i, postJson := range postsJson {
			err = json.Unmarshal([]byte(postJson), &xPost[i])
			if err != nil {
				logger.Error("unmarshal feed error",
					zap.Error(err))
			}
		}
		if len(postsJson) < defaultPostCount {
			yPost, err := mysql.GetCommunityPostByTime(req.CommunityId, xPost[len(xPost)-1].PostId, defaultPostCount-len(xPost))
			if err != nil {
				logger.Error("mysql GetCommunityPostByTime server error",
					zap.Error(err),
					zap.Int64("actorId", req.ActorId),
					zap.Int64("communityId", req.CommunityId))
				logging.SetSpanError(span, err)
				return str.ErrFeedError
			}
			posts = slices.Concat(xPost, yPost)
		} else {
			posts = xPost
		}

	} else {
		posts, err = mysql.GetCommunityPostByTime(req.CommunityId, req.LastPostId, 20)
		if err != nil {
			logger.Error("mysql GetCommunityPostByTime server error",
				zap.Error(err),
				zap.Int64("actorId", req.ActorId),
				zap.Int64("communityId", req.CommunityId))
			logging.SetSpanError(span, err)
			return str.ErrFeedError
		}
	}
	resp.Posts, err = queryDetailed(ctx, posts, req.ActorId, logger)
	if err != nil {
		logger.Error("get feed detail error",
			zap.Error(err),
			zap.Int64("actorId", req.ActorId),
			zap.Int64("communityId", req.CommunityId))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	resp.NewPostId = resp.Posts[len(resp.Posts)-1].PostId
	return nil
}

// GetPostByRelation 获取关注的人的帖子
func (p *FeedSrv) GetPostByRelation(ctx context.Context, req *feedPb.GetPostByRelationRequest, resp *feedPb.GetPostByRelationResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetPostByTimeService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.GetPostByTime")

	rResp, err := relationService.GetFollowList(ctx, &relationPb.GetFollowListRequest{
		UserId: req.ActorId,
	})
	if err != nil {
		logger.Error("get user follow list error",
			zap.Error(err),
			zap.Int64("userId", req.ActorId))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	followList := rResp.FollowList
	var posts []*feedPb.Post
	var wg sync.WaitGroup
	var lock sync.Mutex
	goroutineLimiter := make(chan struct{}, 15)
	for _, follow := range followList {
		goroutineLimiter <- struct{}{}
		wg.Add(1)
		go func(follow *userPb.User) {
			defer func() {
				wg.Done()
				<-goroutineLimiter
			}()

			pResp, err := publishService.ListPost(ctx, &publishPb.ListPostRequest{
				ActorId: req.ActorId,
				UserId:  follow.UserId,
			})
			if err != nil {
				logger.Error("get follow publish list error",
					zap.Error(err),
					zap.Int64("followId", follow.UserId),
					zap.Int64("ActorId", req.ActorId))
				return
			}
			lock.Lock()
			posts = slices.Concat(posts, pResp.Posts)
			lock.Unlock()
		}(follow)
	}
	wg.Wait()
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateTime > posts[j].CreateTime
	})
	resp.Posts = posts
	return nil
}

// QueryPosts 查询帖子的大致信息
func (p *FeedSrv) QueryPosts(ctx context.Context, req *feedPb.QueryPostsRequest, resp *feedPb.QueryPostsResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "QueryPostsService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.QueryPosts")

	var err error
	resp.Posts, err = query(ctx, req.PostIds, req.ActorId, logger)
	if err != nil {
		logger.Error("query posts error",
			zap.Error(err),
			zap.Any("postIds", req.PostIds))
		logging.SetSpanError(span, err)
		return str.ErrFeedError
	}
	return nil
}

func query(ctx context.Context, postIds []int64, actorId int64, logger *zap.Logger) ([]*feedPb.Post, error) {
	posts, err := mysql.QueryPosts(postIds)
	if err != nil {
		return nil, err
	}
	return queryDetailed(ctx, posts, actorId, logger)
}

func queryDetailed(ctx context.Context, posts []*models.Post, actorId int64, logger *zap.Logger) ([]*feedPb.Post, error) {
	respPosts := make([]*feedPb.Post, len(posts))
	userMap := make(map[int64]*userPb.User)
	communityMap := make(map[int64]*communityPb.Community)
	for i, post := range posts {
		respPosts[i] = &feedPb.Post{
			PostId: post.PostId,
		}
		if _, exist := userMap[post.UserId]; !exist {
			userMap[post.UserId] = &userPb.User{}
		}
		if _, exist := communityMap[post.CommunityId]; !exist {
			communityMap[post.CommunityId] = &communityPb.Community{}
		}
	}
	ucwg := sync.WaitGroup{}
	goroutineLimiter := make(chan struct{}, 45)
	defer close(goroutineLimiter)
	for userId := range userMap {
		ucwg.Add(1)
		goroutineLimiter <- struct{}{}
		go func(userId int64) {
			defer func() {
				<-goroutineLimiter
				ucwg.Done()
			}()
			userResp, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
				UserId: userId,
			})
			if err != nil {
				logger.Error("get user info error",
					zap.Error(err),
					zap.Int64("user_id", userId))
			}
			userMap[userId] = userResp.User
		}(userId)
	}
	for communityId := range communityMap {
		goroutineLimiter <- struct{}{}
		ucwg.Add(1)
		go func(communityId int64) {
			defer func() {
				<-goroutineLimiter
				ucwg.Done()
			}()
			communityResp, err := communityService.GetCommunityInfo(ctx, &communityPb.GetCommunityInfoRequest{
				CommunityId: communityId,
			})
			if err != nil {
				logger.Error("get community info error",
					zap.Error(err),
					zap.Int64("community_id", communityId))
			}
			communityMap[communityId] = communityResp.Community
		}(communityId)
	}
	wg := sync.WaitGroup{}
	for i, post := range posts {

		goroutineLimiter <- struct{}{}
		wg.Add(1)
		//like count
		go func(i int, post *models.Post) {
			defer func() {
				wg.Done()
				<-goroutineLimiter
			}()
			likeCountResp, err := likeService.GetLikeCount(ctx, &likePb.GetLikeCountRequest{
				SourceId:   post.PostId,
				SourceType: 1,
			})
			if err != nil {
				logger.Error("get like count error",
					zap.Error(err),
					zap.Int64("post_id", post.PostId))
				return
			}
			respPosts[i].LikeCount = likeCountResp.Count
		}(i, post)

		goroutineLimiter <- struct{}{}
		wg.Add(1)
		//comment count
		go func(i int, post *models.Post) {
			defer func() {
				wg.Done()
				<-goroutineLimiter
			}()
			commentCountResp, err := commentService.CountComment(ctx, &commentPb.CountCommentRequest{
				ActorId: actorId,
				PostId:  post.PostId,
			})
			if err != nil {
				logger.Error("get comment count error",
					zap.Error(err),
					zap.Int64("post_id", post.PostId),
					zap.Int64("actor_id", actorId))
				return
			}
			respPosts[i].CommentCount = commentCountResp.Count
		}(i, post)
		if actorId != 0 {
			wg.Add(1)
			go func(i int, post *models.Post) {
				defer wg.Done()
				//IsLike
				isLikeResp, err := likeService.IsLike(ctx, &likePb.IsLikeRequest{
					ActorId:    actorId,
					SourceId:   post.PostId,
					SourceType: 1,
				})
				if err != nil {
					logger.Error("get feed isLike error",
						zap.Error(err),
						zap.Int64("post_id", post.PostId),
						zap.Int64("actor_id", actorId))
					return
				}
				respPosts[i].IsLike = isLikeResp.Result
				//IsCollect
				isCollectResp, err := collectService.IsCollect(ctx, &collectPb.IsCollectRequest{
					ActorId: actorId,
					PostId:  post.PostId,
				})
				if err != nil {
					logger.Error("get feed isCollect error",
						zap.Error(err),
						zap.Int64("post_id", post.PostId),
						zap.Int64("actor_id", actorId))
					return
				}
				respPosts[i].IsCollect = isCollectResp.Result
			}(i, post)
		}
	}
	ucwg.Wait()
	wg.Wait()
	return respPosts, nil
}

func cleanGetPostByTime() {
	ctx, span := tracing.Tracer.Start(context.Background(), "cleanGetPostByTimeService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "FeedService.cleanGetPostByTime")

	key := "GetPostByTime"
	length, err := redis.Client.LLen(ctx, key).Result()
	if err != nil {
		logger.Error("get GetPostByTime redis list length error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return
	}
	if length < maxNewestPostLength {
		logger.Info("redis GetPostByTime list length is under control",
			zap.Int64("length", length))
		return
	}
	err = redis.Client.RPopCount(ctx, key, int(length/2)).Err()
	if err != nil {
		logger.Error("rpop GetPostByTime redis list error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return
	}
	logger.Info("clean GetPostByTime list successfully",
		zap.Int64("removeLength", length/2))
}
