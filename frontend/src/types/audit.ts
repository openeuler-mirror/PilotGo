/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Mar 19 11:41:25 2024 +0800
 */
export interface AuditItem {
  CreatedAt: string;
  DeletedAt: string;
  ID: number;
  /* 
  * @params Isempty 
  * 0:has children logs  
  * 1:no children logs
  *  */
  Isempty: number; 
  hasChildren?: boolean;
  UpdatedAt: string;
  action: string;
  log_uuid: string;
  message: string;
  module: string;
  parent_uuid: string;
  status: string;
  user_id: number;
}