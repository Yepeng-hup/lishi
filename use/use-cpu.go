package use


import (
	"fmt"
	"log"
	"sysmons/cmd"
	"sysmons/core"
	"sysmons/curb"
)

func UseCPU(s *cmd.System, a *cmd.Alarm, c *cmd.CmdConfig){
	if s.CpuSy() < *c.CPU {
		_, err, b := curb.ReadTxtDbData("cpu")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if !b {
			curb.WriteTxtDbData("cpu", float64(*c.CPU))
			alarmInfo := f1 + a.T + f2 + fmt.Sprintf("可用CPU少于: %v%s", *c.CPU, x) + "\n\n" + fmt.Sprintf("当前可用CPU: %v%s", s.CpuSy(), x)
			if *c.DDToken == "" {
				token, err := core.CatFile(*c.DDTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				if err := core.DingDing(alarmInfo, token); err == nil {
					core.CmdLogs("Cpu utilization rate high send dingding success！")
				}else {
					log.Print(err.Error())
				}
			} else {
				if err := core.DingDing(alarmInfo, *c.DDToken); err == nil {
					core.CmdLogs("Cpu utilization rate high send dingding success！")
				}else{
					log.Print(err.Error())
				}
			}
		}
	} else {
		b, err := curb.DeleteTxtDbData("cpu")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b {
			core.CmdLogs("cpu delete value success.")
		}
		logs := fmt.Sprintf("cpuFree: %d%s", s.CpuSy(), x)
		core.CmdLogs(logs)
	}
}