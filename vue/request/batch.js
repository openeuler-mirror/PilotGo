import { request } from './request'
// 创建批次
export function createBatch(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}
// 获取批次
export function getBatch() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}
// 删除批次
export function deleteBatch(data) {
  return request({
    url: '/user/searchAll',
    method: 'post',
    params: data
  })
}
// 编辑批次
export function updateBatch(data) {
  return request({
    url: '/user/register',
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