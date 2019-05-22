package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
)

type CmdArgs struct {
	//engine
	V bool
	H bool
	//test
	T     bool
	NoAck bool
	//config file
	Cfg string
	//mq
	QName     string
	Pub       bool
	Cus       bool
	MsgData   string
	RepNumber int
	WorkerNumber int
}

func (cmd *CmdArgs) Parse() {
	//engine
	pflag.BoolVarP(&cmd.V, "version", "V", false, "version info")
	pflag.BoolVarP(&cmd.H, "help", "h", false, "show this message")
	//test
	pflag.BoolVarP(&cmd.T, "test", "t", false, `test model for msg publish or consume`)
	pflag.BoolVar(&cmd.NoAck, "no_ack", false, `do not return ack the message to the Broker`)
	//get config file
	pflag.StringVarP(&cmd.Cfg, "config_file", "c", "./config.json", "config file")
	//mq
	pflag.BoolVar(&cmd.Pub, "publish", false, "run publish message")
	pflag.BoolVar(&cmd.Cus, "consume", false, "run consume message")
	pflag.StringVarP(&cmd.QName, "queue", "q", "", `message queue name(`+availQueue()+`)`)
	pflag.StringVarP(&cmd.MsgData, "msg_data", "D", "", `message data used for publish`)
	pflag.IntVarP(&cmd.RepNumber, "rep_number", "n", 1, `message data send rep_number time`)
	pflag.IntVarP(&cmd.WorkerNumber, "worker_number", "w", 1, `worker that work message`)
	pflag.Parse()

	//show info
	cmd.ShowInfo()
}

//show all info
func (cmd *CmdArgs) ShowInfo() {
	switch {
	case cmd.V:
		cmd.PrintVersion()
		os.Exit(0)
	case cmd.H:
		pflag.PrintDefaults()
		os.Exit(0)
	}
}

//engine show
func (cmd *CmdArgs) PrintVersion() {
	log.Println(fmt.Sprintf("%s, %s", "Audit Rule Engine", "v0.0.1 beta"))
}

//queue list
func availQueue() string {
	var q []string
	for _, v := range QueName {
		q = append(q, v)
	}
	return strings.Join(q, "|")
}
