/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface Taginfo {
  path: string
  title: string
  fullpath: string
  query: object
  meta: any | null
}

// 存储的router信息用于sidebar动态生成
export const tagviewStore = defineStore('tagview', () => {
    const taginfos = ref<Taginfo[]>([])
    function $reset() {
      taginfos.value = [];
    }
    return { taginfos, $reset }
})
