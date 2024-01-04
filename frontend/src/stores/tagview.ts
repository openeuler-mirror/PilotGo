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

    return { taginfos }
})
