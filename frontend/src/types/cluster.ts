/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Mon Mar 18 14:50:14 2024 +0800
 */
export interface MachineInfo {
  cpu: string;
  departid: number;
  departname: string;
  id: number;
  ip: string;
  maintstatus: string;
  runstatus: string;
  systeminfo: string;
  uuid: string;
}

export interface DeptTree {
  id: number;
  pid: number;
  label: string;
  children?: DeptTree[]
}