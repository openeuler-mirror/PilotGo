import { request } from './request'

// 请求组织树接口
export function getDeparts(data) {
  return request({
    url: 'machinemanager/departinfo',
    method: 'get',
    params: data
  })
}
// 点击获取当前子节点接口
export function getChildNode(data) {
  return request({
    url: 'machinemanager/depart',
    method: 'get',
    params: data
  })
}
// 添加节点
export function addDepart(data) {
  return request({
    url: 'machinemanager/adddepart',
    method: 'post',
    params: data
  })
}
// 编辑节点
export function updateDepart(data) {
  return request({
    url: 'machinemanager/updatedepart',
    method: 'get',
    params: data
  })
}
// 删除节点
export function deleteDepart(data) {
  return request({
    url: 'machinemanager/t',
    method: 'get',
    params: data
  })
}
// 拖拽节点

// 点击部门刷新列表接口
export function getClusters(data) {
    return request({
      url: 'machinemanager/machineinfo',
      method: 'get',
      params: data
    })
  }

// 添加ip接口
export function insertIp(data) {
  return request({
    url: '/machinemanager/addmachine',
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
    url: '/machinemanager/deletemachinedata',
    method: 'post',
    params: data
  })
}

// 根据ip获取机器信息
export function getDeviceInfo(data) {
  return request({
    url: 'machinemanager/deviceinfo',
    method: 'get',
    params: data
  })
}