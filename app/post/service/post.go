package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"log"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/proto/community/communityPb"
	"star/proto/post/postPb"
	"star/proto/user/userPb"
	"star/utils"
)

type PostSrv struct {
}

var userService userPb.UserService
var communityService communityPb.CommunityService

func (p *PostSrv) New() {
	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	//创建一个社区微服务客户端
	communityMicroService := micro.NewService(micro.Name(str.CommunityServiceClient))
	communityService = communityPb.NewCommunityService(str.CommunityService, communityMicroService.Client())

	cronRunner := cron.New()
	cronRunner.AddFunc("@every 10m", func() {})
	cronRunner.Start()

}

func (p *PostSrv) QueryPostExist(ctx context.Context, req *postPb.QueryPostExistRequest, resp *postPb.QueryPostExistResponse) error {
	key := fmt.Sprintf("QueryPostExist:%d", req.PostId)
	if _, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		return mysql.QueryPostExist(req.PostId)
	}); err != nil {
		if errors.Is(err, str.ErrPostNotExists) {
			return str.ErrPostNotExists
		}
		log.Println("query post exist err:", err)
		return str.ErrPostError
	}
	return nil
}

func (p *PostSrv) CreatePost(ctx context.Context, req *postPb.CreatePostRequest, resp *postPb.CreatePostResponse) error {
	post := &models.Post{
		PostId:      utils.GetID(),
		UserId:      req.UserId,
		Star:        0,
		Collection:  0,
		Title:       req.Tile,
		Content:     req.Content,
		IsScan:      req.IsScan,
		CommunityId: req.CommunityId,
	}
	if err := mysql.InsertPost(post); err != nil {
		utils.Logger.Error("create post error", zap.Int64("user_id", req.UserId), zap.Error(err))
		return str.ErrPostError
	}
	return nil
}

func (p *PostSrv) GetPostByPopularity(ctx context.Context, req *postPb.GetPostListByPopularityRequest, resp *postPb.GetPostListByPopularityResponse) error {
	posts, err := mysql.GetPostByPopularity(str.DefaultLoadPostNumber)
	if err != nil {
		utils.Logger.Error("get post by popularity error", zap.Error(err))
		return str.ErrPostError
	}

	return nil
}

func convertGetPostToPB(ctx context.Context, posts []*models.Post) []*postPb.Post {
	pposts := make([]*postPb.Post, len(posts))
	for i, post := range posts {
		userResp, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
			UserId: post.UserId,
		})
		if err != nil {
			utils.Logger.Error("get user info error", zap.Error(err))
			return nil
		}
		communityResp, err := communityService.GetCommunityInfo(ctx, &communityPb.GetCommunityInfoRequest{
			CommunityId: post.CommunityId,
		})
		if err != nil {
			utils.Logger.Error("get community info error", zap.Error(err))
			return nil
		}
		ppost := &postPb.Post{
			PostId:        post.PostId,
			CommunityId:   post.CommunityId,
			UserId:        post.UserId,
			UserName:      userResp.User.UserName,
			UserImg:       userResp.User.Img,
			CommunityName: communityResp.Community.CommunityName,
			CommunityImg:  communityResp.Community.CommunityImg,
			Tile:          post.Title,
			Content:       post.Content,
		}
		pposts[i] = ppost
	}
	return pposts
}

func updatePopularPost() {
	posts, err := mysql.GetPostByPopularity(str.DefaultLoadPostNumber)
	if err != nil {
		utils.Logger.Error("get post by popularity error", zap.Error(err))
		return
	}
}

func (p *PostSrv) GetPostByTime(ctx context.Context, req *postPb.GetPostListByTimeRequest, resp *postPb.GetPostListByTimeResponse) error {

	return nil
}
