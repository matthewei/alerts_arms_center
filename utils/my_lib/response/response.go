package response

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseCode int

type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceId   interface{}  `json:"trace_id"`
	Stack     interface{}  `json:"stack"`
}

func ResponseError(c *gin.Context, code ResponseCode, err error) {
	trace, _ := c.Get("traceId")
	//traceContext, _ := trace.(*zaplogger.TraceContext)
	traceId := trace
	//if traceContext != nil {
	//	traceId = traceContext.TraceId
	//}

	stack := ""

	resp := &Response{ErrorCode: code, ErrorMsg: err.Error(), Data: "", TraceId: traceId, Stack: stack}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(200, err)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	trace, _ := c.Get("traceId")
	fmt.Println("------------------------", trace)
	traceId := trace
	resp := &Response{ErrorCode: http.StatusOK, ErrorMsg: "", Data: data, TraceId: traceId}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
