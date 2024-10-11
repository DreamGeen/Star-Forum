package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"star/app/utils/logging"
)

var sf *snowflake.Node

// Init 雪花算法初始化
func Init(node int64) (err error) {
	sf, err = snowflake.NewNode(node)
	if err != nil {
		logging.Logger.Error("init snowflake err", zap.Error(err))
		return err
	}
	return nil
}

// GetID 获取用户id
func GetID() int64 {
	return sf.Generate().Int64()
}
