package config

import "audit-center/rabbit"

var QueName = rabbit.QueueName{
	//SOA_PROVIDER
	"SOA_AUDIT_MSG": "auditMessage_OBS",
	//SYS_OBS
	"SOA_AUDIT_BACK_MSG":      "auditResult_SOA_GOODS",
	"OBS_RULE_CHANGE_MSG":     "obsRuleChange_OBS",
	"OBS_PERSON_AUDIT_RESULT": "obsAuditResult_OBS",
	//coupon审核
	"TASK_AUDIT_MSG": "auditMessageCoupon_OBS",
	"TASK_AUDIT_REVOKE_MSG": "auditRevoke_OBS",
	"TASK_AUDIT_BACK_MSG": "auditResultNotify_GB",
}

func IsValidateQueueName(queueName string) bool  {
	for _, queue  := range QueName{
		if queue == queueName {
			return true
		}
	}

	return false
}

/**
	是否是soa队列
 */
func IsSoaQueue(queueName string) bool {
	soaQueues := [...]string{
		QueName["SOA_AUDIT_MSG"],
	}

	for _, soaQueue := range soaQueues{
		if soaQueue == queueName {
			return true
		}
	}

	return  false
}

/**
	是否是GB队列
 */
func IsGbQueue(queueName string) bool {
	soaQueues := [...]string{
		QueName["TASK_AUDIT_MSG"],
	}

	for _, soaQueue := range soaQueues{
		if soaQueue == queueName {
			return true
		}
	}

	return  false
}
