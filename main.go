package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sysmons/cmd"
	"sysmons/core"
)

const(
	f1 = "*** "
	f2 = " ***\n\n"
	x = "%"
	runTimeDay = 30
)

func main(){
	if len(os.Args) < 2 {
		log.Fatal("Please enter parameters or -h to view the usage of parameters.")
	}

	diskDir := flag.String("diskDataDir", "/", "Specify the storage directory to monitor. The default is / .")
	diskHorizon := flag.Int("d", 80, "Specify how many utilization of the disk to send an alarm,The default is 80%.")
	memHorizon := flag.Float64("m", 2.0, "Specify the number of gigabytes of available memory to send an alarm. The default is 2.0G.")
	cleCacheNum := flag.Int("c", 3, "Specify the number[1,2,3] clear system cache. The default is 3 .")
	cpu := flag.Int("cpu", 20, "Specify how much the CPU is lower than to send an alarm, 20% by default.")
	ddToken := flag.String("token", "", "Specify the Token to send DingDing.")
	ddTokenFile := flag.String("token_filePath", "", "File path with token written.")
	title := flag.String("t", "", "push DindDing Keyword Title.")
	processNum := flag.Int("p", 1, "Count the total number of processes.The default is 1 .")
	processName := flag.String("n", "nil", "Process name supports wildcard.The default is nil.")
	flag.Parse()

	d := &cmd.Diskd{
		DiskDir: *diskDir,
	}
	a := &cmd.Alarm{
		T: *title,
	}
	s := cmd.System{}

	//i := s.ProcessCheckTime(*proccssFullName).RunTime
	//u := s.ProcessCheckNum(*processName, *processNum).Response
	//fmt.Println(u)
	//fmt.Println(i)
	//os.Exit(2)

	if s.MemSy().MemFree < *memHorizon{
		alarmInfo := f1+a.T+f2+fmt.Sprintf("Less memory available: %vG", *memHorizon)+"\n\n"+fmt.Sprintf("Available memory: %.2fG", s.MemSy().MemFree)
		if *ddToken == "" {
			token, err := core.CatFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := core.DingDing(alarmInfo, token)
			if er == nil {
				log.Print("mem detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := core.DingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print("mem detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
		err := cmd.CleCache(*cleCacheNum)
		if err != nil {
			log.Println(err.Error())
		}
	}else {
		log.Printf("memFree: %.2fG\n", s.MemSy().MemFree)
	}

	if s.DiskSy(d.DiskDir).DiskUsed > *diskHorizon{
		rel := 100 - s.DiskSy(d.DiskDir).DiskUsed
		alarmInfo := f1+a.T+f2+fmt.Sprintf("Available disks are less than `%s` : %v%s", d.DiskDir, 100-*diskHorizon, x)+"\n\n"+fmt.Sprintf("Available disks `%s`: %d%s, %.2fG",d.DiskDir, rel, x, s.DiskSy(*diskDir).DiskFree)
		if *ddToken == "" {
			token, err := core.CatFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := core.DingDing(alarmInfo, token)
			if er == nil {
				log.Print("Disk detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := core.DingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print("Disk detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
	}else{
		log.Printf("diskFree: %.2fG\n",s.DiskSy(d.DiskDir).DiskFree)
	}

	if s.CpuSy() < *cpu {
		alarmInfo := f1+a.T+f2+fmt.Sprintf("Less CPU available: %v%s", 100-*cpu, x)+"\n\n"+fmt.Sprintf("Available CPU: %v%s", s.CpuSy(), x)
		if *ddToken == "" {
			token, err := core.CatFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := core.DingDing(alarmInfo, token)
			if er == nil {
				log.Print("cpu detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := core.DingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print("cpu detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
	}else{
		log.Printf("cpuFree: %d%s\n",s.CpuSy(), x)
	}


	if s.ProcessCheckNum(*processName, *processNum).Response != true {
		alarmInfo := f1+a.T+f2+fmt.Sprintf("GameServer or Gate process Exit.")
		if *ddToken == "" {
			token, err := core.CatFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := core.DingDing(alarmInfo, token)
			if er == nil {
				log.Print(*processName+" process run num detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := core.DingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print(*processName+" process run num detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
	}else {
		log.Println("GameServer or Gate process run ok.")
	}
}
