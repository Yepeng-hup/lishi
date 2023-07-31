package cmd

import (
	"log"
	"strconv"
	"strings"
	"sysmons/core"
)

func (s *System)ProcessCheckNum(processName string, processNum int)*Process{
	cmd := "ps -ef|awk -F ' ' '{print $8}'|grep "+processName+"|wc -l"
	rel, err := core.RunCommand(cmd)
	if err != nil {
		log.Fatal(err.Error())
	}

	num, err := strconv.Atoi(strings.Replace(rel, "\n", "", -1))
	if err != nil {
		log.Fatal(err.Error())
	}
	if num == processNum {
		p := &Process{
			Response: true,
		}
		return p
	}else {
		p := &Process{
			Response: false,
		}
		return p
	}
}
