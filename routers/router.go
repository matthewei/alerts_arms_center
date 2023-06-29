package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/matthewei/alerts_arms_center/controller"
	"github.com/matthewei/alerts_arms_center/middleware"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type AACRouter struct {
	Router     *gin.Engine
	logger     *zaplogger.ZapLogger
	Controller controller.SnmpController
}

func New(logger *zaplogger.ZapLogger) *AACRouter {
	logger.Info("SetupRouter")
	sysRouter := gin.Default()
	accRouter := new(AACRouter)
	accRouter.Router = sysRouter
	accRouter.logger = logger
	return accRouter
}

func (acc AACRouter) RegisterHealthRouterGroup() {
	acc.logger.Debug("RegisterHealthRouterGroup")
	acc.Router.GET("/health", func(c *gin.Context) {
		acc.logger.Debug("health check")
		c.JSON(200, gin.H{
			"Health": "OK",
		})
	})
}

func (acc AACRouter) RegisterMetricsRouterGroup() {
	acc.logger.Debug("RegisterMetricsRouterGroup")
	acc.Router.GET("/metrics", promHandler(promhttp.Handler()))
}

func promHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func (acc AACRouter) RegisterSnmpParserRouterGroup(sc *controller.SnmpController) {
	acc.logger.Debug("RegisterSnmpParserRouterGroup")
	snmpParserRouterGroup := acc.Router.Group("/snmp")
	snmpParserRouterGroup.Use(
		middleware.RecoveryMiddleware(acc.logger),
		middleware.TraceLoggerMiddleware(),
	)

	sc.SnmpTrapRouterHandler(snmpParserRouterGroup)
}
