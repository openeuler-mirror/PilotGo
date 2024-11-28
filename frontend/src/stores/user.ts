/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface User {
    // id: number
    name?: string
    email?: string
    departmentID?: number
    department?: string
    roleID?: string
    role?:string[]
}

export const userStore = defineStore('user', {
  state: () => {
    return {
      user: {} as User,
    }
  },
  actions: {
    $reset() {
      this.user = {}
    }
  },
  persist:true
})
