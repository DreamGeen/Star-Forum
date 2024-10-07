package service

import (
	"context"
	"errors"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"star/app/storage/redis"
	"star/constant/str"
	"star/proto/collect/collectPb"
	"star/proto/post/postPb"
	"star/utils"
	"strconv"
)

type CollectSrv struct {
}

var postService postPb.PostService

func New() {
	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
	postService = postPb.NewPostService(str.PostService, postMicroService.Client())

}
func (c *CollectSrv) IsCollect(ctx context.Context, req *collectPb.IsCollectRequest, resp *collectPb.IsCollectResponse) error {
	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdStr := fmt.Sprintf("%d", req.PostId)
	ok, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("IsCollect redis service error", zap.Error(err), zap.Int64("post_id", req.PostId), zap.Int64("userId", req.ActorId))
		return err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if ok != 0 {
		resp.Result = true
	} else {
		resp.Result = false
	}
	return nil
}

func (c *CollectSrv) CollectList(ctx context.Context, req *collectPb.CollectListRequest, resp *collectPb.CollectListResponse) error {
	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdsStr, err := redis.Client.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		utils.Logger.Error("redis get user all like posts id error", zap.Error(err), zap.Int64("user_id", req.ActorId))
		return str.ErrLikeError
	}
	if len(postIdsStr) == 0 {
		resp.Posts = nil
		return nil
	}
	postIds := make([]int64, len(postIdsStr))
	for i, postIdStr := range postIdsStr {
		postId, _ := strconv.ParseInt(postIdStr, 10, 64)
		postIds[i] = postId
	}
	queryPostsResp, err := postService.QueryPosts(ctx, &postPb.QueryPostsRequest{
		ActorId: req.ActorId,
		PostIds: postIds,
	})
	if err != nil {
		utils.Logger.Error("query posts detail error", zap.Error(err), zap.Int64("user_id", req.ActorId))
		return str.ErrLikeError
	}
	resp.Posts = queryPostsResp.Posts
	return nil
}

func (c *CollectSrv) CollectAction(ctx context.Context, req *collectPb.CollectActionRequest, resp *collectPb.CollectActionResponse) error {

	key := fmt.Sprintf("user:%d:collect_posts", req.ActorId)
	postIdStr := fmt.Sprintf("%d", req.PostId)
	value, err := redis.Client.ZScore(ctx, key, postIdStr).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("collect_post redis service error", zap.Error(err), zap.Int64("post_id", req.PostId), zap.Int64("userId", req.ActorId))
		return err
	}
	if errors.Is(err, redis2.Nil) {
		err = nil
	}
	if req.ActionType == 1 {
		//收藏
		if value > 0 {
			//重复收藏
			utils.Logger.Warn("user duplicate collect", zap.Int64("post_id", req.PostId), zap.Int64("userId", req.ActorId))
			return nil
		} else {
			//正常收藏
			err = redis.CollectPostAction(ctx, req.ActorId, req.PostId)
			if err != nil {
				utils.Logger.Error("collect_post redis service error", zap.Error(err))
				return str.ErrCollectError
			}
		}
	} else {
		//取消收藏
		if value == 0 {
			//用户未点赞
			utils.Logger.Warn("user did not collect, cancel collecting", zap.Int64("post_id", req.PostId), zap.Int64("userId", req.ActorId))
			return nil
		} else {
			err = redis.UnCollectPostAction(ctx, req.ActorId, req.PostId)
			if err != nil {
				utils.Logger.Error("collect_post redis service error", zap.Error(err))
				return str.ErrCollectError
			}
		}

	}
	return nil
}

func (c *CollectSrv) GetCollectCount(ctx context.Context, req *collectPb.GetCollectCountRequest, resp *collectPb.GetCollectCountResponse) error {
	key := fmt.Sprintf("post:%d:collected_count", req.PostId)
	countStr, err := redis.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get post collect count error", zap.Error(err), zap.Int64("postId", req.PostId))
		return str.ErrCollectError
	}
	if errors.Is(err, redis2.Nil) {
		resp.Count = 0
		return nil
	}
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		utils.Logger.Error("strconv post collect count error", zap.Error(err), zap.Int64("postId", req.PostId), zap.String("countStr", countStr))
		return str.ErrCollectError
	}
	resp.Count = count
	return nil
}
func (c *CollectSrv) GetUserCollectCount(ctx context.Context, req *collectPb.GetUserCollectCountRequest, resp *collectPb.GetUserCollectCountResponse) error  {
	key := fmt.Sprintf("user:%d:collect_posts", req.UserId)
	count,err:=redis.Client.ZCard(ctx, key).Result()
	if err != nil&&!errors.Is(err, redis2.Nil) {
		utils.Logger.Error("redis get user collect count error", zap.Error(err), zap.Int64("user_id", req.UserId))
		return str.ErrCollectError
	}
   if errors.Is(err, redis2.Nil){
	   resp.Count = 0
	   return nil
   }
    resp.Count = count
	return nil
}
