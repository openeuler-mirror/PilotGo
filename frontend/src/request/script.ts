/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import request from "./request";

// 获取全部脚本信息
export function getScripts(data: { page: number; size: number }) {
  return request({
    url: "/script/list_all",
    method: "get",
    params: data,
  });
}

// 获取脚本历史信息
export function getScriptHistorys(data: { script_id: number }) {
  return request({
    url: "/script/list_history",
    method: "get",
    params: data,
  });
}

// 添加脚本
export function addScript(data: any) {
  return request({
    url: "/script/create",
    method: "post",
    data,
  });
}

// 更新脚本信息
export function updateScript(data: any) {
  return request({
    url: "/script/update",
    method: "put",
    data,
  });
}

// 删除脚本
export function deleteScript(data: { script_id: number; version?: string }) {
  return request({
    url: "/script/delete",
    method: "delete",
    data,
  });
}

// 运行脚本
export function runScript(data: any) {
  return request({
    url: "/script_auth/run",
    method: "post",
    data,
  });
}

// 获取黑名单列表
export function getScriptBlackList() {
  return request({
    url: "/script/blacklist",
    method: "get",
  });
}

// 更新黑名单列表
export function updateScriptBlackList(data: { white_list: number[] }) {
  return request({
    url: "/script_auth/update_blacklist",
    method: "put",
    data,
  });
}
