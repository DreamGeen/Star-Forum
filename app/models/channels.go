package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type GroupMessage struct {
	Content     string    `json:"content" db:"content" redis:"content"`
	UserName    string    `json:"userName" db:"userName" redis:"userName"`
	Img         string    `json:"img" db:"img" redis:"img"`
	SendTime    time.Time `json:"sendTime" db:"sendTime" redis:"sendTime"`
	ChatId      int64     `json:"chatId" db:"chatId" redis:"chatId"`
	SdUserId    int64     `json:"sdUserId" db:"sdUserId" redis:"sduserid"`
	CommunityId int64     `json:"communityId" db:"communityId" redis:"communityId"`
}

func (g *GroupMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *GroupMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}


func (g *GroupMessage) Value() (driver.Value, error) {
	return []interface{}{g.Content, g.SendTime, g.ChatId, g.SdUserId, g.CommunityId}, nil
}
