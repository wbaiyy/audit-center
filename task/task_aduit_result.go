package task

import (
	"audit_engine/bucket"
	"audit_engine/config"
	"audit_engine/mydb"
	"audit_engine/rabbit"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//同步审核结果任务
func (tk *ConsumeTask) workUpdateAuditResult(msg []byte) bool {
	db := mydb.DB

	//obs audit result
	var par rabbit.PersonAuditResult
	err := json.Unmarshal(msg, &par)
	if err != nil {
		log.Println(err, "unmarshal person audit result fail")
		return false
	}

	//更新审核记录状态、更新时间
	//人工审核状态转换
	audStat := bucket.ObsAudStat[par.Status]
	sql := "UPDATE audit_message_person SET audit_status=?, update_time=? WHERE message_id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err, "upd prepare fail")
		return false
	}
	rst, err := stmt.Exec(audStat, time.Now().Unix(), par.MsgId)
	if err != nil {
		log.Println(err, "upd audit status exec fail(stmt.exec)")
		return false
	}
	stmt.Close()
	//行记录
	rn, err := rst.RowsAffected()
	if err != nil || rn == 0 {
		log.Println(err, "upd audit status fail(update none row)")
		return false
	}
	log.Printf("update rows num: %d", rn)

	//send msg to soa
	tk.sendBackMsg(par.MsgId, "audit_message_person")

	return true
}

//从db查出messageId信息组装好消息，推送消息给SOA mq
func (tk *ConsumeTask) sendBackMsg(msgId int64, msgTable string) {
	db := mydb.DB

	//select info & return to soa_back_msg
	sql := `
SELECT
  m.site_code,
  m.business_uuid,
  m.audit_mark,
  m.audit_status,
  m.update_time,
  	COALESCE(r.audit_explain, ""),
	COALESCE(r.user_id, 0),
   	COALESCE(r.username, "")
FROM %s as m LEFT JOIN audit_record AS r USING(message_id)
WHERE m.message_id = ? ORDER BY r.id desc LIMIT 1;`

	sql = fmt.Sprintf(sql, msgTable)
	rows := db.QueryRow(sql, msgId)

	//scan rows
	var bkMsg rabbit.AuditBackMsg
	var audStat, audUid int
	var audUser, audDesc string

	err := rows.Scan(
		&bkMsg.SiteCode,
		&bkMsg.BussUuid,
		&bkMsg.AuditMark,
		&audStat,
		&bkMsg.AuditTime,
		&audDesc,
		&audUid,
		&audUser,
	)
	if err != nil {
		log.Println(err, "rows scan fail")
		return
	}
	bkMsg.AuditStatus = bucket.SoaAudStat[audStat]
	bkMsg.AuditRemark = GetAdStatDesc(audStat, audDesc)
	bkMsg.AuditUid = GetAdUid(audUid, 0)
	bkMsg.AuditUser = GetAdUser(audUser, "系统")

	b, err := json.Marshal(bkMsg)
	if err != nil {
		log.Println(err, "marshal result msg data fail")
		return
	}
	log.Printf("audit back msg : %s", b)

	//msg return
	tk.MqGbVh.Publish(config.QueName["SOA_AUDIT_BACK_MSG"], b, 1)
}

//审核状态描述
func GetAdStatDesc(audStat int, defDesc string) string {
	var desc string
	switch audStat {
	default:
		desc = defDesc
	case bucket.AutoPass:
		desc = "系统审核自动通过"
	case bucket.AutoReject:
		desc = "系统审核自动拒绝"
	case bucket.NMatchAutoPass:
		desc = "规则不匹配，系统自动通过"
	}

	return desc
}

//审核Uid
func GetAdUid(uid int, defUid int) int {
	if uid > 0 {
		return uid
	}
	return defUid
}

//审核人
func GetAdUser(user string, defUser string) string {
	if user != "" {
		return user
	}
	return defUser
}
