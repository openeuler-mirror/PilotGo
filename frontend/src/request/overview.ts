/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from './request';

// 获取机器集群概览信息
export function machinesOverview() {
  return request({
    url: '/overview/info',
    method: 'get',
  });
}

// 获取各个部门机器集群概览信息
export function departMachinesOverview() {
  return request({
    url: '/overview/depart_info',
    method: 'get',
  });
}