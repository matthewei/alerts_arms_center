package ginserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/matthewei/alerts_arms_center/routers"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
	"net/http"
	"time"
)

type GinBaseConfigs struct {
	LogLevel string `yaml:"log_level"`
	TimeZone string `yaml:"time_zone"`
}

type GinHttpConfigs struct {
	Address        string        `yaml:"address"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	MaxHeaderBytes uint          `yaml:"max_header_bytes"`
	AllowIps       []string      `yaml:"allow_ips"`
}

type Configurations struct {
	BaseConfig GinBaseConfigs `yaml:"base_config"`
	HttpConfig GinHttpConfigs `yaml:"http_config"`
}

type GinServer struct {
	configuration Configurations
	*http.Server
	logger *zaplogger.ZapLogger
}

func New(configuration Configurations, acc *routers.AACRouter, logger *zaplogger.ZapLogger) *GinServer {
	gin.SetMode(configuration.BaseConfig.LogLevel)
	httpSrvHandler := &http.Server{
		Addr:           configuration.HttpConfig.Address,
		Handler:        acc.Router,
		ReadTimeout:    configuration.HttpConfig.ReadTimeout * time.Second,
		WriteTimeout:   configuration.HttpConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << configuration.HttpConfig.MaxHeaderBytes,
	}
	ginServer := new(GinServer)
	ginServer.configuration = configuration
	ginServer.Server = httpSrvHandler
	ginServer.logger = logger
	return ginServer
}

func (gs GinServer) HttpServerRun() {
	go func() {
		gs.logger.Infof("[INFO] HttpServerRun: %s", gs.configuration.HttpConfig.Address)
		if err := gs.ListenAndServe(); err != nil {
			gs.logger.Errorf(" [ERROR] HttpServerRun: %s err: %v", gs.configuration.HttpConfig.Address, err)
		}
	}()
}

func (gs GinServer) HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := gs.Shutdown(ctx); err != nil {
		gs.logger.Errorf(" [ERROR] HttpServerStop err: %v", err)
	}
	gs.logger.Infof(" [INFO] HttpServerStop stopped")
}
