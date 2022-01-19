import { request } from './request'

// 请求组织树接口

// 添加节点

// 编辑节点

// 删除节点

// 拖拽节点

// 点击部门刷新列表接口
export function getClusters(data) {
    return request({
      url: '/hosts',
      method: 'get',
      data
    })
  }

// 添加ip接口
export function insertIp(data) {
  return request({
    url: '/hosts',
    method: 'post',
    data
  })
}
// 编辑ip接口
export function updateIp({ip, ...data}) {
  return request({
    url: `/hosts/${ip}`,
    method: 'post',
    data
  })
}
// 删除ip接口
export function deleteIp(data) {
  return request({
    url: '/hosts',
    method: 'delete',
    data
  })
}