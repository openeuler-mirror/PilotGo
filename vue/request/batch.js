import { request } from './request'
// 创建批次
export function createBatch(data) {
  return request({
    url: 'batchmanager/createbatch',
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
// 获取批次详情
export function getBatchDetail(data) {
  return request({
    url: '/batchmanager/batchdetail',
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
// 批次详情
export function batchDetail(data) {
  return request({
    url: '/batch/detail',
    method: 'post',
    data
  })
}