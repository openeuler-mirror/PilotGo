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
 * Date: 2022-05-26 10:25:52
 * LastEditTime: 2022-06-02 10:16:10
 * Description: agent config file service
 ******************************************************************************/
package service

import (
	"fmt"
	"time"
)

// ANSI用GBK解码,解决windows => linux 中文乱码
// func ANSI2Gbk(fileName string) ([]byte, error) {
// 	fi, err := os.Open(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer fi.Close()

// 	decoder := mahonia.NewDecoder("gbk") // 把原来ANSI格式的文本文件里的字符，用gbk进行解码。
// 	fd, err := ioutil.ReadAll(decoder.NewReader(fi))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return fd, nil
// }

// 获取时间的日期函数 => 20200426-17:36:04
func NowTime() string {
	time := time.Now()
	year := time.Year()
	month := time.Month()
	day := time.Day()
	hour := time.Hour()
	minute := time.Minute()
	second := time.Second()
	nowtime := fmt.Sprintf("%d%02d%02d-%02d:%02d:%02d", year, month, day, hour, minute, second)
	return nowtime
}
