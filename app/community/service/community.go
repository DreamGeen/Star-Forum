package service

import (
	"context"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/proto/community/communityPb"
	"star/utils"
)

type CommunitySrv struct {
}

func (c *CommunitySrv) CreateCommunity(ctx context.Context, req *communityPb.CreateCommunityRequest, resp *communityPb.EmptyCommunityResponse) error {
	//检查该社区名是否已经存在
	if err := mysql.CheckCommunity(req.CommunityName); err != nil {
		return err
	}
	//构建社区结构体
	community := &models.Community{
		CommunityId:   utils.GetID(),
		CommunityName: req.CommunityName,
		Description:   req.Description,
		LeaderId:      req.LeaderId,
		Img:           str.DefaultCommunityImg,
		Member:        1,
	}
	//将社区插入mysql
	if err := mysql.InsertCommunity(community); err != nil {
		return err
	}
	return nil
}

func (c *CommunitySrv) GetCommunityList(ctx context.Context, req *communityPb.EmptyCommunityRequest, resp *communityPb.GetCommunityListResponse) error {
	//查询community列表

	return nil
}

func (c *CommunitySrv) ShowCommunity(ctx context.Context, req *communityPb.ShowCommunityRequest, resp *communityPb.ShowCommunityResponse) error {
	return nil
}
