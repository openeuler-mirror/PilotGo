/*
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: zhaozhenfang
 * @Date: 2022-05-12 15:28:36
 * @LastEditTime: 2022-06-06 16:35:23
 * @Description: provide agent log manager of pilotgo
 */
import request from './request'

// 获取所有repo源
export function getRepos(data) {
  return request({
    url: '/config/repos',
    method: 'get',
    params: data
  })
}

// 获取repo详情
export function getRepoDetail(data) {
  return request({
    url: '/config/read_file',
    method: 'get',
    params: data
  })
}

// 编辑repo源
export function updateRepo(data) {
  return request({
    url: '/config/file_edit',
    method: 'post',
    data
  })
}

// 下载repo到库
export function saveRepo(data) {
  return request({
    url: '/config/fileSaveAdd',
    method: 'post',
    data
  })
}

// 库文件列表
export function libFileList(data) {
  return request({
    url: '/config/file_all',
    method: 'get',
    params: data
  })
}

// 库文件模糊查询
export function libFileSearch(data) {
  return request({
    url: '/config/file_search',
    method: 'post',
    data
  })
}

// 更新库文件
export function updateLibFile(data) {
  return request({
    url: '/config/file_update',
    method: 'post',
    data
  })
}

// 删除库文件
export function delLibFile(data) {
  return request({
    url: '/config/file_delete',
    method: 'post',
    data
  })
}

// 历史文件列表
export function lastFileList(data) {
  return request({
    url: '/config/lastfile_all',
    method: 'get',
    params: data
  })
}

// 历史文件模糊查询
export function lastFileSearch(data) {
  return request({
    url: '/config/lastfile_search',
    method: 'post',
    data
  })
}