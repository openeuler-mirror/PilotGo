import { request_v1 } from './request_v1'

// 获取插件列表
export function getPlugins() {
  return request_v1({
    url: '/plugins',
    method: 'get',
  })
}

// 添加插件
export function insertPlugin(data) {
  return request_v1({
    url: '/plugins',
    method: 'put',
    data
  })
}

// 启用/停用插件
export function unLoadPlugin(data) {
  return request_v1({
    url: '/plugins/'+data.UUID,
    method: 'post',
    data
  })
}

//删除插件
export function deletePlugins(data) {
  return request_v1({
    url: '/plugins/'+data.UUID,
    method: 'delete',
    data
  })
}