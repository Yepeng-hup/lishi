package main

import(
	"io"
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
	f1 = "*** "
	f2 = " ***\n\n"
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
		return fmt.Errorf("json file change fail: ", err)
	}
	if _, err := http.Post(dingDingToken, "application/json", bytes.NewBuffer(jsonData));err != nil {
		return fmt.Errorf("to dingding fail error: ", err)
	}
	return nil		
}

func runCommand(command string) (string, error) {
    cmd := exec.Command("/bin/bash", "-c", command)
    output, err := cmd.Output()
    if err != nil {
        return "nil", fmt.Errorf("command error: ", err)
    }
    return string(output), nil
}

func makes(num int)float64{
	rel := float64(num) / float64(1024) / float64(1024)
	return rel
}

func catFile(filePath string)(string, error){
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	var data []byte
	for {
		// 将文件中读取的byte存储到buf中
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return "nil", fmt.Errorf(err.Error())
		}
		if n == 0 {
			break
		}
		// 将读取到的结果追加到data切片中
		data = append(data, buf[:n]...)
	}
	r := strings.NewReplacer("\n", "")
	return r.Replace(string(data)), nil
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

	diskd struct {
		DiskDir string
	}

	system struct {}

	alarm struct {
		t string
	}
)

func (s *system) memSy() *memInfo {
	memtotal, err := runCommand(`free | awk '{print $2}'|awk 'NR==2 {print}'`)
	if err != nil {
		log.Fatal(err.Error())
	}
	memtotals,err := strconv.Atoi(strings.Replace(memtotal, "\n", "", -1))
	if err != nil{
		log.Fatal("memTotal str change error: ", err)
	}
	totalRel := makes(memtotals)

	memused, err := runCommand(`free | awk '{print $3}'|awk 'NR==2 {print}'`)
	if err != nil {
		log.Fatal(err.Error())
	}
	memuseds,err := strconv.Atoi(strings.Replace(memused, "\n", "", -1))
	if err != nil{
		log.Fatal("memUsed str change error: ", err)
	}
	usedRel := makes(memuseds)

	memfree, err := runCommand(`free | awk '{print $4}'|awk 'NR==2 {print}'`)
	if err != nil {
		log.Fatal(err.Error())
	}
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
	diskused, err := runCommand(cmd1)
	if err != nil {
		log.Fatal(err.Error())
	}
	diskuseds,err := strconv.Atoi(strings.Replace(diskused, "\n", "", -1))
	if err != nil{
		log.Fatal("diskfree str change error: ", err)
	}


	cmd2 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $4}'"
	diskfree, err := runCommand(cmd2)
	if err != nil {
		log.Fatal(err.Error())
	}
	diskfrees,err := strconv.Atoi(strings.Replace(diskfree, "\n", "", -1))
	if err != nil{
		log.Fatal("diskfree str change error: ", err)
	}
	diskFreeRel := makes(diskfrees)

	cmd3 := "df | egrep -w "+"\""+ dataDir+"\"" +"|awk '{print $2}'"
	disktotal, err := runCommand(cmd3)
	if err != nil {
		log.Fatal(err.Error())
	}
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
	cpuFree, err := runCommand(`top -b -n 1|grep -w "Cpu"|awk -F ',' '{print $4}'|cut -f 1 -d "."`)
	if err != nil {
		log.Fatal(err.Error())
	}
	cpuFrees, err := strconv.Atoi(strings.Replace(strings.Replace(cpuFree, "\n", "", -1), " ", "", -1))
	if err != nil{
		log.Fatal("cpu str change error: ", err)
	}
	return cpuFrees
}


func main(){
	if len(os.Args) < 2 {
		log.Fatal("Please enter parameters or -h to view the usage of parameters.")
	}

	diskDir := flag.String("diskDataDir", "/", "Specify the storage directory to monitor. The default is / .")
	diskHorizon := flag.Int("d", 80, "Specify how many utilization of the disk to send an alarm,The default is 80%.")
	memHorizon := flag.Float64("m", 2.0, "Specify the number of gigabytes of available memory to send an alarm. The default is 2.0G.")
	cpu := flag.Int("cpu", 20, "Specify how much the CPU is lower than to send an alarm, 20% by default.")
	ddToken := flag.String("token", "", "Specify the Token to send DingDing.")
	ddTokenFile := flag.String("token_filePath", "", "File path with token written.")
	title := flag.String("t", "", "push DindDing Keyword Title.")
	flag.Parse()

	d := &diskd{
		DiskDir: *diskDir,
	}
	a := &alarm{
		t: *title,
	}
	s := system{}

	if s.memSy().MemFree < *memHorizon{
		alarmInfo := f1+a.t+f2+fmt.Sprintf("Less memory available: %vG", *memHorizon)+"\n\n"+fmt.Sprintf("Available memory: %.2fG", s.memSy().MemFree)
		if *ddToken == "" {
			token, err := catFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := dingDing(alarmInfo, token)
			if er == nil {
				log.Print("mem detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := dingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print("mem detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
	}else {
		log.Printf("memFree: %.2G", s.memSy().MemFree)
	}

	if s.diskSy(d.DiskDir).DiskUsed > *diskHorizon{
		rel := 100 - s.diskSy(d.DiskDir).DiskUsed
		alarmInfo := f1+a.t+f2+fmt.Sprintf("Available disks are less than `%s` : %v%s", d.DiskDir, 20, x)+"\n\n"+fmt.Sprintf("Available disks `%s`: %d%s, %.2fG",d.DiskDir, rel, x, s.diskSy(*diskDir).DiskFree)
		if *ddToken == "" {
			token, err := catFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := dingDing(alarmInfo, token)
			if er == nil {
				log.Print("Disk detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := dingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print("Disk detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
	}else{
		log.Printf("diskFree: %.2fG\n",s.diskSy(d.DiskDir).DiskFree)
	}

	if s.cpuSy() < *cpu {
		alarmInfo := f1+a.t+f2+fmt.Sprintf("Less CPU available: %v%s", 20, x)+"\n\n"+fmt.Sprintf("Available CPU: %v%s", s.cpuSy(), x)
		if *ddToken == "" {
			token, err := catFile(*ddTokenFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			er := dingDing(alarmInfo, token)
			if er == nil {
				log.Print("cpu detonate send dingding success!")
			}else {
				log.Print(er.Error())
			}
		}else {
			err := dingDing(alarmInfo, *ddToken)
			if err == nil {
				log.Print("cpu detonate send dingding success!")
			}else {
				log.Print(err.Error())
			}
		}
	}else{
		log.Printf("cpuFree: %d%s\n",s.cpuSy(), x)
	}
}
