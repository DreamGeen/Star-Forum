package models

import (
	"database/sql/driver"
	"time"
)

type Counts struct {
	MentionCount    uint32 ` db:"mentionCount" json:"mention_count"`
	LikeCount       uint32 `db:"likeCount" json:"like_count"`
	ReplyCount      uint32 `db:"replyCount" json:"reply_count"`
	SystemCount     uint32 `db:"systemCount" json:"system_count"`
	PrivateMsgCount uint32 `db:"privateMsgCount" json:"private_msg_count"`
	TotalCount      uint32 `json:"total_count"`
}
type PrivateMessage struct {
	Id            int64     `db:"private_message_id" json:"private_msg_id"`
	SenderId      int64     `db:"sender_id" json:"sender_id"`
	RecipientId   int64     `db:"recipient_id" json:"receiver_id"`
	Content       string    `db:"content" json:"content"`
	Status        bool      `db:"status" json:"status"`
	SendTime      time.Time `db:"send_time" json:"send_time"`
	PrivateChatId int64     `db:"private_chat_id" json:"private_chat_id"`
}

type PrivateChat struct {
	Id             int64     `db:"private_chat_id" json:"private_chat_id"`
	User1Id        int64     `db:"user1_id" json:"user1_id"`
	User2Id        int64     `db:"user2_id" json:"user2_id"`
	LastMsgContent string    `db:"last_message_content" json:"last_message_content"`
	LastSendTime   time.Time `db:"last_message_time" json:"last_message_time"`
}

type SystemMessage struct {
	Id          int64     `db:"system_notice_id" json:"system_notice_id"`
	RecipientId int64     `db:"recipient_id" json:"recipient_id"`
	ManagerId   int64     `db:"manager_id" json:"manager_id"`
	Type        string    `db:"type" json:"type"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	Status      bool      `db:"status" json:"status"`
	PublishTime time.Time `db:"publish_time" json:"publish_time"`
}

type SystemMessageUser struct {
	Id              int64     `db:"user_notice_id" json:"user_notice_id"`
	SystemMessageId int64     `db:"system_notice_id" json:"system_notice_id"`
	RecipientId     int64     `db:"recipient_id" json:"recipient_id"`
	Status          bool      `db:"status" json:"status"`
	PullTime        time.Time `db:"pull_time" json:"pull_time"`
}

func (s *SystemMessageUser) Value() (driver.Value, error) {
	return []interface{}{s.Id, s.SystemMessageId, s.RecipientId, s.Status}, nil
}

type RemindMessage struct {
	Id          int64     `db:"id" json:"remind_message_id"`
	SourceId    int64     `db:"source_id" json:"source_id"`
	SenderId    int64     `db:"sender_id" json:"sender_id"`
	RecipientId int64     `db:"recipient_id" json:"recipient_id"`
	SourceType  string    `db:"source_type" json:"source_type"`
	Content     string    `db:"content" json:"content"`
	Url         string    `db:"url" json:"url"`
	Status      bool      `db:"status" json:"status"`
	RemindTime  time.Time `db:"remind_time" json:"remind_time"`
	IsDeleted   bool      `json:"is_deleted"`
}

func GetPrivateChat(m *PrivateMessage) *PrivateChat {
	var user1Id, user2Id int64
	if m.SenderId < m.RecipientId {
		user1Id = m.SenderId
		user2Id = m.RecipientId
	} else {
		user1Id = m.RecipientId
		user2Id = m.SenderId
	}
	return &PrivateChat{
		Id:             m.PrivateChatId,
		User1Id:        user1Id,
		User2Id:        user2Id,
		LastMsgContent: m.Content,
		LastSendTime:   m.SendTime,
	}
}
