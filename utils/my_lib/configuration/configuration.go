package configuration

import (
	"github.com/matthewei/alerts_arms_center/ginserver"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/configfileparser"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/msgkafka"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/site"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/trapsender"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
)

type AlertsArmsCenterConfigurations struct {
	SiteConfigurations      site.SiteConfigurations
	GinServerConfigurations ginserver.Configurations
	ZapLogConfigurations    zaplogger.Configurations
	MsgKafkaConfigurations  msgkafka.Configurations
	SnmpTrapConfigurations  trapsender.Configurations
}

func ParseConfiguration() (*AlertsArmsCenterConfigurations, *zaplogger.ZapLogger, error) {
	AlertsArmsCenterConfiguration := AlertsArmsCenterConfigurations{}
	parser, err := configfileparser.New()
	if err != nil {
		panic(err)
	}
	err = parser.ReadSection("ZapLogConfig", &AlertsArmsCenterConfiguration.ZapLogConfigurations)
	if err != nil {
		return nil, nil, err
	}
	logger, err := zaplogger.New(AlertsArmsCenterConfiguration.ZapLogConfigurations)
	if err != nil {
		return nil, nil, err
	}
	logger.Debug("read ZapLogConfig", AlertsArmsCenterConfiguration.ZapLogConfigurations)

	err = parser.ReadSection("SiteConfig", &AlertsArmsCenterConfiguration.SiteConfigurations)
	if err != nil {
		return nil, nil, err
	}
	logger.Debug("read SiteConfig", AlertsArmsCenterConfiguration.SiteConfigurations)

	err = parser.ReadSection("GinServerConfig", &AlertsArmsCenterConfiguration.GinServerConfigurations)
	if err != nil {
		return nil, nil, err
	}
	logger.Debug("read GinServerConfig", AlertsArmsCenterConfiguration.GinServerConfigurations)

	err = parser.ReadSection("KafkaConfig", &AlertsArmsCenterConfiguration.MsgKafkaConfigurations)
	if err != nil {
		return nil, nil, err
	}
	logger.Debug("read KafkaConfig", AlertsArmsCenterConfiguration.MsgKafkaConfigurations)

	err = parser.ReadSection("SnmpConfig", &AlertsArmsCenterConfiguration.SnmpTrapConfigurations)
	if err != nil {
		return nil, nil, err
	}
	logger.Debug("read SnmpConfig", AlertsArmsCenterConfiguration.SnmpTrapConfigurations)

	return &AlertsArmsCenterConfiguration, logger, nil
}
