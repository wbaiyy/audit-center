package task

import (
	"audit_engine/mydb"
	"audit_engine/rabbit"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

//接收消息审核任务
func (tk *ConsumeTask) workAuditMessage(msg []byte) bool {
	//审核数据
	var audMsg rabbit.AuditMsg
	err := json.Unmarshal(msg, &audMsg)
	if err != nil {
		log.Println(err, "unmarshal audit message fail")
		return false
	}
	log.Printf("auditMsg: %+v\n", audMsg)

	//业务数据
	var audBd rabbit.BusinessData
	err = json.Unmarshal([]byte(audMsg.BussData), &audBd)
	if err != nil {
		log.Println(err, "unmarshal business data fail")
		return false
	}
	log.Printf("bussData: %+v\n", audBd)

	//hash map 规则
	hashRuleTypes := tk.GetRuleItems()
	audType, ok := hashRuleTypes[audMsg.AuditMark]
	if !ok {
		log.Println(audMsg.AuditMark, "hash key not exist")
		return false
	}
	log.Printf("ruleList: %+v", audType)

	//规则校验(rt)
	audStat, rulMch := RunRuleMatch(&audBd, &audType)
	log.Println("matchResult(20：引擎通过,21：引擎拒绝,22：规则全不匹配，自动通过,30：转人工审核)----->", audStat)

	//自动通过|驳回|转人工审核（写db)
	tk.insertAuditMsg(audMsg, audBd, &audType, audStat, rulMch)

	return true
}

//审核消息入库
func (tk *ConsumeTask) insertAuditMsg(audMsg rabbit.AuditMsg, audBd rabbit.BusinessData, audType *AuditType, audStat int, rulMch RuleMatch) {
	db := mydb.DB

	//检测审核规则是否为空
	if len(audType.RuleList) == 0 {
		log.Println("audit rule list is empty")
	}

	//自动通过或拒绝
	fields := []string{
		"site_code",
		"rule_id",
		"template_id",
		"audit_sort",
		"audit_mark",
		"audit_name",
		"business_uuid",
		"business_data",
		"create_user",
		"workflow_id",
		"audit_status",
		"module",
		"create_time",
		"message_remark", //add at 2019-1-19
	}
	args := []interface{}{
		audMsg.SiteCode,
		rulMch.RuleId,
		audType.TypeId,
		audType.AuditSort,
		audType.AuditMark,
		fmt.Sprintf("sku:%s-%s", audBd.GoodSn, audType.TypeTitle),
		audMsg.BussUuid,
		audMsg.BussData,
		audMsg.CreateUser,
		rulMch.FlowId,
		audStat,
		audMsg.Module,
		time.Now().Unix(),
		audMsg.Remark,
	}

	//系统审核结束，非人工审核，db新增系统审核明细
	audOver := audStat != AuditStatus[ObsAudit]
	if audOver {
		fields = append(fields,
			"audit_desc",
			"update_user",
			"update_time",
		)
		args = append(args,
			GetAdStatDesc(audStat, ""),
			"系统",
			time.Now().Unix(),
		)
	}

	messageTable := "audit_message_person"
	if audOver {
		messageTable = "audit_message_system"
	}

	//sql组装
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		messageTable,
		strings.Join(fields, ","),
		strings.Repeat("?,", len(fields))[:2*len(fields)-1],
	)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err, "insert into audit_message Prepare fail")
		return
	}
	defer stmt.Close()

	//sql执行
	result, err := stmt.Exec(args...)
	if err != nil {
		log.Println(err, "insert into audit_message Exec fail")
		return
	}
	lastId, err := result.LastInsertId()
	log.Printf("success insert id: %d", lastId)

	//自动通过或拒绝，发布消息
	if audOver {
		tk.sendBackMsg(lastId, messageTable)
	}
}
