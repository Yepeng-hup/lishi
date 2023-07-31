package cmd

import (
	"log"
	"strconv"
	"strings"
	"sysmons/core"
)


func (s *System) MemSy() *MemInfo {
	memtotal, err := core.RunCommand(`free | awk '{print $2}'|awk 'NR==2 {print}'`)
	if err != nil {
		log.Fatal(err.Error())
	}
	memtotals,err := strconv.Atoi(strings.Replace(memtotal, "\n", "", -1))
	if err != nil{
		log.Fatal("memTotal str change error: ", err)
	}
	totalRel := core.Makes(memtotals)

	memused, err := core.RunCommand(`free | awk '{print $3}'|awk 'NR==2 {print}'`)
	if err != nil {
		log.Fatal(err.Error())
	}
	memuseds,err := strconv.Atoi(strings.Replace(memused, "\n", "", -1))
	if err != nil{
		log.Fatal("memUsed str change error: ", err)
	}
	usedRel := core.Makes(memuseds)

	memfree, err := core.RunCommand(`free | awk '{print $4}'|awk 'NR==2 {print}'`)
	if err != nil {
		log.Fatal(err.Error())
	}
	memfrees,err := strconv.Atoi(strings.Replace(memfree, "\n", "", -1))
	if err != nil{
		log.Fatal("memFree str change error: ", err)
	}
	freeRel := core.Makes(memfrees)

	mem := &MemInfo{
		MemTotal: totalRel,
		MemUsed: usedRel,
		MemFree: freeRel,
	}
	return mem
}
