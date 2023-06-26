/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-11-18 13:03:16
 * LastEditTime: 2023-06-26 19:09:58
 * Description: Execute instruction function
 ******************************************************************************/
package utils

import (
	"fmt"
	"io/ioutil"
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

	b1, err := ioutil.ReadAll(stdout)
	if err != nil {
		return 0, "", "", err
	}
	s1 := strings.TrimRight(string(b1), "\n")

	b2, err := ioutil.ReadAll(stderr)
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
func RunScript(absPath string) (*CmdResult, error) {
	cmd := exec.Command("/bin/bash", absPath)
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

	b1, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	s1 := strings.TrimRight(string(b1), "\n")

	b2, err := ioutil.ReadAll(stderr)
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
