import { request } from './request'
// 用户登录
export function loginByEmail(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}
// 用户退出
export function logout() {
  return request({
    url: '/user/logout',
    method: 'post',
  })
}
// 获取当前登录用户信息
export function getCurUser() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}
// 获取全部用户信息
export function getUsers(data) {
  return request({
    url: '/user/searchAll',
    method: 'get',
    params: data
  })
}
// 添加用户
export function addUser(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}
// 编辑用户
export function updateUser(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}
// 删除用户
export function delUser(data) {
  return request({
    url: '/user/delete',
    method: 'post',
    data
  })
}
// 重置密码
export function resetPwd(data) {
  return request({
    url: '/user/reset',
    method: 'post',
    data
  })
}
// 批量导入用户
export function importUser(data) {
  return request({
    url: '/user/import',
    method: 'post',
    data
  })
}