package bucket

// 10：操作人撤销,
// 20：规则引擎校验，自动通过,21：规则引擎校验，自动拒绝,22：规则全不匹配，自动通过,
// 30：人工审核中,31：人工审核通过,32：人工审核驳回

const (
	//result
	Pass   = 2
	Reject = 3

	//cancel
	ApplyCancel = 10

	//engine
	AutoPass       = 20
	AutoReject     = 21
	NMatchAutoPass = 22

	//audit
	Auditing    = 30
	AuditPass   = 31
	AuditReject = 32
)

//处理状态，1通过，2驳回
var ObsAudStat = map[int]int{
	1: AuditPass,
	2: AuditReject,
}

//soa最终状态
var SoaAudStat = map[int]int{
	//拒绝
	ApplyCancel: Reject,
	AutoReject:  Reject,
	AuditReject: Reject,
	//通过
	AutoPass:       Pass,
	AuditPass:      Pass,
	NMatchAutoPass: Pass,
}
