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
 * @Date: 2022-02-25 16:33:46
 * @LastEditTime: 2022-06-17 16:25:47
 * @Description: provide agent log manager of pilotgo
 */
import request from './request'
// 创建批次
export function createBatch(data) {
  return request({
    url: 'macList/createbatch',
    method: 'post',
    data
  })
}
// 获取批次
export function getBatches(data) {
  return request({
    url: '/batchmanager/batchinfo',
    method: 'get',
    params: data
  })
}
// 获取批次-不分页
export function getAllBatches() {
  return request({
    url: '/batchmanager/selectbatch',
    method: 'get',
  })
}
// 获取批次详情
export function getBatchDetail(data) {
  return request({
    url: '/batchmanager/batchmachineinfo',
    method: 'get',
    params: data
  })
}
// 删除批次
export function delBatches(data) {
  return request({
    url: 'batchmanager/deletebatch',
    method: 'post',
    data
  })
}
// 编辑批次
export function updateBatch(data) {
  return request({
    url: 'batchmanager/updatebatch',
    method: 'post',
    data
  })
}