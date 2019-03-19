package config

import "audit_engine/rabbit"

var QueName = rabbit.QueueName{
	//SOA_PROVIDER
	"SOA_AUDIT_MSG": "auditMessage_OBS",
	//SYS_OBS
	"SOA_AUDIT_BACK_MSG":      "auditResult_SOA_GOODS",
	"OBS_RULE_CHANGE_MSG":     "obsRuleChange_OBS",
	"OBS_PERSON_AUDIT_RESULT": "obsAuditResult_OBS",
}
