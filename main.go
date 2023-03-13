package main

import(
	"os"
	"os/exec"
	"log"
	"strconv"
	"encoding/json"
	"net/http"
	"bytes"
	"strings"
	"flag"
	"fmt"
)

const(
	alarm = "*** 告警了 ***\n\n"
	x = "%"
)

func dingDing(content,token string) error {
	dingDingToken := "https://oapi.dingtalk.com/robot/send?access_token="+token
    dataInfo := map[string]interface{} {
        "msgtype": "markdown",
        "markdown": map[string]string{
            "title": "@",
            "text": content,
        },
        "at": map[string]string{
            "isAtAll": "False",
        },
    }

	jsonData, err := json.Marshal(dataInfo)
	if err != nil {
		log.Fatal("json file change fail: ", err)
		os.Exit(1)
	}
	if _, err := http.Post(dingDingToken, "application/json", bytes.NewBuffer(jsonData));err != nil {
		log.Fatalln("to dingding fail error: ", err)
	}
	return nil		
}

func runCommand(command string) string {
    cmd := exec.Command("/bin/bash", "-c", command)
    output, err := cmd.Output()
    if err != nil {
        log.Fatal("command error: ", err)
    }
    return string(output)
}

func makes(num int)float64{
	rel := float64(num) / float64(1024) / float64(1024)
	return rel
}

type (
	memInfo struct {
		MemTotal float64
		MemUsed float64
		MemFree float64
	}

	diskInfo struct {
		DiskTotal float64
		DiskUsed int
		DiskFree float64
	}

	system struct {
		DiskDir string
	}
)

type systemResources interface {
	memSy() *memInfo
    diskSy(dataDir string) *diskInfo
    cpuSy() int
}


func (s *system) memSy() *memInfo {
	memtotal := runCommand(`free | awk '{print $2}'|awk 'NR==2 {print}'`)
	memtotals,err := strconv.Atoi(strings.Replace(memtotal, "\n", "", -1))
	if err != nil{
		log.Fatal("memTotal str change error: ", err)
	}
	totalRel := makes(memtotals)

	memused := runCommand(`free | awk '{print $3}'|awk 'NR==2 {print}'`)
	memuseds,err := strconv.Atoi(strings.Replace(memused, "\n", "", -1))
	if err != nil{
		log.Fatal("memUsed str change error: ", err)
	}
	usedRel := makes(memuseds)

	memfree := runCommand(`free | awk '{print $4}'|awk 'NR==2 {print}'`)
	memfrees,err := strconv.Atoi(strings.Replace(memfree, "\n", "", -1))
	if err != nil{
		log.Fatal("memFree str change error: ", err)
	}
	freeRel := makes(memfrees)

	mem := &memInfo{
		MemTotal: totalRel,
		MemUsed: usedRel,
		MemFree: freeRel,
	}
	return mem
}

func (s *system)diskSy(dataDir string)*diskInfo {
	cmd1 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $5}'|head -c -2"
	diskused := runCommand(cmd1)
	diskuseds,err := strconv.Atoi(strings.Replace(diskused, "\n", "", -1))
	if err != nil{
		log.Fatal("diskfree str change error: ", err)
	}


	cmd2 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $4}'"
	diskfree := runCommand(cmd2)
	diskfrees,err := strconv.Atoi(strings.Replace(diskfree, "\n", "", -1))
	if err != nil{
		log.Fatal("diskfree str change error: ", err)
	}
	diskFreeRel := makes(diskfrees)

	cmd3 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $2}'"
	disktotal := runCommand(cmd3)
	disktotals, err := strconv.Atoi(strings.Replace(disktotal, "\n", "", -1))
	if err != nil{
		log.Fatal("disktotal str change error: ", err)
	}
	diskTotalRel := makes(disktotals)

	disk := &diskInfo{
		DiskUsed: diskuseds,
		DiskFree: diskFreeRel,
		DiskTotal: diskTotalRel,
	}
	return disk
}

func (s *system) cpuSy()int{
	cpuFree := runCommand(`top -b -n 1|grep -w "Cpu"|awk -F ',' '{print $4}'|cut -f 1 -d "."`)
	cpuFrees, err := strconv.Atoi(strings.Replace(strings.Replace(cpuFree, "\n", "", -1), " ", "", -1))
	if err != nil{
		log.Fatal("cpu str change error: ", err)
	}
	return cpuFrees
}

func main(){
	diskDir := flag.String("diskDataDir", "/", "Specify the storage directory to monitor. The default is / .")
	diskHorizon := flag.Int("d", 80, "Specify how many% utilization of the disk to send an alarm,The default is 80.")
	memHorizon := flag.Float64("m", 2.0, "Specify the number of gigabytes of available memory to send an alarm. The default is 2.0G.")
	cpu := flag.Int("cpu", 20, "Specify how much the CPU is lower than to send an alarm, 20% by default.")
	ddToken := flag.String("token", "", "Specify the Token to send DingDing!")
	flag.Parse()

	s := &system{
		DiskDir: *diskDir,
	}
	if s.memSy().MemFree < *memHorizon{
		alarmInfo := alarm+fmt.Sprintf("Less memory available: %vG", *memHorizon)+"\n\n"+fmt.Sprintf("Available memory: %.2fG", s.memSy().MemFree)
		err := dingDing(alarmInfo, *ddToken)
		if err == nil {
			log.Print("Memory detonate send dingding success!")
		}
	}else {
		log.Printf("memFree: %.2G", s.memSy().MemFree)
	}

	if s.diskSy(s.DiskDir).DiskUsed > *diskHorizon{
		rel := 100 - s.diskSy(s.DiskDir).DiskUsed
		alarmInfo := alarm+fmt.Sprintf("Available disks are less than `%s` : %v%s", s.DiskDir, 20, x)+"\n\n"+fmt.Sprintf("Available disks `%s`: %d%s, %.2fG",s.DiskDir, rel, x, s.diskSy(*diskDir).DiskFree)
		err := dingDing(alarmInfo, *ddToken)
		if err == nil {
			log.Print("Disk detonate send dingding success!")
		}
	}else{
		log.Printf("diskFree: %.2fG\n",s.diskSy(s.DiskDir).DiskFree)
	}

	if s.cpuSy() < *cpu {
		alarmInfo := alarm+fmt.Sprintf("Less CPU available: %v%s", 20, x)+"\n\n"+fmt.Sprintf("Available CPU: %v%s", s.cpuSy(), x)
		err := dingDing(alarmInfo, *ddToken)
		if err == nil {
			log.Print("CPU detonate send dingding success!")
		}
	}else{
		log.Printf("cpuFree: %d%s\n",s.cpuSy(), x)
	}
}