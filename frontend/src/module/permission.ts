import { ref } from "vue";
import { ElMessage } from 'element-plus';

import { getPermission } from "@/request/user";
import { RespCodeOK } from "@/request/request";

const userPermissions = ref<any>({})

export function hasPermisson(permission: string): boolean {
    let words = permission.split("/")
    let resource = words[0]
    let operate = words[1]

    if ((resource in userPermissions.value) && userPermissions.value[resource].includes(operate)) {
        return true;
    }
    return false;
}

export function updatePermisson(): void {
    getPermission().then((resp: any) => {
        if (resp.code === RespCodeOK) {
            userPermissions.value = resp.data
        } else {
            ElMessage.error("failed to get machines overview info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machines overview info:" + err.msg)
    })

}