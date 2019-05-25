package main

import (
	"audit-center/cache"
	"audit-center/config"
	"audit-center/mydb"
	"audit-center/rabbit"
	"audit-center/task"
	"log"
	"strings"
	"sync"
)

var (
	cmd config.CmdArgs
	cfg config.CFG
	tk  task.ConsumeTask
)

func main() {
	//cmdline 解析
	cmd.Parse()

	//config 获取
	cfg.InitByCmd(cmd)

	//mq 初始化
	tk = task.ConsumeTask{TkCfg: cfg}
	tk.Bootstrap()
	defer tk.Stop()

	//mysql 初始化
	mydb.DB = mydb.Connect(cfg.Mysql["obs"])
	mydb.GoodsDB= mydb.Connect(cfg.Mysql["goods"])

	//cache 初始化
	cache.Storage = cache.New()

	//task 执行
	Dispatch()
	log.Println("Audit Rule Engine END !!!")
	//running()
}

//根据队列处理
func Dispatch() {
	var wg sync.WaitGroup
	queues := strings.Split(cmd.QName, ",")
	for _, queueName := range queues  {
		if config.IsValidateQueueName(queueName) {
			wg.Add(1)
			go func(queueName string) {
				defer wg.Done()
				running(queueName)
			}(queueName)
		}
	}
	wg.Wait()
}



//基于命令行的指令执行对应任务
func running(queueName string) {
	// 队列连接
	var mq rabbit.MQ
	if config.IsSoaQueue(queueName) {
		mq = tk.MqSoaVh
	} else if config.IsGbQueue(queueName) {
		mq = tk.MqGbVh
	} else {
		mq = tk.MqObsVh
	}

	// 队列创建
	mq.GetChannel(queueName)
	q := mq.Create(queueName)

	// 任务分派
	switch {
	case cmd.Pub:
		log.Println("message publish to queue:", q.Name)
		mq.Publish(q.Name, SimulateData(q.Name, &cmd), cmd.RepNumber)
	case cmd.Cus:
		//log.Println("message consume from queue:", q.Name)
		//mq.Consume(q.Name, tk.GetWork(q.Name, cmd.T), cmd.NoAck)
		tk.RunWork(mq.Consume(q.Name), tk.GetWork(q.Name, cmd.T), cmd.WorkerNumber, cmd.NoAck, queueName)
	default:
		log.Fatalln("[x]", "queue must be consume or publish")
	}
}

// select queue and prepare message data
func SimulateData(qn string, cmd *config.CmdArgs) []byte {
	if cmd.MsgData != "" {
		return []byte(cmd.MsgData)
	}
	switch qn {
	case config.QueName["SOA_AUDIT_BACK_MSG"]:
		return rabbit.SimAuditBackMsg()
	case config.QueName["SOA_AUDIT_MSG"]:
		return rabbit.SimAuditMsg()
	case config.QueName["OBS_RULE_CHANGE_MSG"]:
		return []byte(`{"action":"upd|del|add","templete_id":1}`)
	case config.QueName["OBS_PERSON_AUDIT_RESULT"]:
		return []byte(`{"message_id":1,"status": 2}`)
	default:
		return []byte("Heyman, Cool")
	}
}
