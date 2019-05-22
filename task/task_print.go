package task

import "log"

func (tk *ConsumeTask) workPrintMessage(msg []byte) bool {
	log.Println("working...")
	log.Println(string(msg))
	log.Println("done")
	return true
}
