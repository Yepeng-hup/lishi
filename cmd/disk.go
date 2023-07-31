package cmd

import (
	"log"
	"strconv"
	"strings"
	"sysmons/core"
)

func (s *System)DiskSy(dataDir string)*DiskInfo {
	cmd1 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $5}'|head -c -2"
	diskused, err := core.RunCommand(cmd1)
	if err != nil {
		log.Fatal(err.Error())
	}
	diskuseds,err := strconv.Atoi(strings.Replace(diskused, "\n", "", -1))
	if err != nil{
		log.Fatal("diskfree str change error: ", err)
	}


	cmd2 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $4}'"
	diskfree, err := core.RunCommand(cmd2)
	if err != nil {
		log.Fatal(err.Error())
	}
	diskfrees,err := strconv.Atoi(strings.Replace(diskfree, "\n", "", -1))
	if err != nil{
		log.Fatal("diskfree str change error: ", err)
	}
	diskFreeRel := core.Makes(diskfrees)

	cmd3 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $2}'"
	disktotal, err := core.RunCommand(cmd3)
	if err != nil {
		log.Fatal(err.Error())
	}
	disktotals, err := strconv.Atoi(strings.Replace(disktotal, "\n", "", -1))
	if err != nil{
		log.Fatal("disktotal str change error: ", err)
	}
	diskTotalRel := core.Makes(disktotals)

	disk := &DiskInfo{
		DiskUsed: diskuseds,
		DiskFree: diskFreeRel,
		DiskTotal: diskTotalRel,
	}
	return disk
}

