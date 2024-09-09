package mysql

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
	"star/constant/str"
	"star/models"
	"star/utils"
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
	checkPrivateChatExistSQL = "select  count(1) from private_chat where user1_id=? and user2_id = ?"
	insertPrivateChatSQL     = "insert into  private_chat(user1_id,user2_id,last_message_content, last_message_time)values (?,?,?,?)"
	insertSystemMsgSQL       = "insert into manager_system_notice(system_notice_id, title, content, type, status, recipient_id, manager_id, publish_time)values(?,?,?,?,?,?,?,?)"
	insertSystemMsgUserSQL   = "insert into user_system_notice(user_notice_id,system_notice_id, recipient_id,status)values(?,?,?,?)"
	queryAllUserIdSQL        = "select userId from user"
)

func ListMessageCount(userId int64) (*models.Counts, error) {
	counts := new(models.Counts)
	if err := Client.Get(counts, listMessageCountSQL, userId, userId, userId, userId, userId); err != nil {
		zap.L().Error("get msg count error:", zap.Error(err))
		return nil, str.ErrMessageError
	}
	log.Println(counts.ReplyCount)
	counts.TotalCount = counts.MentionCount + counts.LikeCount + counts.SystemCount + counts.ReplyCount + counts.PrivateMsgCount
	return counts, nil
}

func InsertPrivateMsg(message *models.PrivateMessage) (err error) {
	var tx *sqlx.Tx
	tx, err = Client.Beginx()
	if err != nil {
		zap.L().Error("begin tx error:", zap.Error(err))
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
	//插入私信
	if _, err = tx.Exec(insertPrivateMsgSQL, message.Id, message.SenderId, message.RecipientId, message.Status, message.SendTime); err != nil {
		zap.L().Error("insert private_msg error:", zap.Error(err))
		return str.ErrMessageError
	}
	privateChat := models.GetPrivateChat(message)

	//检查会话是否存在
	var exist int
	if err = tx.Get(&exist, checkPrivateChatExistSQL, privateChat.User1Id, privateChat.User2Id); err != nil {
		//查询出错
		zap.L().Error("select private_chat error:", zap.Error(err))
		return str.ErrMessageError
	}
	if exist == 0 {
		//不存在则插入会话
		if _, err = tx.Exec(insertPrivateChatSQL, privateChat.User1Id, privateChat.User2Id, privateChat.LastMsgContent, privateChat.LastSendTime); err != nil {
			zap.L().Error("insert private_chat error:", zap.Error(err))
			return str.ErrMessageError
		}
	} else {
		//存在则更新数据
		if _, err = tx.Exec(updatePrivateChatSQL, privateChat.LastMsgContent, privateChat.User1Id, privateChat.User2Id); err != nil {
			zap.L().Error("update private_chat err:", zap.Error(err))
			return str.ErrMessageError
		}
	}
	if err = tx.Commit(); err != nil {
		zap.L().Error("commit tx error:", zap.Error(err))
		return str.ErrMessageError
	}
	return nil
}

func InsertSystemMsg(message *models.SystemMessage) (err error) {
	var tx *sqlx.Tx
	tx, err = Client.Beginx()
	if err != nil {
		zap.L().Error("begin tx error:", zap.Error(err))
		return str.ErrMessageError
	}

	defer func() {
		if p := recover(); p != nil {
			zap.L().Error("recovered from panic, transaction rolled back:", zap.Any("panic", p))
			tx.Rollback()
			err = str.ErrMessageError
		} else if err != nil {
			zap.L().Error("Transaction rolled back due to error:", zap.Error(err))
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec(insertSystemMsgSQL, message.Id, message.Title, message.Content, message.Type, message.Status, message.RecipientId, message.ManagerId, message.PublishTime); err != nil {
		zap.L().Error("insert system_msg error:", zap.Error(err))
	}

	if message.Type == "single" {
		// 单个用户
		if _, err = tx.Exec(insertSystemMsgUserSQL, utils.GetID(), message.Id, message.RecipientId, false); err != nil {
			zap.L().Error("insert user_system_msg error", zap.Error(err))
		}
	} else {
		// 全体用户
		var userIds []int64
		if err = tx.Select(&userIds, queryAllUserIdSQL); err != nil {
			zap.L().Error("query all users id error", zap.Error(err))
		}
		systemMessageUsers := make([]*models.SystemMessageUser, 0, len(userIds))
		for _, userId := range userIds {
			systemMessageUser := &models.SystemMessageUser{
				Id:              utils.GetID(),
				SystemMessageId: message.Id,
				RecipientId:     userId,
				Status:          message.Status,
			}
			systemMessageUsers = append(systemMessageUsers, systemMessageUser)
		}
		var query string
		var args []interface{}
		query, args, err = sqlx.In(insertSystemMsgUserSQL, systemMessageUsers)
		if err != nil {
			zap.L().Error("insert system_msg error:", zap.Error(err))
		}
		if _, err = tx.Exec(query, args...); err != nil {
			zap.L().Error("insert system_msg error:", zap.Error(err))
		}
	}
	if err = tx.Commit(); err != nil {
		zap.L().Error("commit tx error:", zap.Error(err))
		return str.ErrMessageError
	}
	return nil
}
