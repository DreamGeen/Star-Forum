package models

type Community struct {
	CommunityId   int64  `json:"community_id" db:"communityId"`
	LeaderId      int64  `json:"leader_id" db:"leaderId"`
	LastMsgId     int64  `json:"last_msg_id" db:"lastMsgId"`
	Member        int64  `json:"member" db:"member"`
	Description   string `json:"description" db:"description"`
	CommunityName string `json:"community_name" db:"communityName"`
	Img           string `json:"img" db:"img"`
}
