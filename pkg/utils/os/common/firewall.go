/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 19:05:07 2023 +0800
 */
package common

type FireWalldConfig struct {
	Status      string
	Nic         string
	DefaultZone string
	Zones       []string
	Services    []string
}

type FirewalldCMDList struct {
	Service []string    `json:"services"` // 列出允许通过这个防火墙的服务
	Ports   interface{} `json:"ports"`    //列出允许通过这个防火墙的目标端口。（即 需要对外开放的端口）
	Sources []string    `json:"sources"`  // 允许通过的IP或mac
}

func (firewalldConfig *FireWalldConfig) Set() {
	firewalldConfig.Status = ""
	firewalldConfig.Nic = ""
	firewalldConfig.DefaultZone = ""
	firewalldConfig.Zones = nil
	firewalldConfig.Services = nil
}
