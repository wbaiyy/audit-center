package tool

import (
	"log"
)

//致命错误
func FatalLog(err error, msg ...string) {
	if err != nil {
		log.Fatalln("[Error]", msg, err)
	}
}
