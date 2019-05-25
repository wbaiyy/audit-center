package task

import (
	"audit-center/cache"
	"audit-center/mydb"
	"audit-center/tool"
	"database/sql"
	"fmt"
	"log"
	"time"
)

//规则哈希表
type AuditTypeList map[string]AuditType

//规则类型
type AuditType struct {
	TypeId    int
	TypeTitle string
	AuditSort int
	AuditMark string
	RuleList  []AuditRule
}

//规则条目
type AuditRule struct {
	RuleId   int
	TypeId   int
	RuleRel  int //1与 2或
	RuleProc int //rule成立后的处理方式，1 系统通过，2 系统驳回，3 转人工审核
	FlowId   int
	Profit   float64
	ItemList []RuleItem
}

type RuleItem struct {
	ItemId      int
	RuleId      int
	CompareType int
	Field       string
	Operate     string
	Value       string
}

//规则项(compare_type 1:阈值 2:字段）
func (tk *ConsumeTask) GetRuleItems(auditMark string) AuditTypeList {
	db := mydb.DB

	//---------审核类型
	sql := fmt.Sprintf(`select id, title, sort,audit_mark from audit_template where audit_mark = "%s";`, auditMark)
	rows, err := db.Query(sql)
	if err != nil {
		tool.FatalLog(err, "SELECT audit_template")
	}

	var aTypes []AuditType
	var typeIds []interface{}
	for rows.Next() {
		var at AuditType
		rows.Scan(&at.TypeId, &at.TypeTitle, &at.AuditSort, &at.AuditMark)
		aTypes = append(aTypes, at)
		//审核ID
		typeIds = append(typeIds, at.TypeId)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("aTypes:\n", aTypes)
	rows.Close()

	//----------规则条目
	if len(typeIds) == 0 {
		return nil
	}
	stmt, err := db.Prepare("select id ,template_id, items_relation, process_type, workflow_id, base_profit_margin " +
		"from audit_rule WHERE template_id IN (" + mydb.Concat(typeIds) + ") " +
		"ORDER BY sort ASC;")
	if err != nil {
		log.Fatal(err)
	}
	rows, err = stmt.Query(typeIds...)
	if err != nil {
		log.Fatal(err)
	}

	var aRules []AuditRule
	var rids []interface{}
	var ar AuditRule
	ruleGroups := make(map[int][]AuditRule, len(aTypes))
	for rows.Next() {
		rows.Scan(&ar.RuleId, &ar.TypeId, &ar.RuleRel, &ar.RuleProc, &ar.FlowId, &ar.Profit)
		aRules = append(aRules, ar)
		//规则条目 list
		rids = append(rids, ar.RuleId)

		//分组
		ruleGroups[ar.TypeId] = append(ruleGroups[ar.TypeId], ar)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("aRules:\n", aRules)
	//log.Println("ruleGroups:\n", ruleGroups)
	rows.Close()
	stmt.Close()

	//分组+哈希表填充
	hashAuditTypeList := make(AuditTypeList, 10)
	for i, at := range aTypes {
		aTypes[i].RuleList = ruleGroups[at.TypeId]
		hashAuditTypeList[at.AuditMark] = aTypes[i]
	}
	//log.Println("hashAuditTypeList:\n", hashAuditTypeList)

	//--------比较项
	if len(rids) == 0 {
		return nil
	}
	sql = "select id, rule_id,compare_type,field,operation,value from audit_rule_item WHERE rule_id IN (" + mydb.Concat(rids) + ")"
	stmt, err = db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err = stmt.Query(rids...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []RuleItem
	itemGroups := make(map[int][]RuleItem, len(aRules))
	for rows.Next() {
		var k RuleItem
		rows.Scan(&k.ItemId, &k.RuleId, &k.CompareType, &k.Field, &k.Operate, &k.Value)
		items = append(items, k)
		itemGroups[k.RuleId] = append(itemGroups[k.RuleId], k)
	}
	//log.Println("itemGroups:\n", itemGroups)

	//哈希表填充
	for k, t := range hashAuditTypeList {
		for kk, r := range t.RuleList {
			hashAuditTypeList[k].RuleList[kk].ItemList = itemGroups[r.RuleId]
		}
	}
	//log.Printf("hashAuditTypeList: %+v\n", hashAuditTypeList)

	return hashAuditTypeList
}

func getCostPrice(virWhCode string, goodSn string, chargePrice float64) float64 {
	c := cache.Storage
	key := fmt.Sprintf("goods:%s:%s", virWhCode, goodSn)
	goodsCostPrice, found  := c.Get(key)
	fmt.Println(goodsCostPrice, found)
	if found {
		return  goodsCostPrice.(float64)
	}
	dbPrice := getCostPriceFromDb(virWhCode, goodSn)
	if dbPrice == 0 {
		return chargePrice
	}

	c.Set(key, dbPrice , 2 * time . Minute)
	return dbPrice
}

func getCostPriceFromDb(virWhCode string, goodSn string) float64{
	db := mydb.GoodsDB
	var goodsCostPrice float64
	err := db.QueryRow(
		`SELECT shiji_price FROM public_ods_gb_v_purchase_bill_base_months WHERE v_wh_code=? AND sku=?`, virWhCode, goodSn ).
		Scan(&goodsCostPrice)

	switch {
	case err == sql.ErrNoRows:
		return 0
	case err != nil:
		tool.FatalLog(err, "QueryRow [table->public_ods_gb_v_purchase_bill_base_months] error:")
	}

	return goodsCostPrice
}
