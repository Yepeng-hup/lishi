package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)



func DingDing(content,token string) error {
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

func RunCommand(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "nil", fmt.Errorf("command error: ", err)
	}
	return string(output), nil
}

func Makes(num int)float64{
	rel := float64(num) / float64(1024) / float64(1024)
	return rel
}

func CatFile(filePath string)(string, error){
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

