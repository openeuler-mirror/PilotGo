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
 * Date: 2021-05-20 09:08:08
 * LastEditTime: 2022-05-20 16:25:41
 * Description: 机器操作日志业务逻辑
 ******************************************************************************/
package service

import (
	"fmt"
	"net/http"
	"strconv"
)

// 计算批量机器操作的状态：成功数，总数目，比率
func BatchActionStatus(StatusCodes []string) (status string) {
	var StatusOKCounts int
	for _, success := range StatusCodes {
		if success == strconv.Itoa(http.StatusOK) {
			StatusOKCounts++
		}
	}
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(StatusOKCounts)/float64(len(StatusCodes))), 64)
	rate := strconv.FormatFloat(num, 'f', 2, 64)
	status = strconv.Itoa(StatusOKCounts) + "," + strconv.Itoa(len(StatusCodes)) + "," + rate
	return
}
