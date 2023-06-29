package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/matthewei/alerts_arms_center/utils/my_lib/response"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
)

// RecoveryMiddleware捕获所有panic，并且返回错误信息
func RecoveryMiddleware(logger *zaplogger.ZapLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//先做一下日志记录
				logger.Errorf("error", fmt.Sprint(err))
				logger.Errorf("stack", string(debug.Stack()))
				response.ResponseError(ctx, 500, errors.New(fmt.Sprint(err)))
				return
			}
		}()
		ctx.Next()
	}
}
