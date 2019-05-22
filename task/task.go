package task

import (
	"audit-center/config"
	"audit-center/rabbit"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type ConsumeTask struct {
	TkCfg   config.CFG
	MqGbVh  rabbit.MQ //gb vhost
	MqSoaVh rabbit.MQ //soa vhost
	MqObsVh rabbit.MQ //obs vhost
}

//初始化队列任务环境
func (tk *ConsumeTask) Bootstrap() {
	log.Println("task bootstrap...")

	//初始化一个rabbit连接
	cfg := tk.TkCfg
	tk.MqSoaVh.Init(cfg.RabbitMq["soa"])
	tk.MqGbVh.Init(cfg.RabbitMq["gb"])
	tk.MqObsVh.Init(cfg.RabbitMq["obs"])
}

//停止则回收相关资源
func (tk *ConsumeTask) Stop() {
	log.Println("task clean...")

	tk.MqSoaVh.Close()
	tk.MqGbVh.Close()
	tk.MqObsVh.Close()
}

//基于queue队列名分配工作任务
func (tk *ConsumeTask) GetWork(qn string, test bool) (workFn func([]byte) bool) {
	switch {
	case test:
		workFn = tk.workPrintMessage
	case qn == config.QueName["SOA_AUDIT_MSG"] || qn == config.QueName["TASK_AUDIT_MSG"]:  //审核消息
		workFn = tk.workAuditMessage
	case qn == config.QueName["OBS_PERSON_AUDIT_RESULT"]:   //人工审核消息结果
		workFn = tk.workUpdateAuditResult
	case qn == config.QueName["TASK_AUDIT_REVOKE_MSG"]:		//审核消息撤销
		workFn = tk.workAuditMessageRevoke
	}
	return workFn
}

//处理消息
func (tk *ConsumeTask) RunWork(messages <-chan amqp.Delivery, workMethod func([]byte) bool ,workerNum int, noAck bool, qn string)  {

	var wg sync.WaitGroup
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go func(i int) {
			log.Printf("==> [%s] task-%d start...", qn, i)
			defer wg.Done()
			for  {
				d, ok := <- messages
				if !ok {
					break
				}
				success := workMethod(d.Body)
				log.Printf("<== [%s] task-%d done, result: [%v]!!", qn, i, success)

				if success && !noAck {
					d.Ack(false)
				}
			}
		}(i)
	}
	wg.Wait()
}


