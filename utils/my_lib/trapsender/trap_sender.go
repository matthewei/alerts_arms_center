package trapsender

import (
	"github.com/k-sone/snmpgo"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/alertbuckets"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/site"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/timeutil"
	"github.com/matthewei/alerts_arms_center/utils/my_lib/zaplogger"
	"github.com/matthewei/alerts_arms_center/utils/third_lib/telemetry"
	"reflect"
	"time"
)

type Configurations struct {
	SnmpConfig     *SnmpConfigurations
	TrapOIDsConfig TrapOIDsConfigurations
}

type SnmpConfigurations struct {
	SNMPDestination            []string      `yaml:"snmp_destination"`
	SNMPRetries                uint          `yaml:"snmp_retries"`
	SNMPTimeout                time.Duration `yaml:"snmp_timeout"`
	SNMPVersion                string        `yaml:"snmp_version"`
	SNMPCommunity              string        `yaml:"snmp_community"`
	SNMPAuthenticationEnabled  bool          `yaml:"snmp_authentication_enabled"`
	SNMPAuthenticationProtocol string        `yaml:"snmp_authentication_protocol"`
	SNMPAuthenticationUsername string        `yaml:"snmp_authentication_username"`
	SNMPAuthenticationPassword string        `yaml:"snmp_authentication_password"`
	SNMPPrivateEnabled         bool          `yaml:"snmp_private_enabled"`
	SNMPPrivateProtocol        string        `yaml:"snmp_private_protocol"`
	SNMPPrivatePassword        string        `yaml:"snmp_private_password"`
	SNMPSecurityEngineID       string        `yaml:"snmp_security_engine_id"`
	SNMPContextEngineID        string        `yaml:"snmp_context_engine_id"`
	SNMPContextName            string        `yaml:"snmp_context_name"`
}

type TrapSender struct {
	configurations     Configurations
	siteConfigurations site.SiteConfigurations
	TrapOID            SnmpgoTrapOID
	logger             *zaplogger.ZapLogger
}

func New(configuration Configurations, logger *zaplogger.ZapLogger, siteConfig site.SiteConfigurations) *TrapSender {
	TrapSender := new(TrapSender)
	//1. snmp config
	TrapSender.configurations.SnmpConfig = generateSnmpConfig(configuration, logger)
	//2. trap oid config
	TrapSender.configurations.TrapOIDsConfig = configuration.TrapOIDsConfig
	//3. generate trap oid
	TrapSender.TrapOID = generateSnmpTrapsOid(configuration, logger)
	if TrapSender.TrapOID == (SnmpgoTrapOID{}) {
		return nil
	}
	TrapSender.configurations = configuration
	TrapSender.logger = logger
	TrapSender.siteConfigurations = siteConfig
	return TrapSender
}

func (ts TrapSender) SendAlertTraps(alerts *alertbuckets.AlertBuckets) error {
	traps := ts.generateSnmpTraps(alerts)
	ts.logger.Debug("traps", traps)
	connections, err := ts.connect()
	if err != nil {
		ts.logger.Error("failed to connect to snmp server!", err)
		return err
	}
	defer func() {
		for _, connection := range connections {
			connection.Close()
		}
	}()
	for _, connection := range connections {
		for _, trap := range traps {
			err = connection.V2Trap(trap)
			if err != nil {
				ts.logger.Error("failed to send trap!-> ", err)
				telemetry.SNMPErrorTotal.WithLabelValues().Inc()
				return err
			}
			ts.logger.Debug("send trap successfully!")
			telemetry.SNMPSentTotal.WithLabelValues().Inc()
		}
	}
	return nil
}

// generateSnmpTrap generate snmp trap
func (ts TrapSender) generateSnmpTraps(alerts *alertbuckets.AlertBuckets) []snmpgo.VarBinds {
	var (
		traps []snmpgo.VarBinds
	)
	for _, alert := range alerts.AlertBuckets {
		varBinds := ts.generateVarBinds(alert)
		traps = append(traps, varBinds)
	}
	return traps
}

