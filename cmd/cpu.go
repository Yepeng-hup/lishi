package cmd

import (
	"log"
	"strconv"
	"strings"
	"sysmons/core"
)

func (s *System) CpuSy()int{
	cpuFree, err := core.RunCommand(`top -b -n 1|grep -w "Cpu"|awk -F ',' '{print $4}'|cut -f 1 -d "."`)
	if err != nil {
		log.Fatal(err.Error())
	}
	cpuFrees, err := strconv.Atoi(strings.Replace(strings.Replace(cpuFree, "\n", "", -1), " ", "", -1))
	if err != nil{
		log.Fatal("cpu str change error: ", err)
	}
	return cpuFrees
}

