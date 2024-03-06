import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface BatchMidifiedMessage {
    modified: boolean
}

export const BatchMidifiedMessageStore = defineStore('batch_modified_message', () => {
    const msg = ref<BatchMidifiedMessage>({ modified: false })

    return { msg }
})

export function notify_batch_modified() {
    BatchMidifiedMessageStore().msg.modified = !BatchMidifiedMessageStore().msg.modified
}