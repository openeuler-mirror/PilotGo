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
// 删除批次
export function delBatch(data) {
  return request({
    url: '/batch/delete',
    method: 'post',
    params: data
  })
}
// 编辑批次
export function updateBatch(data) {
  return request({
    url: '/batch/update',
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