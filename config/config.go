package config

import (
	"audit-center/mydb"
	"audit-center/rabbit"
	"audit-center/tool"
	"github.com/spf13/viper"
	"log"
)

type EngineInfo struct {
	Name    string
	Version string
}

type CFG struct {
	cmd        CmdArgs
	Test       bool
	ConfigFile string
	RabbitMq   map[string]rabbit.Config
	Mysql      map[string]mydb.Config
}

//config init
func (cfg *CFG) InitByCmd(cmd CmdArgs) {
	//read config file
	viper.SetConfigFile(cmd.Cfg)
	err := viper.ReadInConfig()
	tool.FatalLog(err, "viper read config error")

	//test
	cfg.cmd = cmd
	cfg.Test = cmd.T
	cfg.ConfigFile = cmd.Cfg

	//init rabbitmq config
	cfg.RabbitMq = make(map[string]rabbit.Config)
	for _, v := range []string{"soa", "gb", "obs"} {
		cfg.RabbitMq[v] = rabbit.Config{
			Host:  viper.GetString("rabbitmq." + v + ".host"),
			Port:  viper.GetInt("rabbitmq." + v + ".port"),
			User:  viper.GetString("rabbitmq." + v + ".user"),
			Pass:  viper.GetString("rabbitmq." + v + ".pass"),
			Vhost: viper.GetString("rabbitmq." + v + ".vhost"),
		}
	}

	//init mysql config
	cfg.Mysql = make(map[string]mydb.Config)
	for _, v := range []string{"obs", "goods"}  {
		cfg.Mysql[v] = mydb.Config{
			Host:        viper.GetString("mysql." + v + ".host"),
			Port:        viper.GetInt("mysql." + v + ".port"),
			User:        viper.GetString("mysql." + v + ".user"),
			Pass:        viper.GetString("mysql." + v + ".pass"),
			DbName:      viper.GetString("mysql." + v + ".dbname"),
			Protocol:    viper.GetString("mysql." + v + ".protocol"),
			ConnMaxLife: viper.GetInt("mysql." + v + ".conn_max_life"),
		}
	}

	//print env info
	cfg.printEnv()
}

func (cfg *CFG) printEnv() {
	cfg.cmd.PrintVersion()
	log.Printf("cmdline: %+v\n", cfg.cmd)
}
