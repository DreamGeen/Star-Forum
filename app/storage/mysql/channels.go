package mysql

import (
	"log"
	"star/models"
)

const (
	insertMsgSQL = "insert into group_messages(chatId,sendtime,content,sdUserId,communityId) values (?,?,?,?,?)"
	queryMsgSQL  = `
                            select g.sendTime,g.content,g.sdUserId,u.userName,u.img
                           from group_messages as g,user as u
                           where g.sdUserId=u.userId and g.communityId=? 
                            order by g.chatId desc limit ?
                               `

	queryLastMsgIdSQL    = "select communityId,lastMsgId from community where communityId=?"
	getMessagesBeforeSQL = `
                             select g.sendTime,g.content,g.sdUserId,u.userName,u.img
                             from group_messages as g,user as u
                             where  g.sdUserId=u.userId and g.chatId< ? and g.communityId=?
                             order by g.chatId desc limit ?
                            `
	batchInsertMsgSQL = "insert into group_messages(content,sendTime,chatId,sdUserId,communityId) values (:content,:sendTime,:chatId,:sdUserId,:communityId)"
)

func SaveMsg(message *models.GroupMessage) error {
	_, err := Client.Exec(insertMsgSQL, message.SendTime, message.Content, message.SdUserId, message.CommunityId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SaveBathMsg 批量插入消息
func SaveBathMsg(messages []*models.GroupMessage) error {
	_, err := Client.NamedExec(batchInsertMsgSQL, messages)
	if err != nil {
		log.Println(err)
	}
	//query, args, err := sqlx.In(batchInsertMsgSQL, messages...)
	//fmt.Println(query)
	//fmt.Println(args)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//_, err = db.Exec(query, args...)
	return err
}

// GetRecentMsg 加载聊天记录
func GetRecentMsg(communityId int64, number int) ([]*models.GroupMessage, error) {
	messages := make([]*models.GroupMessage, number)
	if err := Client.Select(&messages, queryMsgSQL, communityId, number); err != nil {
		log.Println("load msg error,err:", err)
		return nil, err
	}
	return messages, nil
}

func GetMessagesBefore(communityId, lastChatId int64, number int64) ([]*models.GroupMessage, error) {
	messages := make([]*models.GroupMessage, number)
	if err := Client.Select(&messages, getMessagesBeforeSQL, lastChatId, communityId, number); err != nil {
		return nil, err
	}
	return messages, nil
}
