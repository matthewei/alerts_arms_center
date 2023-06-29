package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/matthewei/alerts_arms_center/services"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/response"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/trapsender"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
	"github.com/matthewei/alerts_arms_center/utils/third_lib/telemetry"
	"net/http"
	"strconv"
)

type SnmpController struct {
	logger     *zaplogger.ZapLogger
	trapSender *trapsender.TrapSender
}

func New(logger *zaplogger.ZapLogger, trapsender *trapsender.TrapSender) *SnmpController {
	sc := new(SnmpController)
	sc.logger = logger
	sc.trapSender = trapsender
	return sc
}

func (sc *SnmpController) SnmpTrapRouterHandler(router *gin.RouterGroup) {
	sc.logger.Debug("SnmpTrapRouterHandler")
	router.POST("/trap", sc.snmpTrapHandler)
}

func (sc SnmpController) snmpTrapHandler(ctx *gin.Context) {
	zaplogger.WithContext(ctx).Debug("snmpTrapHandler")
	err := services.SnmpTrapService(ctx, sc.logger, sc.trapSender)
	if err != nil {
		telemetry.RequestTotal.WithLabelValues(strconv.FormatInt(http.StatusBadRequest, 10)).Inc()
		response.ResponseError(ctx, http.StatusBadRequest, err)
	}
	telemetry.RequestTotal.WithLabelValues("200").Inc()
	response.ResponseSuccess(ctx, "")
}
