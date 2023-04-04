package baseos

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

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

// TODO 完善94-96逻辑，getmemoryconfig接口存在空行bug
func (b *BaseOS) GetMemoryConfig() *common.MemoryConfig {
	output, err := utils.RunCommand("cat /proc/meminfo")
	if err != nil {
		fmt.Printf("failed to get memory config: %s\n", err)
	}
	outputlines := strings.Split(output, "\n")
	m := &common.MemoryConfig{}
	reg := regexp.MustCompile(`[a-zA-Z\s]+`)
	for _, line := range outputlines {
		//一次获取一行,_ 获取当前行是否被读完
		if line == "" {
			continue
		}
		k := strings.Split(line, ":")[0]
		v := reg.ReplaceAllString(strings.Split(line, ":")[1], "")
		vint64, _ := strconv.ParseInt(v, 10, 64)
		moduleMatch(k, vint64, m)
	}
	return m
}
