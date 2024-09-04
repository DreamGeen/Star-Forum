package httpHandler

import (
	"github.com/gin-gonic/gin"
	"log"
	"star/app/gateway/client"
	"star/constant/str"
	"star/models"
	"star/proto/community/communityPb"
	"star/utils"
)

func CreateCommunityHandler(c *gin.Context) {
	community := new(models.Community)
	if err := c.ShouldBindJSON(community); err != nil {
		log.Println(err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	if err := validateDescriptionAndName(community); err != nil {
		log.Println(err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	userId, err := utils.GetUserId(c)
	if err != nil {
		log.Println(err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	req := &communityPb.CreateCommunityRequest{
		CommunityName: community.CommunityName,
		Description:   community.Description,
		LeaderId:      userId,
	}
	if _, err := client.CreateCommunity(c, req); err != nil {
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
}

//func GetCommunityListHandler(c *gin.Context) {
//	resp, err := client.GetCommunityList()
//	if err != nil {
//		str.Response(c, err, str.Empty, nil)
//		return
//	}
//	str.Response(c, nil, "communitys", resp.Communitys)
//}

//func ShowCommunityHandler(c *gin.Context) {
//	communityIdString := c.Param("id")
//	communityId, _ := strconv.Atoi(communityIdString)
//	resp, err := client.ShowCommunity()
//	if err != nil {
//		str.Response(c, err, str.Empty, nil)
//		return
//	}
//	str.Response(c, nil, str.Empty, resp)
//}

// validateDescription 效验简介字数
func validateDescriptionAndName(community *models.Community) error {
	if len(community.Description) < 2 {
		return str.ErrDescriptionShort
	}
	if len(community.Description) > 50 {
		return str.ErrDescriptionLong
	}
	if len(community.CommunityName) < 1 {
		return str.ErrCommunityNameEmpty
	}
	if len(community.CommunityName) > 10 {
		return str.ErrCommunityNameLong
	}
	return nil
}
