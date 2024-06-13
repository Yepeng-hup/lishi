package use

import (
	"log"
	"sysmons/cmd"
	"sysmons/core"
	"sysmons/curb"
)

func UseProcess(s *cmd.System, a *cmd.Alarm, c *cmd.CmdConfig) {
	if !s.ProcessCheckNum(*c.ProcessName, *c.ProcessNum).Response {
		_, err, b := curb.ReadTxtDbData("process")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if !b {
			curb.WriteTxtDbData("process", float64(*c.ProcessNum))
			alarmInfo := f1 + a.T + f2 + "have process Exit."
			if *c.DDToken == "" {
				token, err := core.CatFile(*c.DDTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				if err := core.DingDing(alarmInfo, token); err == nil {
					core.CmdLogs(*c.ProcessName + " Process run num high send dingding success!")
				}else{
					core.CmdLogs(err.Error())
				}
			} else {
				if err := core.DingDing(alarmInfo, *c.DDToken); err == nil {
					core.CmdLogs(*c.ProcessName + " Process run num high send dingding success!")
				}else{
					core.CmdLogs(err.Error())
				}
			}
		}
	} else {
		b, err := curb.DeleteTxtDbData("process")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b {
			core.CmdLogs("process delete value success.")
		}
		core.CmdLogs("process run ok.")
	}
}
