package rabbit

import (
	"audit-center/tool"
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
	return []byte(`{"auditMark":"goods-price-check","bussData":"{\"calculatePrice\":33.99,\"catId\":11293,\"changeType\":1,\"chargePrice\":666.99000,\"freightPrice\":0,\"goodSn\":\"249353702\",\"pipelineCode\":\"GB\",\"rate\":0.31,\"saleMark\":1,\"sysLabelId\":8,\"virWhCode\":\"1433363\"}","bussUuid":"3210961","createTime":1560928203971,"createUid":2,"createUser":"wangbei","module":"goods","siteCode":"GB"}`)
}
