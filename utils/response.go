package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应信息结构体
type Response struct {
	Code ResponseCode `json:"code"` //响应码
	Msg  interface{}  `json:"msg"`  //响应信息
	Data interface{}  `json:"data"` //补充响应信息
}

// ResponseMessage 默认响应(传响应码)
func ResponseMessage(c *gin.Context, code ResponseCode) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  code.getMsg(),
		Data: nil,
	},
	)
}

// ResponseMessageWithData 可以设置补充响应信息(传错误)
func ResponseMessageWithData(c *gin.Context, code ResponseCode, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  code.getMsg(),
		Data: data,
	},
	)
}

// ResponseErr 默认响应(传错误)
func ResponseErr(c *gin.Context, err error) {
	c.JSON(http.StatusOK, &Response{
		Code: getCode(err),
		Msg:  err.Error(),
		Data: nil,
	})
}

// ResponseErrWithData 可以设置补充响应信息(传错误)
func ResponseErrWithData(c *gin.Context, err error) {
	c.JSON(http.StatusOK, &Response{
		Code: getCode(err),
		Msg:  err.Error(),
		Data: nil,
	})
}
