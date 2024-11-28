/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from './request'

// 获取日志列表
export function getLogs(data:any) {
  return request({
    url: '/log/log_all',
    method: 'get',
    params: data
  })
}


// 获取子日志
export function getLogChildrens(data:{uuid:string}) {
  return request({
    url: '/log/log_child',
    method: 'get',
    params: data
  })
}