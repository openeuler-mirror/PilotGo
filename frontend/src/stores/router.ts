import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface Menu {
  path: string
  title: string
  hidden: boolean
  panel: string
  icon: string
  subMenus: Menu[] | null
  type?:string // 标识插件，权限到位后可删除
}

// 存储的router信息用于sidebar和router动态生成
export const routerStore = defineStore('router', {
  state: () => {
    return {
      menus: [] as Menu[],
      routers: [] as any
    }
  },
  actions:  {
    reset(){
      this.menus = [];
      this.routers = [];
    }
  },
  persist: true
})

