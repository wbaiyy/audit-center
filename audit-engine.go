package main

import (
	"audit_engine/config"
	"audit_engine/mydb"
	"audit_engine/rabbit"
	"audit_engine/task"
	"log"
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
	mydb.DB = mydb.Connect(cfg.Mysql)

	//task 执行
	running()
}

//基于命令行的指令执行对应任务
func running() {
	// 队列连接
	var mq rabbit.MQ
	if cmd.QName == config.QueName["SOA_AUDIT_MSG"] {
		mq = tk.MqSoaVh
	} else {
		mq = tk.MqGbVh
	}

	// 队列创建
	q := mq.Create(cmd.QName)

	// 任务分派
	switch {
	case cmd.Pub:
		log.Println("message publish to queue:", q.Name)
		mq.Publish(q.Name, SimulateData(q.Name, &cmd), cmd.RepNumber)
	case cmd.Cus:
		log.Println("message consume from queue:", q.Name)
		mq.Consume(q.Name, tk.GetWork(q.Name, cmd.T), cmd.NoAck)
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
