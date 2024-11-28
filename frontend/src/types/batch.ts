/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Fri Mar 15 11:38:55 2024 +0800
 */
// 批次列表
export interface BatchItem {
  CreatedAt:String;
  DeletedAt:null;
  Depart:String;
  DepartName: String;
  ID: number;
  UpdatedAt: String;
  description: String;
  manager: String;
  name: String;
}

// 批次详情列表
export interface BatchMachineInfo {
  CPU: String;
  departid: number;
  id:number;
  ip: String;
  machineuuid: String;
  maintstatus: String;
  runstatus: String;
  sysinfo: String;
}