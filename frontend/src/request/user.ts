import request from './request';

// 用户登录
export function loginByEmail(data: any) {
    return request({
        url: '/user/login',
        method: 'post',
        data,
    });
}

// 用户登出
export function logout() {
    return request({
        url: '/user/logout',
        method: 'get',
    })
}

// 获取全部用户信息
export function getUsers(data: any) {
    return request({
        url: '/user/searchAll',
        method: 'get',
        params: data,
    });
}

// 按邮箱查找用户
export function searchUser(data: any, params: any) {
    return request({
        url: '/user/userSearch',
        method: 'post',
        params: params,
        data,
    });
}

// 添加用户
export function addUser(data: any) {
    return request({
        url: '/user/register',
        method: 'post',
        data
    })
}

// 更新用户信息
export function updateUser(data: any) {
    return request({
        url: '/user/update',
        method: 'post',
        data
    })
}

// 重置用户密码
export function resetUserPasswd(data: any) {
    return request({
        url: '/user/reset',
        method: 'post',
        data
    })
}

// 删除用户
export function deleteUser(data: any) {
    return request({
        url: '/user/delete',
        method: 'post',
        data
    })
}

// 获取当前登录用户信息
export function getCurrentUser() {
    return request({
        url: '/user/info',
        method: 'get'
    })
}

// 获取用户权限
export function getPermission() {
    return request({
        url: '/user/permission',
        method: 'post',
    })
}