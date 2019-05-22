package task

import (
	"audit-center/rabbit"
	"fmt"
	"log"
)

const (
	RuleMatched    = 1
	RuleNotMatched = 2
)

//1与 2或
const (
	RelAnd = 1
	RelOr  = 2
)

const (
	SysPass    = 1 //系统匹配通过
	SysReject  = 2 //系统拒绝
	ObsAudit   = 3 //obs审核
	SysDefPass = 4 //系统未匹配通过
)

var AuditStatus = map[int]int{
	SysPass:    20, //规则引擎校验，自动通过
	SysReject:  21, //规则引擎校验，自动拒绝
	SysDefPass: 22, //规则全不匹配，自动通过
	ObsAudit:   30, //人工审核中
}

//规则匹配结果
type RuleMatch struct {
	RMatch      bool
	RuleId      int
	FlowId      int
	RuleGo      int
	Profit      float64
	Explain     string
	ItemMatches []ItemMatch
}

//规则项匹配结果
type ItemMatch struct {
	ItemId  int
	IMatch  bool
	Explain string
}

//bussData 转成对应项的string值
func bussDataToString(field string, bussData *rabbit.BusinessData, baseRate float64) string {
	switch field {
	//商品价格审核字段
	case "catId":
		return fmt.Sprintf("%d", bussData.CatId)
	case "changeType":
		return fmt.Sprintf("%d", bussData.ChangeType)
	case "chargePrice":
		return fmt.Sprintf("%0.4f", bussData.ChargePrice)
	case "pipelineCode":
		return bussData.PipelineCode
	case "priceLoss":
		return GetPriceLoss(bussData.ChargePrice, bussData.Rate, baseRate)
	case "rate":
		return fmt.Sprintf("%0.4f", bussData.Rate)
	case "sysLabelId":
		return fmt.Sprintf("%d", bussData.SysLabelId)
	case "virWhCode":
		return bussData.VirWhCode
	case "saleMark":
		return fmt.Sprintf("%d", bussData.SaleMark)
	//COUPON审核字段
	case "templateId":
		return fmt.Sprintf("%d", bussData.TemplateId)
	case "perUserReceiveCount":
		return fmt.Sprintf("%d", bussData.PerUserReceiveCount)
	case "limitCount":
		return fmt.Sprintf("%d", bussData.LimitCount)
	case "userLimitCount":
		return fmt.Sprintf("%d", bussData.UserLimitCount)
	case "includeGoodsCount":
		return fmt.Sprintf("%d", bussData.IncludeGoodsCount)
	case "fullCount":
		return fmt.Sprintf("%d", bussData.FullCount)
	case "fullAmount":
		return fmt.Sprintf("%0.4f", bussData.FullAmount)
	case "reducePercent":
		return fmt.Sprintf("%0.4f", bussData.ReducePercent)
	case "reduceAmount":
		return fmt.Sprintf("%0.4f", bussData.ReduceAmount)
	case "reduceCount":
		return fmt.Sprintf("%d", bussData.ReduceCount)
	case "fixedPrice":
		return fmt.Sprintf("%0.4f", bussData.FixedPrice)
	}
	return "=X="
}

//get priceLoss
func GetPriceLoss(chargePrice float64, rate float64, baseRate float64) string {
	//亏损金额 = before:（利润率 - 基础利润率） × 计费价格 /  6.1, now:（ 基础利润率  - 利润率） × 计费价格 /  6.1
	return fmt.Sprintf("%0.4f", chargePrice*(baseRate-rate)/6.1)
}

//rule多条规则比较
//返回结果:
// int:	1 系统通过，2 系统驳回，3 转人工审核
// RuleMatch: 匹配的规则明细
func RunRuleMatch(bussData *rabbit.BusinessData, auditType *AuditType) (int, RuleMatch) {

	var rml []RuleMatch
	var result int

	for i, rule := range auditType.RuleList {
		var iml []ItemMatch

		//item结果
		for _, item := range rule.ItemList {
			field := bussDataToString(item.Field, bussData, rule.Profit)
			match := ValueCompare(field, item.Operate, item.Value)
			im := ItemMatch{
				ItemId:  item.ItemId,
				IMatch:  match,
				Explain: fmt.Sprintf(`(bussData.%v) [%v %v %v]`, item.Field, field, item.Operate, item.Value),
			}
			iml = append(iml, im)
		}

		//rule的验证结果
		switch rule.RuleRel {
		case RelAnd:
			for _, im := range iml {
				if !im.IMatch { //与条件，只要有一个不匹配，直接不匹配
					result = RuleNotMatched
					break
				}
				result = RuleMatched
			}
		case RelOr:
			for _, im := range iml {
				if im.IMatch { //或条件，只要有一个匹配，直接匹配
					result = RuleMatched
					break
				}
				result = RuleNotMatched
			}
		}

		//基于规则引擎校验的结果进行进一步处理
		rml = append(rml, RuleMatch{
			RMatch:      result == RuleMatched,
			RuleId:      rule.RuleId,
			FlowId:      rule.FlowId,
			Profit:      rule.Profit,
			RuleGo:      rule.RuleProc,
			Explain:     fmt.Sprintf("itemsRel=%d (1:and 2:or)", rule.RuleRel),
			ItemMatches: iml,
		})

		log.Printf("rule[%d]: %+v\n", i, rml[i])

		if result == RuleMatched { //任一条rule通过，则按Rule Process处理
			return AuditStatus[rule.RuleProc], rml[len(rml)-1]
		}
	}

	//如果都不匹配，默认规则放行
	return AuditStatus[SysDefPass], RuleMatch{}
}
