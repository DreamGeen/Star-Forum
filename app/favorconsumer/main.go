package favorconsumer

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/app/storage/mq"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/utils"
	"sync"
	"time"
)

type collectMsg struct {
	count   int64
	userIds map[int64]struct{}
}

var conn *amqp091.Connection
var channel *amqp091.Channel

var likePostBuffer map[int64]int

var likeCommentBuffer map[int64]int

var collectPostBuffer map[int64]*collectMsg
var deleteCollect map[int64][]int64

var muLikePost sync.RWMutex
var muLikeComment sync.RWMutex
var muCollect sync.RWMutex

func failOnError(err error, msg string) {
	if err != nil {
		utils.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		utils.Logger.Error("close rabbitmq conn error", zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		utils.Logger.Error("close rabbitmq channel error", zap.Error(err))
		panic(err)
	}
}
func main() {
	var err error
	conn, err = amqp091.Dial(mq.ReturnRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer CloseMQ()

	go likePostConsumer()
	go likeCommentConsumer()
	go asyncUpdatePostLike()
	go asyncUpdateCommentLike()
	go collectPostConsumer()
	go asyncUpdatePostCollect()
}

func likePostConsumer() {
	delivery, err := channel.Consume(str.LikePost, str.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var message *models.Post
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			utils.Logger.Error("unmarshal json fail")
			continue
		}
		muLikePost.Lock()
		likePostBuffer[message.PostId] += message.Star
		muLikePost.Unlock()
	}
}
func likeCommentConsumer() {
	delivery, err := channel.Consume(str.LikeComment, str.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var comment *models.Comment
		if err := json.Unmarshal(msg.Body, &comment); err != nil {
			utils.Logger.Error("unmarshal json fail")
			continue
		}
		muLikeComment.Lock()
		likeCommentBuffer[comment.PostId] += comment.Star
		muLikeComment.Unlock()
	}
}

func collectPostConsumer() {
	delivery, err := channel.Consume(str.CollectPost, str.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var post *models.Collect
		if err := json.Unmarshal(msg.Body, &post); err != nil {
			utils.Logger.Error("unmarshal json fail")
			continue
		}
		muCollect.Lock()
		collectBuffer := collectPostBuffer[post.PostId]
		collectBuffer.count += post.Collection
		if post.Collection > 0 {
			collectBuffer.userIds[post.UserId] = struct{}{}
		} else {
			if _, exists := collectBuffer.userIds[post.UserId]; exists {
				delete(collectBuffer.userIds, post.UserId)
			} else {
				deleteCollect[post.PostId] = append(deleteCollect[post.PostId], post.UserId)
			}
		}
		muCollect.Unlock()
	}
}
func asyncUpdatePostLike() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			muLikePost.Lock()
			for postId, count := range likePostBuffer {
				if err := mysql.UpdatePostLike(postId, count); err != nil {
					utils.Logger.Error("async post update like fail", zap.Error(err), zap.Int64("post_id", postId), zap.Int("count", count))
				}
				delete(likePostBuffer, postId)
			}
			muLikePost.Unlock()
		}
	}
}
func asyncUpdateCommentLike() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			muLikeComment.Lock()
			for commentId, count := range likeCommentBuffer {
				if err := mysql.UpdateCommentLike(commentId, count); err != nil {
					utils.Logger.Error("async comment update like fail", zap.Error(err), zap.Int64("comment_id", commentId), zap.Int("count", count))
				}
				delete(likeCommentBuffer, commentId)
			}
			muLikeComment.Unlock()
		}
	}
}

func asyncUpdatePostCollect() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			muCollect.Lock()
			for postId, userCollect := range collectPostBuffer {
				if err := mysql.AddCollect(postId, userCollect); err != nil {
					utils.Logger.Error("async post add collectMsg fail",
						zap.Error(err),
						zap.Int64("post_id", postId),
						zap.Int("count", userCollect.count))
				}
				delete(collectPostBuffer, postId)
			}
			for postId, userIds := range deleteCollect {
				if err := mysql.DeleteCollect(postId, userIds); err != nil {
					utils.Logger.Error("async post delete collectMsg fail",
						zap.Error(err),
						zap.Int64("post_id", postId))
				}
				delete(deleteCollect, postId)
			}
			muCollect.Unlock()
		}
	}
}
