/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Apr 21 05:57:14 2022 +0800
 */
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

// 数组去重
func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
