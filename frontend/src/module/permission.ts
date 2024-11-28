/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import { ref } from "vue";
import { ElMessage } from 'element-plus';

import { getPermission } from "@/request/user";
import { RespCodeOK } from "@/request/request";
import { updateSidebarItems } from "@/router";

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
          userPermissions.value = resp.data;
            updateSidebarItems();
        } else {
            ElMessage.error("failed to get machines overview info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machines overview info:" + err.msg)
    })

}