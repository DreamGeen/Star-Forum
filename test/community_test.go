package test

import (
	"context"
	"star/app/gateway/client"
	"star/proto/community/communityPb"
	"testing"
)

func TestCreateCommunity(t *testing.T) {
	client.Init()
	_, err := client.CreateCommunity(context.Background(), &communityPb.CreateCommunityRequest{
		CommunityName: "天天",
		Description:   "测试社区",
		LeaderId:      1820019310731464704,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
