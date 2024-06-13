package use

import (
	"fmt"
	"log"
	"sysmons/cmd"
	"sysmons/core"
	"sysmons/curb"
)


func UseDisk(d *cmd.Diskd, s *cmd.System, a *cmd.Alarm, c *cmd.CmdConfig) {
	if int(s.DiskSy(d.DiskDir).DiskFree) < *c.DiskHorizon {
		_, err, b := curb.ReadTxtDbData("disk")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if !b {
			curb.WriteTxtDbData("disk", float64(*c.DiskHorizon))
			alarmInfo := f1 + a.T + f2 + fmt.Sprintf("`%s` 磁盘空间少于: %v%s", d.DiskDir, *c.DiskHorizon, "G") + "\n\n" + fmt.Sprintf("`%s` 当前可用空间: %.2fG", d.DiskDir, s.DiskSy(*c.DiskDir).DiskFree)
			if *c.DDToken == "" {
				token, err := core.CatFile(*c.DDTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				if err := core.DingDing(alarmInfo, token); err == nil {
					core.CmdLogs("Disk utilization rate high send dingding success!")
				}else {
					core.CmdLogs(err.Error())
				}
			} else {
				if err := core.DingDing(alarmInfo, *c.DDToken); err == nil {
					core.CmdLogs("Disk utilization rate high send dingding success!")
				}else{
					core.CmdLogs(err.Error())
				}
			}
		}
	} else {
		b, err := curb.DeleteTxtDbData("disk")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b {
			core.CmdLogs("disk delete value success.")
		}
		logs := fmt.Sprintf("diskFree: %.2fG", s.DiskSy(d.DiskDir).DiskFree)
		core.CmdLogs(logs)
	}
}