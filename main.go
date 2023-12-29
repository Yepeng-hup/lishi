package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sysmons/cmd"
	"sysmons/core"
	"sysmons/curb"
)

const(
	f1 = "*** "
	f2 = " ***\n\n"
	x = "%"
	runTimeDay = 30
)

func main(){
	if len(os.Args) < 2 {
		fmt.Println("Please enter parameters or -h to view the usage of parameters.")
		os.Exit(1)
	}

	diskDir := flag.String("diskDataDir", "/", "Specify the storage directory to monitor. The default is / .")
	diskHorizon := flag.Int("d", 50, "Specify how many utilization of the disk to send an alarm,The default is 80%.")
	memHorizon := flag.Float64("m", 2.0, "Specify the number of gigabytes of available memory to send an alarm. The default is 2.0G.")
	cleCacheNum := flag.Int("c", 0, "Specify the number[1,2,3] clear system cache.[0] do nothing.The default is 0 .")
	cpu := flag.Int("cpu", 20, "Specify how much the CPU is lower than to send an alarm, 20% by default.")
	ddToken := flag.String("token", "", "Specify the Token to send DingDing.")
	ddTokenFile := flag.String("token_filePath", "", "File path with token written.")
	title := flag.String("t", "", "push DindDing Keyword Title.")
	processNum := flag.Int("p", 1, "Count the total number of processes.The default is 1 .")
	processName := flag.String("n", "nil", "Process name supports wildcard.The default is nil.")
	//suppression := flag.Int("s", 1, "Is alarm suppression and noise reduction enabled[0,1].The default is 1 .")

	flag.Parse()

	d := &cmd.Diskd{
		DiskDir: *diskDir,
	}
	a := &cmd.Alarm{
		T: *title,
	}
	s := cmd.System{}

	if s.MemSy().MemFree < *memHorizon {
		_, err, b := curb.ReadTxtDbData("mem")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b != true {
			curb.WriteTxtDbData("mem", *memHorizon)
			alarmInfo := f1+a.T+f2+fmt.Sprintf("可用内存少于: %vG", *memHorizon)+"\n\n"+fmt.Sprintf("当前可用内存: %.2fG", s.MemSy().MemFree)
			if *ddToken == "" {
				token, err := core.CatFile(*ddTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				er := core.DingDing(alarmInfo, token)
				if er == nil {
					core.CmdLogs("mem detonate send dingding success!")
				}else {
					log.Print(er.Error())
				}
			}else {
				err := core.DingDing(alarmInfo, *ddToken)
				if err == nil {
					core.CmdLogs("mem detonate send dingding success!")
				}else {
					log.Print(err.Error())
				}
			}
			err := cmd.CleCache(*cleCacheNum)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}else {
		b, err := curb.DeleteTxtDbData("mem")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b == true{
			core.CmdLogs("mem delete value success.")
		}
		logs := fmt.Sprintf("memFree: %.2fG", s.MemSy().MemFree)
		core.CmdLogs(logs)
	}

	if int(s.DiskSy(d.DiskDir).DiskFree) < *diskHorizon {
		_, err, b := curb.ReadTxtDbData("disk")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b != true {
			curb.WriteTxtDbData("disk", float64(*diskHorizon))
			alarmInfo := f1 + a.T + f2 + fmt.Sprintf("`%s` 磁盘空间少于: %v%s", d.DiskDir, *diskHorizon, "G") + "\n\n" + fmt.Sprintf("`%s` 当前可用空间: %.2fG", d.DiskDir, s.DiskSy(*diskDir).DiskFree)
			if *ddToken == "" {
				token, err := core.CatFile(*ddTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				er := core.DingDing(alarmInfo, token)
				if er == nil {
					core.CmdLogs("Disk detonate send dingding success!")
				} else {
					log.Print(er.Error())
				}
			} else {
				err := core.DingDing(alarmInfo, *ddToken)
				if err == nil {
					core.CmdLogs("Disk detonate send dingding success!")
				} else {
					log.Print(err.Error())
				}
			}
		}
	}else{
		b, err := curb.DeleteTxtDbData("disk")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b == true{
			core.CmdLogs("disk delete value success.")
		}
		logs := fmt.Sprintf("diskFree: %.2fG",s.DiskSy(d.DiskDir).DiskFree)
		core.CmdLogs(logs)
	}

	if s.CpuSy() < *cpu {
		_, err, b := curb.ReadTxtDbData("cpu")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b != true {
			curb.WriteTxtDbData("cpu", float64(*cpu))
			alarmInfo := f1 + a.T + f2 + fmt.Sprintf("可用CPU少于: %v%s", *cpu, x) + "\n\n" + fmt.Sprintf("当前可用CPU: %v%s", s.CpuSy(), x)
			if *ddToken == "" {
				token, err := core.CatFile(*ddTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				er := core.DingDing(alarmInfo, token)
				if er == nil {
					core.CmdLogs("cpu利用率爆炸发送钉钉成功！")
				} else {
					log.Print(er.Error())
				}
			} else {
				err := core.DingDing(alarmInfo, *ddToken)
				if err == nil {
					core.CmdLogs("cpu利用率爆炸发送钉钉成功！")
				} else {
					log.Print(err.Error())
				}
			}
		}
	}else{
		b, err := curb.DeleteTxtDbData("cpu")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b == true{
			core.CmdLogs("cpu delete value success.")
		}
		logs := fmt.Sprintf("cpuFree: %d%s",s.CpuSy(), x)
		core.CmdLogs(logs)
	}


	if s.ProcessCheckNum(*processName, *processNum).Response != true {
		_, err, b := curb.ReadTxtDbData("process")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b != true {
			curb.WriteTxtDbData("process", float64(*processNum))
			alarmInfo := f1 + a.T + f2 + fmt.Sprintf("have process Exit.")
			if *ddToken == "" {
				token, err := core.CatFile(*ddTokenFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				er := core.DingDing(alarmInfo, token)
				if er == nil {
					core.CmdLogs(*processName + " process run num detonate send dingding success!")
				} else {
					log.Print(er.Error())
				}
			} else {
				err := core.DingDing(alarmInfo, *ddToken)
				if err == nil {
					core.CmdLogs(*processName + " process run num detonate send dingding success!")
				} else {
					log.Print(err.Error())
				}
			}
		}
	}else {
		b, err := curb.DeleteTxtDbData("process")
		if err != nil {
			core.CmdLogs(err.Error())
		}
		if b == true{
			core.CmdLogs("cpu delete value success.")
		}
		core.CmdLogs("process run ok.")
	}
}
