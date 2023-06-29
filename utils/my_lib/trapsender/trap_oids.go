package trapsender

import "github.com/k-sone/snmpgo"

type TrapOIDsConfigurations struct {
	SnmpTrapOid         string `yaml:"snmp_trap_oid"`
	SiteNameOid         string `yaml:"region_name_oid"`
	SysModuleOid        string `yaml:"sys_module_oid"`
	AlertNameOid        string `yaml:"alarm_name_oid"`
	AlertStatusOid      string `yaml:"alarm_status_oid"`
	AlertTypeOid        string `yaml:"alarm_type_oid"`
	AlertResourceOid    string `yaml:"alarm_resource_oid"`
	AlertHostOid        string `yaml:"alarm_host_oid"`
	AlertSeverityOid    string `yaml:"alarm_level_oid"`
	AlertSendTimeOid    string `yaml:"alert_send_time_oid"`
	AlertSummaryOid     string `yaml:"alarm_summary_oid"`
	AlertDescriptionOid string `yaml:"alarm_description_oid"`
	AlertImpactOid      string `yaml:"alarm_impact_oid"`
}

type SnmpgoTrapOID struct {
	SnmpTrapOid         *snmpgo.Oid
	SiteNameOid         *snmpgo.Oid
	SysModuleOid        *snmpgo.Oid
	AlertNameOid        *snmpgo.Oid
	AlertStatusOid      *snmpgo.Oid
	AlertTypeOid        *snmpgo.Oid
	AlertResourceOid    *snmpgo.Oid
	AlertHostOid        *snmpgo.Oid
	AlertSeverityOid    *snmpgo.Oid
	AlertSendTimeOid    *snmpgo.Oid
	AlertSummaryOid     *snmpgo.Oid
	AlertDescriptionOid *snmpgo.Oid
	AlertImpactOid      *snmpgo.Oid
}
