package os

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

type MemoryConfig struct {
	MemTotal     int64
	MemFree      int64
	MemAvailable int64
	Buffers      int64
	Cached       int64
	SwapCached   int64
}

func (conf *MemoryConfig) String() string {
	b, err := json.Marshal(*conf)
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	return out.String()
}

func reverse(res []byte) []byte {
	key := 0
	length := len(res) - 1
	tmp := make([]byte, length+1)

	for {
		tmp[key] = res[length]
		key++
		if length == 0 {
			break
		}
		length--
	}
	return tmp
}
func reserveRead(analyStr []byte) ([]byte, []byte) {
	length := len(analyStr) - 1
	tmp1 := make([]byte, 0)
	tmp2 := make([]byte, 0)

	key := 0
	for {
		if string(analyStr[length-2]) == " " && string(analyStr[length-1]) == "k" && string(analyStr[length]) == "B" {
			length = length - 3
		}
		tmp2 = append(tmp2, analyStr[length])
		if length == 0 {
			break
		}
		if string(analyStr[length-1]) == " " && string(analyStr[length-2]) == " " {
			break
		}
		length--
	}
	key = 0
	for {
		tmp1 = append(tmp1, analyStr[key])
		key++
		if string(analyStr[key]) == ":" {
			break
		}
	}
	return tmp1, reverse(tmp2)
}
func moduleMatch(name string, value int64, memconf *MemoryConfig) {
	if name == "MemTotal" {
		memconf.MemTotal = value
	} else if name == "MemFree" {
		memconf.MemFree = value
	} else if name == "MemAvailable" {
		memconf.MemAvailable = value
	} else if name == "Buffers" {
		memconf.Buffers = value
	} else if name == "Cached" {
		memconf.Cached = value
	} else if name == "SwapCached" {
		memconf.SwapCached = value
	}
}

func bytesToInt(bys []byte) int64 {
	length := float64(len(bys)) - 1
	var x float64
	for _, value := range bys {
		tmp := math.Pow(10, length)
		x = x + (float64(value)-48)*tmp
		length--
	}
	return int64(x)

}

func GetMemoryConfig() MemoryConfig {
	cmd := exec.Command("/bin/sh", "-c", "cat /proc/meminfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
	}
	//使用带缓冲的读取器
	outputBuf := bufio.NewReader(stdout)
	m := MemoryConfig{}
	for {
		//一次获取一行,_ 获取当前行是否被读完
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			// 判断是否到文件的结尾了否则出错
			if err.Error() != "EOF" {
				fmt.Printf("Error :%s\n", err)
			}
			break
		}
		a, b := reserveRead(output)
		moduleMatch(string(a), bytesToInt(b), &m)

	}
	return m
}
