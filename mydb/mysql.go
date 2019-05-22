package mydb

import (
	"audit-center/tool"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

type Config struct {
	Host        string
	Port        int
	User        string
	Pass        string
	Protocol    string
	DbName      string
	ConnMaxLife int
}

var DB *sql.DB

func Connect(dbcf Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbcf.User, dbcf.Pass, dbcf.Protocol, dbcf.Host, dbcf.Port, dbcf.DbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		tool.FatalLog(err, "connect to mysql fail")
	}

	//最大打开的连接数100
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	//连接的最大生命周期
	log.Println("db connection max alive time:", time.Duration(dbcf.ConnMaxLife)*time.Second)
	db.SetConnMaxLifetime(time.Duration(dbcf.ConnMaxLife) * time.Second)
	return db
}

func Close(db sql.DB) {
	db.Close()
}

//连接
func Concat(ids []interface{}) string {
	return strings.Join(strings.Split(strings.Repeat("?", len(ids)), ""), ",")
}
