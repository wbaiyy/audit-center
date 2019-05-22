package task

import (
	"audit-center/bucket"
	"audit-center/mydb"
	"audit-center/rabbit"
	"audit-center/tool"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func (tk *ConsumeTask) workAuditMessageRevoke(msg []byte) bool {
	var auditRevokeMsg rabbit.PersonRevokeAudit
	err := json.Unmarshal(msg, &auditRevokeMsg)
	if err !=nil {
		fmt.Println(err, "Unmarshal revoke message error")
		return false
	}

	auditStatus := tk.getMessageStatus(auditRevokeMsg.BussUuid, auditRevokeMsg.RevokeMark)
	if auditStatus != bucket.Auditing {
		logs.Warning(fmt.Sprintf("[status error] datebase status:%d, :message:%s",auditStatus, msg))
		return true
	}

	tk.updateMessageRevoke(auditRevokeMsg)

 	return true
}
/**
	通过业务ID和审核标识获取消息状态
 */
func (tk *ConsumeTask) getMessageStatus(bussUuid interface{}, auditMark string) int {
	db := mydb.DB

	var auditStatus int
	err := db.QueryRow("SELECT audit_status FROM audit_message_person " +
		"where business_uuid=? AND audit_mark=?", bussUuid, auditMark).Scan(&auditStatus)

	switch {
		case err == sql.ErrNoRows:
			logs.Warning(fmt.Sprintf("No message with that business_uuid:%v", bussUuid))
			return 0
		case err != nil:
			tool.FatalLog(err, "QueryRow error:")
	}

	return auditStatus
}
/**
	更新消息为撤销状态
 */
func (tk *ConsumeTask) updateMessageRevoke (auditRevokeMsg rabbit.PersonRevokeAudit) {
	db := mydb.DB

	args := []interface{}{
		bucket.ApplyCancel,
		auditRevokeMsg.RevokeRemark,
		auditRevokeMsg.RevokeUser,
		auditRevokeMsg.RevokeTime,
		auditRevokeMsg.BussUuid,
		auditRevokeMsg.RevokeMark,
	}

	_, err :=db.Exec("update audit_message_person set audit_status=?, message_remark=?,update_user=?,update_time=?" +
		" where business_uuid=? AND audit_mark=?", args...)

	tool.FatalLog(err, "Exce Revoke sql error")
}
