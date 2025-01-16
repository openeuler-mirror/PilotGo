/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from "./request";

// 获取批次分页
export function getBatches(data: any) {
  return request({
    url: "/batchmanager/batchinfo",
    method: "get",
    params: data,
  });
}

// 获取所有批次不分页
export function getAllBatches() {
  return request({
    url: "/batchmanager/batchinfo_nopage",
    method: "get",
  });
}

// 获取批次详情
export function getBatchDetail(data: any) {
  return request({
    url: "/batchmanager/batchmachineinfo",
    method: "get",
    params: data,
  });
}

// 删除批次
export function deleteBatch(data: any) {
  return request({
    url: "/batchmanager/deletebatch",
    method: "post",
    data,
  });
}

// 创建批次
export function createBatch(data: any) {
  return request({
    url: "batchmanager/createbatch",
    method: "post",
    data,
  });
}

// 编辑批次信息
export function updateBatch(data: any) {
  return request({
    url: "batchmanager/updatebatch",
    method: "post",
    data,
  });
}
