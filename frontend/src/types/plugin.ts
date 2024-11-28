/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Feb 20 11:33:42 2024 +0800
 */

export interface Extention {
  name: string;
  permission: string;
  type: string;
  url: string;
  parentName?: string;
}

export type ExtArr = Array<Extention>;