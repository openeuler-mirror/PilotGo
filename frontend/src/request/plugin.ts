/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from './request'

// 分页查询插件列表
export function getPluginsPaged(data: any) {
  return request({
    url: '/plugins_paged',
    method: 'get',
    params: data,
  })
}

// 获取插件列表
export function getPlugins() {
  return request({
    url: '/plugins',
    method: 'get',
  })
}

export function addPlugin(data: any) {
  return request({
    url: '/plugins',
    method: 'put',
    data
  })
}

// 启用/停用插件
export function togglePlugin(data: any) {
  return request({
    url: '/plugins/' + data.uuid,
    method: 'post',
    data
  })
}

//删除插件
export function deletePlugins(data: any) {
  return request({
    url: '/plugins/' + data.UUID,
    method: 'delete',
    data
  })
}

// 插件扩展方法
export function pluginExtAPI(data: { url: string, uuids: string[] }) {
  return request({
    url: data.url,
    method: 'post',
    data: data.uuids
  })
}