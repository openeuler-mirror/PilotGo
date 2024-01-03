import request from './request';

// 获取批次
export function getBatches(data: any) {
    return request({
        url: '/batchmanager/batchinfo',
        method: 'get',
        params: data
    })
}

// 获取批次详情
export function getBatchDetail(data: any) {
    return request({
        url: '/batchmanager/batchmachineinfo',
        method: 'get',
        params: data
    })
}

// 删除批次
export function deleteBatch(data: any) {
    return request({
        url: '/batchmanager/deletebatch',
        method: 'post',
        data,
    })
}

// 创建批次
export function createBatch(data: any) {
    return request({
        url: 'macList/createbatch',
        method: 'post',
        data
    })
}

// 编辑批次信息
export function updateBatch(data: any) {
    return request({
        url: 'batchmanager/updatebatch',
        method: 'post',
        data
    })
}
