package rabbit

type QueueName map[string]string
type QueueWork map[string]func([]byte)

//-----------------
//soa
//-----------------
//审核响应消息结构
type AuditBackMsg struct {
	SiteCode    string `json:"siteCode"`    //站点
	BussUuid    string `json:"bussUuid"`    //唯一业务ID(default uuid=0)
	AuditMark   string `json:"auditMark"`   //消息审核模板类型
	AuditStatus int    `json:"auditStatus"` //审核状态(1.审核中 2.审核通过 3.审核不通过)
	AuditRemark string `json:"auditRemark"` //审核备注
	AuditUid    int    `json:"auditUid"`    //审核人ID
	AuditUser   string `json:"auditUser"`   //审核人（默认为系统）
	AuditTime   int    `json:"auditTime"`   //发起时间
}

//审核消息结构
type AuditMsg struct {
	SiteCode   string //站点
	AuditMark  string //消息审核模板类型
	BussUuid   string `json:"bussUuid"` //唯一业务ID
	BussData   string `json:"bussData"` //业务审核基础数据
	Module     string //消息来源模块
	CreateUid  int    //消息创建者UID
	CreateUser string //消息创建者用户
	CreateTime int    //消息创建时间
	Remark     string // 用户消息备注
}

//业务数据结构
type BusinessData struct {
	GoodsBusinessData
	CouponBusinessData
}

//价格审核业务数据结构
type GoodsBusinessData struct {
	CalculatePrice float64 //计算结果价格
	ChargePrice    float64 //计费价格
	PriceLoss      float64 //亏损金额，SOA无此数据，需审核中心计算
	CatId          int     //分类ID
	VirWhCode      string  //销售仓库
	PipelineCode   string  //网站渠道
	GoodSn         string  //SKU
	SysLabelId     int     //价格系统标签（价格类型）
	Rate           float64 //利润率
	ChangeType     int     //价格变更类型：1：人工 2：系统
	SaleMark       int     //销售标识：1：正常 3：清仓
	FreightPrice   float64 //包邮运费
}

//COUPON审核业务数据
type CouponBusinessData struct {
	TemplateId          int		//COUPON模板
	PerUserReceiveCount int		//单用户领取次数
	LimitCount          int		//总使用次数
	UserLimitCount      int		//单用户使用次数
	IncludeGoodsCount   int		//适用SKU数
	FullAmount          float64	//满金额
	FullCount           int		//满数量
	ReducePercent       float64 //减百分比
	ReduceAmount        float64 //减金额
	ReduceCount         int		//减数量
	FixedPrice          float64 //一口价成本
}

//-----------------
//obs
//-----------------

//人工审核通过数据结构 //{"message_id":"1","status":"2"}
type PersonAuditResult struct {
	MsgId  int64 `json:"message_id"`
	Status int
}

//人工撤销审核数据结构 //{"message_id":"1","status":"2"}
type PersonRevokeAudit struct {
	SiteCode     string      //站点编码
	BussUuid     interface{} //唯一业务ID(default uuid=0)
	RevokeRemark string      //撤销备注
	RevokeUid    int         //审核人ID
	RevokeUser   string      //审核人名
	RevokeTime   int         //撤销时间
	RevokeMark   string      // 撤销模板类型  ‘goods-price-check’ => '价格审核' ‘promotion-coupon-check’ => 'coupon审核'
}
