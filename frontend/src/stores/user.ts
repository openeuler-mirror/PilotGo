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
