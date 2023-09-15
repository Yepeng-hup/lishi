package core

import (
	"log"
	"os"
	"time"
)

var t = time.Now()

func CmdLogs(logText string)bool{
	file, err := os.OpenFile("cmd.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file error: ",err.Error())
		return false
	}
	defer file.Close()
	writeText := t.Format("2006-01-02 15:04:05")+" --> "+logText+"\n"
	if _, err := file.WriteString(writeText); err != nil {
		log.Println("write file error: ",err.Error())
		return false
	}
	return true
}
