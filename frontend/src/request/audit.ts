import request from './request'

// 获取日志列表
export function getLogs(data:any) {
  return request({
    url: '/log/log_all',
    method: 'get',
    params: data
  })
}


// 获取子日志
export function getLogChildrens(data:{uuid:string}) {
  return request({
    url: '/log/log_child',
    method: 'get',
    params: data
  })
}