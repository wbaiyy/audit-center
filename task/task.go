package task

import (
	"audit_engine/config"
	"audit_engine/rabbit"
	"log"
)

type ConsumeTask struct {
	TkCfg   config.CFG
	MqGbVh  rabbit.MQ //gb vhost
	MqSoaVh rabbit.MQ //soa vhost
}

//初始化队列任务环境
func (tk *ConsumeTask) Bootstrap() {
	log.Println("task bootstrap...")

	//初始化一个rabbit连接
	cfg := tk.TkCfg
	tk.MqSoaVh.Init(cfg.RabbitMq["soa"])
	tk.MqGbVh.Init(cfg.RabbitMq["gb"])
}

//停止则回收相关资源
func (tk *ConsumeTask) Stop() {
	log.Println("task clean...")

	tk.MqSoaVh.Close()
	tk.MqGbVh.Close()
}

//基于queue队列名分配工作任务
func (tk *ConsumeTask) GetWork(qn string, test bool) (workFn func([]byte) bool) {
	switch {
	case test:
		workFn = tk.workPrintMessage
	case qn == config.QueName["SOA_AUDIT_MSG"]:
		workFn = tk.workAuditMessage
	case qn == config.QueName["OBS_PERSON_AUDIT_RESULT"]:
		workFn = tk.workUpdateAuditResult
	}
	return workFn
}
