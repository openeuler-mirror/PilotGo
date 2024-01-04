<template>
    <div>
        <el-tree ref="treeRef" :data="authData" show-checkbox node-key="menuName" default-expand-all
            :check-strictly="strictly" :default-expanded-keys="[2, 3]" :default-checked-keys="checkedKeys"
            @check-change="onCheckChanged">
        </el-tree>

        <div class="footer" v-if="showEdit">
            <el-button @click="onResetSelected">重 置</el-button>
            <el-button type="primary" @click="onChangePermission">确 定</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';

import { RespCodeOK } from "@/request/request";
import { changeRolePermission } from "@/request/role";

import { authData } from "./authData";

const treeRef = ref()

const props = defineProps({
    role: {
        type: Object,
        default: {},
    },
    showEdit: {
        type: Boolean,
        default: true,
    }
})

const roleData = ref([])
const checkedKeys = ref<string[]>([])
const strictly = ref(true)

onMounted(() => {
    checkedKeys.value = props.role.menus.split(',').concat(props.role.buttons);
})

function onChangePermission() {
    let menus: string[] = []
    let buttons: string[] = []
    let checkedNodes = treeRef.value.getCheckedNodes();
    checkedNodes.filter((item: any) => item.isMenu).forEach((item: any) => {
        menus.push(item.menuName)
    });
    checkedNodes.filter((item: any) => item.btnId).forEach((item: any) => {
        buttons.push(item.menuName)
    });

    changeRolePermission({
        role: props.role.role,
        buttonId: buttons,
        menus: menus,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            // TODO:

            ElMessage.success("change role permission success:" + resp.msg)
        } else {
            ElMessage.error("failed to change role permission:" + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to change role permission:" + err.msg)
    })
}

function onCheckChanged(data: any, checked: boolean) {
    // 选中button权限时，自动选中页面权限
    if (checked) {
        if (!('isMenu' in data) || !data.isMenu) {
            let checkedMenu: any = null
            for (let i = 0; i < authData.length; i++) {
                let childrens = authData[i].children
                for (let n = 0; n < childrens.length; n++) {
                    if (childrens[n].id == data.id) {
                        checkedMenu = authData[i]
                    }
                }
            }

            if (checkedMenu) {
                treeRef.value.setChecked(checkedMenu.menuName, true)
            }
        }
    }
}

function onResetSelected() {
    treeRef.value.setCheckedNodes([]);
}

</script>

<style lang="scss" scoped>
.footer {
    text-align: right;
}
</style>