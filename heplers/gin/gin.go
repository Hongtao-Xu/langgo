package helperGin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body" `
}

//SendOk 发送成功的消息
func SendOk(c *gin.Context, body interface{}) {
	resp := ApiResponse{}
	resp.Code = 0
	resp.Message = ""
	resp.Body = body
	c.JSON(http.StatusOK, resp)
}

//SendBad 发送失败的消息
func SendBad(c *gin.Context, code int, message string, body interface{}) {
	resp := ApiResponse{}
	resp.Code = code
	resp.Message = message
	resp.Body = body
	c.AbortWithStatusJSON(http.StatusOK, resp)
}
