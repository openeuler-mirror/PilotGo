package common

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// 通过 /proc/cpuinfo来获取CPU型号
type CPUInfo struct {
	ModelName string
	CpuNum    int
}

func (cpu *CPUInfo) String() string {
	b, err := json.Marshal(*cpu)
	if err != nil {
		return fmt.Sprintf("%+v", *cpu)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *cpu)
	}
	return out.String()
}
