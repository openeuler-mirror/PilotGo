/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Nov 26 09:45:08 2021 +0800
 */
package utils

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type CmdResult struct {
	RetCode int
	Stdout  string
	Stderr  string
}

func RunCommand(s string) (int, string, string, error) {
	cmd := exec.Command("/bin/bash", "-c", "export LANG=en_US.utf8 ; "+s)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, "", "", err
	}

	exitCode := 0
	err = cmd.Start()
	if err != nil {
		return 0, "", "", err
	}

	b1, err := io.ReadAll(stdout)
	if err != nil {
		return 0, "", "", err
	}
	s1 := strings.TrimRight(string(b1), "\n")

	b2, err := io.ReadAll(stderr)
	if err != nil {
		return 0, "", "", err
	}
	s2 := strings.TrimRight(string(b2), "\n")

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		e, ok := err.(*exec.ExitError)
		if !ok {
			return 0, "", "", err
		}
		exitCode = e.ExitCode()
	}

	return exitCode, s1, s2, nil
}

// 运行指定的shell脚本文件
func RunScript(absPath string, params []string) (*CmdResult, error) {
	arr_command := []string{}
	arr_command = append(arr_command, absPath)
	arr_command = append(arr_command, params...)

	cmd := exec.Command("/bin/bash", arr_command...)
	cmd.Env = append(cmd.Env, "LANG=en_US.utf8")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	exitCode := 0
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	b1, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	s1 := strings.TrimRight(string(b1), "\n")

	b2, err := io.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	s2 := strings.TrimRight(string(b2), "\n")

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		e, ok := err.(*exec.ExitError)
		if !ok {
			return nil, err
		}
		exitCode = e.ExitCode()
	}

	return &CmdResult{
		RetCode: exitCode,
		Stdout:  s1,
		Stderr:  s2,
	}, nil
}

// 检查字符串是否包含特殊字符
func CheckString(str string) bool {
	// 定义允许的字符范围
	for _, char := range str {
		// 检查字符是否在允许的字符范围内
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '.' || char == '+' || char == '_') {
			return false // 如果包含特殊字符，则返回 true
		}
	}
	return true
}
