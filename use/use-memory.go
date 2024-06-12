package use

import (
	"fmt"
	"log"
	"sysmons/cmd"
	"sysmons/core"
	"sysmons/curb"
)

func UseMemory(s *cmd.System, a *cmd.Alarm, c *cmd.CmdConfig)  {
	if s.MemSy().MemFree < *c.MemHorizon {
		_, err, b := curb.ReadTxtDbData("mem")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if !b {
			curb.WriteTxtDbData("mem", *c.MemHorizon)
			alarmInfo := f1 + a.T + f2 + fmt.Sprintf("可用内存少于: %vG", *c.MemHorizon) + "\n\n" + fmt.Sprintf("当前可用内存: %.2fG", s.MemSy().MemFree)
			if *c.DDToken == "" {
				token, err := core.CatFile(*c.DDTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				if err := core.DingDing(alarmInfo, token); err == nil {
					core.CmdLogs("Mem utilization rate high send dingding success!")
				}else {
					log.Print(err.Error())
				}
			} else {
				if err := core.DingDing(alarmInfo, *c.DDToken); err == nil {
					core.CmdLogs("Mem utilization rate high send dingding success!")
				}else {
					log.Print(err.Error())
				}
			}
			if err := cmd.CleCache(*c.CleCacheNum); err != nil {
				log.Println(err.Error())
			}
		}
	} else {
		b, err := curb.DeleteTxtDbData("mem")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b {
			core.CmdLogs("mem delete value success.")
		}
		logs := fmt.Sprintf("memFree: %.2fG", s.MemSy().MemFree)
		core.CmdLogs(logs)
	}
}