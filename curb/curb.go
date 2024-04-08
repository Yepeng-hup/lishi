package curb

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	fileDbtxtPath     = "/tmp/db.txt"
	destfileDbtxtPath = "/tmp/new.txt"
)

func WriteTxtDbData(monName string, monValue float64) bool {
	file, err := os.OpenFile(fileDbtxtPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file error: ", err.Error())
		return false
	}
	defer file.Close()
	dbData := fmt.Sprintf("%s: %f\n", monName, monValue)
	if _, err := file.WriteString(dbData); err != nil {
		log.Println("write file error: ", err.Error())
		return false
	}
	return true
}

func ReadTxtDbData(monName string) (float64, error, bool) {
	re := regexp.MustCompile(`:\s+(\d+\.\d+)`)
	file, err := os.Open(fileDbtxtPath)
	if err != nil {
		return 0, fmt.Errorf("open file error: %s", err.Error()), false
	}
	defer file.Close()

	// 创建一个读取器
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if strings.Contains(line, monName) {
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				memValue := match[1]
				mem, err := strconv.ParseFloat(memValue, 64)
				if err != nil {
					return 0, fmt.Errorf("str convert float64 fail,%s", err.Error()), false
				}
				return mem, nil, true

			} else {
				return 0, fmt.Errorf("error: No matching values found."), false
			}
		}
		if err != nil {
			break
		}
	}
	return 0, fmt.Errorf("error: not this [%s] value", monName), false
}

func DeleteTxtDbData(monName string) (bool, error) {
	inputFileName := fileDbtxtPath
	// 临时文件
	outputFileName := destfileDbtxtPath

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return false, fmt.Errorf("Unable to open the original file,%s", err.Error())
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return false, fmt.Errorf("Unable to create new file,%s", err.Error())
	}

	// 创建扫描器和写入器
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		// 检查行是否包含指定的字符串
		if !strings.Contains(line, monName) {
			// 如果不包含指定的字符串，则将该行写入新文件
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return false, fmt.Errorf("Error writing new file,%s", err.Error())
			}
		}
	}

	// 检查扫描和写入是否发生错误
	if err := scanner.Err(); err != nil {

		return false, fmt.Errorf("Error scanning original file,%s", err.Error())
	}
	if err := writer.Flush(); err != nil {
		return false, fmt.Errorf("Error writing new file,%s", err.Error())
	}

	data, err := ioutil.ReadFile(outputFileName)
	if err != nil {
		return false, fmt.Errorf("read file err,%v\n", err)
	}
	err = ioutil.WriteFile(inputFileName, data, 0666)
	if err != nil {
		return false, fmt.Errorf("write file error,%v\n", err.Error())
	}

	inputFile.Close()
	outputFile.Close()

	if err := os.Remove(outputFileName); err != nil {
		return false, fmt.Errorf("error deleting original file,%s", err.Error())
	}
	return true, nil
}
