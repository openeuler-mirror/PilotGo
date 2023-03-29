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
 * LastEditTime: 2022-04-20 14:10:23
 * Description: Execute instruction function
 ******************************************************************************/
package utils

import (
	"bytes"
	"os/exec"
)

func RunCommand(s string) (string, error) {
	// TODO
	cmd := exec.Command("/bin/bash", "-c", "export LANG=en_US.utf8 ; " + s)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}
