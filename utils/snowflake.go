package utils

import (
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var sf *snowflake.Node

// Init 雪花算法初始化
func Init(node int64) (err error) {
	sf, err = snowflake.NewNode(node)
	if err != nil {
		zap.L().Error("init snowflake err", zap.Error(err))
		return err
	}
	return nil
}

// GetID 获取用户id
func GetID() int64 {
	return sf.Generate().Int64()
}
