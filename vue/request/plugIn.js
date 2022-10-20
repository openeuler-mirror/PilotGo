import request from './request'

// 获取插件列表
export function getPlugins() {
  return request({
    url: '/plugin/list',
    method: 'get',
  })
}
//安装插件
export function loadPlugin(data) {
  return request({
    url: '/plugin/load',
    method: 'post',
    data
  })
}
//卸载插件
export function unLoadPlugin(data) {
  return request({
    url: '/plugin/unload',
    method: 'post',
    data
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