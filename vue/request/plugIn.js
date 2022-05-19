import request from './request'

// 获取插件列表
export function getPlugins() {
  return request({
    url: '/plugin',
    method: 'get',
  })
}
// 添加插件
export function insertPlugin(data) {
  return request({
    url: '/plugin',
    method: 'post',
    data
  })
}
//删除插件
export function deletePlugins(data) {
  return request({
    url: '/plugin',
    method: 'delete',
    data
  })
}