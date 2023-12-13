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
 * Date: 2022-06-17 13:03:16
 * LastEditTime: 2022-06-17 14:10:23
 * Description: common func
 ******************************************************************************/
package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// string数组转为int数组
func String2Int(strArr string) []int {
	strArrs := strings.Split(strArr, ",")
	res := make([]int, len(strArrs))

	for index, val := range strArrs {
		res[index], _ = strconv.Atoi(val)
	}

	return res
}

func Int2String(intSlice []int) string {
	stringSlice := make([]string, len(intSlice))
	for i, v := range intSlice {
		stringSlice[i] = fmt.Sprint(v)
	}
	return strings.Join(stringSlice, ",")
}
