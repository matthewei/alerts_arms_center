package main

import (
	"github.com/matthewei/alerts_arms_center/controller"
	"github.com/matthewei/alerts_arms_center/ginserver"
	"github.com/matthewei/alerts_arms_center/routers"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/configuration"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/trapsender"
	"github.com/matthewei/alerts_arms_center/utils/third_lib/telemetry"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	telemetry.Init()
	configuration, logger, err := configuration.ParseConfiguration()
	if err != nil {
		return
	}
	logger.Debug("configuration", configuration)
	siteName := configuration.SiteConfigurations.SiteName
	logger.Info("alerts_arms_center-> ", siteName)
	trapSender := trapsender.New(configuration.SnmpTrapConfigurations, logger, configuration.SiteConfigurations)
	//判断trapSender是否为空，如果为空，直接返回
	if trapSender == nil {
		logger.Error("trapSender is nil")
		return
	}
	accController := controller.New(logger, trapSender)

	accRouter := routers.New(logger)
	{
		accRouter.RegisterHealthRouterGroup()
		accRouter.RegisterMetricsRouterGroup()
		accRouter.RegisterSnmpParserRouterGroup(accController)
	}

	server := ginserver.New(configuration.GinServerConfigurations, accRouter, logger)
	server.HttpServerRun()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.HttpServerStop()

}
