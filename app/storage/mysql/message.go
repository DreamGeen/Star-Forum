package mysql

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"star/constant/str"
	"star/models"
	"star/utils"
	"time"
)

const (
	listMessageCountSQL = `
                         select
                         (select count(1) from private_messages  where recipient_id=? and status=false) as privateMsgCount,
                         (select count(1) from  user_system_notice where recipient_id=? and status=false) as  systemCount,
                         (select count(1) from like_remind where  recipient_id=? and status=false) as  likeCount,
                         (select count(1) from mention_remind where  recipient_id=? and status=false) as  mentionCount,
                         (select count(1) from   reply_remind where  recipient_id=? and status=false) as replyCount;
                          `
	insertPrivateMsgSQL      = "insert into private_msg(private_message_id,sender_id,recipient_id,content,status,send_time) values (?,?,?,?,?,?)"
	updatePrivateChatSQL     = "update private_chat set last_message_content=? where user1_id=? and user2_id = ?"
	checkPrivateChatExistSQL = "select  private_chat_id from private_chat where user1_id=? and user2_id = ?"
	insertPrivateChatSQL     = "insert into  private_chat(private_chat_id,user1_id,user2_id,last_message_content, last_message_time)values (?,?,?,?)"
	insertSystemMsgSQL       = "insert into manager_system_notice(system_notice_id, title, content, type, status, recipient_id, manager_id, publish_time)values(?,?,?,?,?,?,?,?)"
	insertSystemMsgUserSQL   = "insert into user_system_notice(user_notice_id,system_notice_id, recipient_id,status)values(?,?,?,?)"
	queryAllUserIdSQL        = "select userId from user"
	insertLikeMessageSQL     = "insert into like_remind(id, source_id, source_type, content, url, status, sender_id, recipient_id, remind_time)values(?,?,?,?,?,?,?,?,?)"
	insertMentionMessageSQL  = "insert into mention_remind(id, source_id, source_type, content, url, status, sender_id, recipient_id, remind_time)values(?,?,?,?,?,?,?,?,?)"
	insertReplyMessageSQL    = "insert into reply_remind(id, source_id, source_type, content, url, status, sender_id, recipient_id, remind_time)values(?,?,?,?,?,?,?,?,?)"
	loadMessageSQL           = `
              select  private_message_id, sender_id, recipient_id, content, status, send_time, private_chat_id   from private_messages
              where  private_chat_id=? and send_time<? order by send_time desc limit ?;
            `
	getChatListSQL = `
    ( 
     SELECT
        user2_id AS other_user_id,
        last_message_content,
        last_message_time
     FROM private_chat
     WHERE user1_id = ?
    )
   UNION ALL
    (
      SELECT
        user1_id AS other_user_id,
        last_message_content,
        last_message_time
      FROM private_chat
      WHERE user2_id = ?
    )   ORDER BY last_message_time DESC
`
	getAllPrivateChat     = "select user1_id,user2_id from private_chat"
	insertBatchSystemUser = "insert into user_system_notice(user_notice_id,system_notice_id, recipient_id,status) values(:user_notice_id,:system_notice_id,:recipient_id,:status)"
)

func ListMessageCount(userId int64) (*models.Counts, error) {
	counts := new(models.Counts)
	if err := Client.Get(counts, listMessageCountSQL, userId, userId, userId, userId, userId); err != nil {
		utils.Logger.Error("get msg count error:", zap.Error(err))
		return nil, str.ErrMessageError
	}
	counts.TotalCount = counts.MentionCount + counts.LikeCount + counts.SystemCount + counts.ReplyCount + counts.PrivateMsgCount
	return counts, nil
}

func InsertPrivateMsg(message *models.PrivateMessage) (err error) {
	var tx *sqlx.Tx
	tx, err = Client.Beginx()
	if err != nil {
		utils.Logger.Error("begin tx error:", zap.Error(err))
		return str.ErrMessageError
	}
	// 使用 defer 确保在发生错误时回滚事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			zap.L().Error("Recovered from panic, transaction rolled back:", zap.Any("panic", p))
			err = str.ErrMessageError
		} else if err != nil {
			tx.Rollback()
			zap.L().Error("Transaction rolled back due to error:", zap.Error(err))
		}
	}()
	privateChat := models.GetPrivateChat(message)
	//检查会话是否存在
	var privateChatId int64
	if err = tx.Get(&privateChatId, checkPrivateChatExistSQL, privateChat.User1Id, privateChat.User2Id); err != nil {
		//查询出错
		utils.Logger.Error("select private_chat error:", zap.Error(err))
		return str.ErrMessageError
	}
	if privateChatId == 0 {
		//不存在则插入会话
		//生成会话id
		privateChatId = utils.GetID()
		privateChat.Id = privateChatId
		message.PrivateChatId = privateChatId
		if _, err = tx.Exec(insertPrivateChatSQL, privateChat.User1Id, privateChat.User2Id, privateChat.LastMsgContent, privateChat.LastSendTime); err != nil {
			utils.Logger.Error("insert private_chat error:", zap.Error(err))
			return str.ErrMessageError
		}
	} else {
		message.PrivateChatId = privateChatId
		//存在则更新数据
		if _, err = tx.Exec(updatePrivateChatSQL, privateChat.LastMsgContent, privateChat.User1Id, privateChat.User2Id); err != nil {
			utils.Logger.Error("update private_chat err:", zap.Error(err))
			return str.ErrMessageError
		}
	}
	//插入私信
	if _, err = tx.Exec(insertPrivateMsgSQL, message.PrivateChatId, message.Id, message.SenderId, message.RecipientId, message.Status, message.SendTime); err != nil {
		utils.Logger.Error("insert private_msg error:", zap.Error(err))
		return str.ErrMessageError
	}
	if err = tx.Commit(); err != nil {
		utils.Logger.Error("commit tx error:", zap.Error(err))
		return str.ErrMessageError
	}
	return nil
}

