import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface User {
    // id: number
    name?: string
    email?: string
    departmentID?: number
    department?: string
    roleID?: string
}

export const userStore = defineStore('user', () => {
    const user = ref<User>({})

    function $reset() {
        user.value = {}
    }

    return { user, $reset }
})
