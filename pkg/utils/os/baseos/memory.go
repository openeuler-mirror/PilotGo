package baseos

import (
	"fmt"
	"math"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

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
func moduleMatch(name string, value int64, memconf *common.MemoryConfig) {
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

// TODO 完善94-96逻辑，getmemoryconfig接口存在空行bug
func (b *BaseOS) GetMemoryConfig() *common.MemoryConfig {
	output, err := utils.RunCommand("cat /proc/meminfo")
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command: %s\n", err)
	}
	outputlines := strings.Split(output, "\n")
	m := &common.MemoryConfig{}
	for _, line := range outputlines {
		//一次获取一行,_ 获取当前行是否被读完
		if len(line) <= 1 {
			continue
		}
		a, b := reserveRead([]byte(line))
		moduleMatch(string(a), bytesToInt(b), m)
	}
	return m
}