func InsertSystemMsg(message *models.SystemMessage) (err error) {
	var tx *sqlx.Tx
	tx, err = Client.Beginx()
	if err != nil {
		utils.Logger.Error("begin tx error:", zap.Error(err))
		return str.ErrMessageError
	}

	defer func() {
		if p := recover(); p != nil {
			utils.Logger.Error("recovered from panic, transaction rolled back:", zap.Any("panic", p))
			tx.Rollback()
			err = str.ErrMessageError
		} else if err != nil {
			utils.Logger.Error("Transaction rolled back due to error:", zap.Error(err))
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec(insertSystemMsgSQL, message.Id, message.Title, message.Content, message.Type, message.Status, message.RecipientId, message.ManagerId, message.PublishTime); err != nil {
		utils.Logger.Error("insert system_msg error:", zap.Error(err))
		return str.ErrMessageError
	}

	if message.Type == "single" {
		// 单个用户
		if _, err = tx.Exec(insertSystemMsgUserSQL, utils.GetID(), message.Id, message.RecipientId, false); err != nil {
			utils.Logger.Error("insert user_system_msg error", zap.Error(err))
			return str.ErrMessageError
		}
	} else {
		// 全体用户
		var userIds []int64
		if err = tx.Select(&userIds, queryAllUserIdSQL); err != nil {
			utils.Logger.Error("query all users id error", zap.Error(err))
			return str.ErrMessageError
		}
		systemMessageUsers := make([]interface{}, 0, len(userIds))
		for _, userId := range userIds {
			systemMessageUser := &models.SystemMessageUser{
				Id:              utils.GetID(),
				SystemMessageId: message.Id,
				RecipientId:     userId,
				Status:          message.Status,
			}
			systemMessageUsers = append(systemMessageUsers, systemMessageUser)
		}
		if _, err = tx.NamedExec(insertBatchSystemUser, systemMessageUsers); err != nil {
			utils.Logger.Error("insert system_msg error:", zap.Error(err))
			return str.ErrMessageError
		}
	}
	if err = tx.Commit(); err != nil {
		utils.Logger.Error("commit tx error:", zap.Error(err))
		return str.ErrMessageError
	}
	return nil
}

func InsertLikeMessage(message *models.RemindMessage) error {
	return insertRemindMessage(insertLikeMessageSQL, message)
}

func InsertMentionMessage(message *models.RemindMessage) error {
	return insertRemindMessage(insertMentionMessageSQL, message)
}

func InsertReplyMessage(message *models.RemindMessage) error {
	return insertRemindMessage(insertReplyMessageSQL, message)
}

func insertRemindMessage(query string, message *models.RemindMessage) error {
	if _, err := Client.Exec(query, message.Id, message.SourceId, message.SourceType, message.Content,
		message.Url, message.Status, message.SenderId, message.RecipientId, message.RemindTime); err != nil {
		utils.Logger.Error("insert remind message error:", zap.Error(err), zap.Any("message", message))
		return str.ErrMessageError
	}
	return nil
}

func LoadMessage(privateChatId int64, lastMsgTime time.Time, limit int) ([]*models.PrivateMessage, error) {
	var privateMessages []*models.PrivateMessage
	if err := Client.Select(&privateMessages, loadMessageSQL, privateChatId, lastMsgTime, limit); err != nil {
		utils.Logger.Error("load private_message error:", zap.Error(err), zap.Int64("privateChatId", privateChatId))
		return nil, str.ErrMessageError
	}
	return privateMessages, nil
}

func GetChatList(userId int64) ([]*models.PrivateChat, error) {
	var list []*models.PrivateChat
	if err := Client.Select(&list, getChatListSQL, userId, userId); err != nil {
		utils.Logger.Error("getChatList error:", zap.Error(err))
		return nil, str.ErrMessageError
	}
	return list, nil
}

func GetAllPrivateChat() ([]*models.PrivateChat, error) {
	var chats []*models.PrivateChat
	if err := Client.Select(&chats, getChatListSQL); err != nil {
		utils.Logger.Error("GetAllPrivateChat error:", zap.Error(err))
		return nil, str.ErrMessageError
	}
	return chats, nil
}
