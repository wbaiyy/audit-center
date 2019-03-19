package rabbit

import (
	"audit_engine/tool"
	"encoding/json"
)

//审核响应消息模拟
func SimAuditBackMsg() []byte {
	var msg = AuditBackMsg{
		SiteCode:    "GB",
		BussUuid:    "13710",
		AuditStatus: 2,
		AuditRemark: "系统审核通过",
		AuditUid:    0,
		AuditUser:   "系统",
		AuditTime:   1535439935871,
	}

	b, err := json.Marshal(msg)
	tool.FatalLog(err, "publish json marshal fail")
	return b
}

//审核消息模拟
func SimAuditMsg() []byte {
	return []byte(`{"auditMark":"goods-price-check","bussData":"{\"calculatePrice\":52.02,\"catId\":11286,\"changeType\":1,\"chargePrice\":59.00000,\"goodSn\":\"YL4225902\",\"pipelineCode\":\"GB\",\"rate\":3.68,\"sysLabelId\":-1,\"virWhCode\":\"1433363\"}","bussUuid":"13710","createTime":1535427621607,"createUser":"huang","createUserId":0,"module":"goods","siteCode":"GB"}`)
}
