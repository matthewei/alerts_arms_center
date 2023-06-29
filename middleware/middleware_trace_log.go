package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
	"go.uber.org/zap"
)

type TraceLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (tw TraceLogWriter) Write(p []byte) (int, error) {
	if n, err := tw.body.Write(p); err != nil {
		return n, err
	}
	return tw.ResponseWriter.Write(p)
}

func TraceLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 每个请求生成的请求traceId具有全局唯一性
		u1, _ := uuid.NewUUID()
		traceId := u1.String()
		zaplogger.NewContext(ctx, zap.String("traceId", traceId))

		// 为日志添加请求的地址以及请求参数等信息
		zaplogger.NewContext(ctx, zap.String("request.method", ctx.Request.Method))
		headers, _ := json.Marshal(ctx.Request.Header)
		zaplogger.NewContext(ctx, zap.String("request.headers", string(headers)))
		zaplogger.NewContext(ctx, zap.String("request.url", ctx.Request.URL.String()))

		// 将请求参数json序列化后添加进日志上下文
		if ctx.Request.Form == nil {
			ctx.Request.ParseMultipartForm(32 << 20)
		}
		form, _ := json.Marshal(ctx.Request.Form)
		zaplogger.NewContext(ctx, zap.String("request.params", string(form)))

		ctx.Next()
	}
}