func (ts TrapSender) generateVarBinds(alert alertbuckets.AlertBucket) snmpgo.VarBinds {
	var (
		varBinds snmpgo.VarBinds
	)
	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, ts.TrapOID.SnmpTrapOid))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.SnmpTrapOid, snmpgo.NewOctetString([]byte(ts.siteConfigurations.SiteName))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.SiteNameOid, snmpgo.NewOctetString([]byte(ts.siteConfigurations.SiteName))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.SysModuleOid, snmpgo.NewOctetString([]byte(ts.siteConfigurations.SysModule))))

	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertNameOid, snmpgo.NewOctetString([]byte(alert.Labels["alertname"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertTypeOid, snmpgo.NewOctetString([]byte(alert.Labels["alertType"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertResourceOid, snmpgo.NewOctetString([]byte(alert.Labels["resource"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertHostOid, snmpgo.NewOctetString([]byte(alert.Labels["hostname"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertSeverityOid, snmpgo.NewOctetString([]byte(alert.Labels["severity"]))))
	if alert.Status == "firing" {
		varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertStatusOid, snmpgo.NewCounter64(0)))
		TimeUnixMilli := timeutil.RFC3339TOUnixMilli(alert.StartsAt.Format(time.RFC3339))
		varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertSendTimeOid, snmpgo.NewCounter64(uint64(TimeUnixMilli))))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertStatusOid, snmpgo.NewCounter64(1)))
		TimeUnixMilli := timeutil.RFC3339TOUnixMilli(alert.EndsAt.Format(time.RFC3339))
		varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertSendTimeOid, snmpgo.NewCounter64(uint64(TimeUnixMilli))))
	}
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertSummaryOid, snmpgo.NewOctetString([]byte(alert.Annotations["summary"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertDescriptionOid, snmpgo.NewOctetString([]byte(alert.Annotations["description"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(ts.TrapOID.AlertImpactOid, snmpgo.NewOctetString([]byte(alert.Annotations["impact"]))))
	ts.logger.Debug("generateVarBinds varBinds", varBinds)
	return varBinds
}

// connect connect to snmp server
func (ts TrapSender) connect() ([]*snmpgo.SNMP, error) {
	snmpArguments := []snmpgo.SNMPArguments{}
	for _, destination := range ts.configurations.SnmpConfig.SNMPDestination {
		snmpArgument := snmpgo.SNMPArguments{
			Address: destination,
			Retries: ts.configurations.SnmpConfig.SNMPRetries,
			Timeout: ts.configurations.SnmpConfig.SNMPTimeout,
		}
		if ts.configurations.SnmpConfig.SNMPVersion == "V2c" {
			snmpArgument.Version = snmpgo.V2c
			snmpArgument.Community = ts.configurations.SnmpConfig.SNMPCommunity
		}
		if ts.configurations.SnmpConfig.SNMPVersion == "V3" {
			snmpArgument.Version = snmpgo.V3
			snmpArgument.UserName = ts.configurations.SnmpConfig.SNMPAuthenticationUsername
			if ts.configurations.SnmpConfig.SNMPAuthenticationEnabled && ts.configurations.SnmpConfig.SNMPPrivateEnabled {
				snmpArgument.SecurityLevel = snmpgo.AuthPriv
			} else if ts.configurations.SnmpConfig.SNMPAuthenticationEnabled {
				snmpArgument.SecurityLevel = snmpgo.AuthNoPriv
			} else {
				snmpArgument.SecurityLevel = snmpgo.NoAuthNoPriv
			}
			if ts.configurations.SnmpConfig.SNMPPrivateEnabled {
				snmpArgument.AuthProtocol = snmpgo.AuthProtocol(ts.configurations.SnmpConfig.SNMPAuthenticationProtocol)
				snmpArgument.AuthPassword = ts.configurations.SnmpConfig.SNMPAuthenticationPassword
			}
			snmpArgument.SecurityEngineId = ts.configurations.SnmpConfig.SNMPSecurityEngineID
			snmpArgument.ContextEngineId = ts.configurations.SnmpConfig.SNMPContextEngineID
			snmpArgument.ContextName = ts.configurations.SnmpConfig.SNMPContextName
		}
		snmpArguments = append(snmpArguments, snmpArgument)
	}
	snmpslist := []*snmpgo.SNMP{}
	for _, snmpArgument := range snmpArguments {
		snmp, err := snmpgo.NewSNMP(snmpArgument)
		if err != nil {
			return nil, err
		}
		err = snmp.Open()
		if err != nil {
			ts.logger.Error("snmp open failed!-> ", err)
			return nil, err
		}
		snmpslist = append(snmpslist, snmp)
	}
	return snmpslist, nil
}

// generateSnmpConfig generates the SNMP configuration
func generateSnmpConfig(configurations Configurations, logger *zaplogger.ZapLogger) *SnmpConfigurations {
	snmpConfig := configurations.SnmpConfig
	isV2c := snmpConfig.SNMPVersion == "2c"

	snmpDestionations := []string{}
	for _, destination := range snmpConfig.SNMPDestination {
		snmpDestionations = append(snmpDestionations, destination)
	}
	snmpConfigurations := new(SnmpConfigurations)
	snmpConfigurations = &SnmpConfigurations{
		SNMPVersion:     snmpConfig.SNMPVersion,
		SNMPDestination: snmpDestionations,
		SNMPRetries:     snmpConfig.SNMPRetries,
		SNMPTimeout:     snmpConfig.SNMPTimeout,
	}

	if isV2c {
		snmpConfigurations.SNMPCommunity = snmpConfig.SNMPCommunity
	}

	if !isV2c {
		snmpConfigurations.SNMPAuthenticationPassword = snmpConfig.SNMPAuthenticationPassword
		snmpConfigurations.SNMPSecurityEngineID = snmpConfig.SNMPSecurityEngineID
		snmpConfigurations.SNMPContextEngineID = snmpConfig.SNMPContextEngineID
		snmpConfigurations.SNMPContextName = snmpConfig.SNMPContextName
	}
	if isV2c && (snmpConfig.SNMPAuthenticationEnabled || snmpConfig.SNMPPrivateEnabled) {
		logger.Error("SNMP authentication or private only available with SNMP v3")
		return nil
	}
	if !snmpConfig.SNMPAuthenticationEnabled && snmpConfig.SNMPPrivateEnabled {
		logger.Errorln("SNMP private encryption requires authentication enabled.")
	}
	if snmpConfig.SNMPAuthenticationEnabled {
		snmpConfigurations.SNMPAuthenticationEnabled = snmpConfig.SNMPAuthenticationEnabled
		snmpConfigurations.SNMPAuthenticationProtocol = snmpConfig.SNMPAuthenticationProtocol
		snmpConfigurations.SNMPAuthenticationUsername = snmpConfig.SNMPAuthenticationUsername
		snmpConfigurations.SNMPAuthenticationPassword = snmpConfig.SNMPAuthenticationPassword
	}
	if snmpConfig.SNMPPrivateEnabled {
		snmpConfigurations.SNMPPrivateEnabled = snmpConfig.SNMPPrivateEnabled
		snmpConfigurations.SNMPPrivateProtocol = snmpConfig.SNMPPrivateProtocol
		snmpConfigurations.SNMPPrivatePassword = snmpConfig.SNMPPrivatePassword
	}
	logger.Debug("snmpConfiguration: ", snmpConfigurations)
	return snmpConfigurations
}

// generate trap oid
func generateSnmpTrapsOid(configuration Configurations, logger *zaplogger.ZapLogger) SnmpgoTrapOID {
	trapoid := SnmpgoTrapOID{}
	oid := reflect.TypeOf(trapoid)
	config := reflect.TypeOf(configuration.TrapOIDsConfig)
	for i := 0; i < config.NumField(); i++ {
		for j := 0; j < oid.NumField(); j++ {
			if config.Field(i).Name == oid.Field(j).Name {
				newOid, err := snmpgo.NewOid(reflect.ValueOf(configuration.TrapOIDsConfig).Field(i).Interface().(string))
				if err != nil {
					logger.Error("generate trap oid error: ", err)
					return SnmpgoTrapOID{}
				}
				// generate trap oid
				valueOf := reflect.ValueOf(newOid)
				reflect.ValueOf(&trapoid).Elem().Field(j).Set(valueOf)
			}
		}
	}
	return trapoid
}
