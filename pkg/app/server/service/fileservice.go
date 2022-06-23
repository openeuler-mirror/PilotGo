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
	"strings"
	"time"
)

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

// dhcp方式配置网络
func NetworkDHCP(net []interface{}) (text string) {
	for _, n := range net {
		nn := n.(map[string]interface{})
		for key, value := range nn {
			if key == "BOOTPROTO" {
				text += key + "=" + "dhcp" + "\n"
			} else if key == "IPADDR" {
				break
			} else if key == "NETMASK" {
				break
			} else if key == "GATEWAY" {
				break
			} else if key == "DNS1" {
				break
			} else if key == "DNS2" {
				break
			} else {
				text += key + "=" + value.(string) + "\n"
			}
		}
	}
	return
}

// static方式配置网络
func NetworkStatic(net []interface{}, ip string, netmask string, gateway string, dns1 string, dns2 string) (text string) {
	for _, n := range net {
		nn := n.(map[string]interface{})
		for key, value := range nn {
			if key == "BOOTPROTO" {
				text += key + "=" + "static" + "\n"
			} else if key == "IPADDR" {
				text += key + "=" + ip + "\n"
			} else if key == "NETMASK" {
				text += key + "=" + netmask + "\n"
			} else if key == "GATEWAY" {
				text += key + "=" + gateway + "\n"
			} else if key == "DNS1" {
				text += key + "=" + dns1 + "\n"
			} else if key == "DNS2" && len(dns2) != 0 {
				text += key + "=" + dns2 + "\n"
			} else {
				text += key + "=" + value.(string) + "\n"
			}
		}
	}
	if ok := strings.Contains(text, "IPADDR"); !ok {
		t := "IPADDR" + "=" + ip + "\n"
		text += t
	}
	if ok := strings.Contains(text, "NETMASK"); !ok {
		t := "NETMASK" + "=" + netmask + "\n"
		text += t
	}
	if ok := strings.Contains(text, "GATEWAY"); !ok {
		t := "GATEWAY" + "=" + gateway + "\n"
		text += t
	}
	if ok := strings.Contains(text, "DNS1"); !ok {
		t := "DNS1" + "=" + dns1 + "\n"
		text += t
	}
	if ok := strings.Contains(text, "DNS2"); !ok {
		if len(dns2) != 0 {
			t := "DNS2" + "=" + dns2 + "\n"
			text += t
		}
	}
	return
}
