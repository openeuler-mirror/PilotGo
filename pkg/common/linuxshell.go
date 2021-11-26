package common

/**
 * @Author: zhang han
 * @Date: 2021/11/18 10:22
 * @Description: linux使用命令的配置
 */

import (
	"bytes"
	"os/exec"
)

func Shells(s string) (string, error) {

	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}
