package services

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/alertbuckets"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/response"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/trapsender"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
	"io"
	"net/http"
)

func SnmpTrapService(context *gin.Context, logger *zaplogger.ZapLogger, trapSender *trapsender.TrapSender) error {
	req := context.Request
	defer req.Body.Close()
	// alertmanager发送的原始数据
	postBody, err := io.ReadAll(req.Body)
	if err != nil {
		response.ResponseError(context, http.StatusUnprocessableEntity, err)
		return err
	}
	logger.Debug("alertmanager-postBody", string(postBody))
	alertsData := alertbuckets.AlertsData{}
	if err := json.Unmarshal(postBody, &alertsData); err != nil {
		response.ResponseError(context, http.StatusUnprocessableEntity, err)
		return err
	}
	logger.Debug("alertsData", alertsData)
	// 取出相关的alerts信息
	alertsBuckets := alertbuckets.New(alertsData.Alerts)
	logger.Debug("alertsBuckets", alertsBuckets)
	// Send the alerts to the snmp-trapper:
	err = trapSender.SendAlertTraps(alertsBuckets)
	if err != nil {
		return err
	}
	return nil
}
