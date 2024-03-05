<template>
    <div class="container">
        <PGTable :data="roles" title="角色列表" :showSelect="true" :total="total" :onPageChanged="onPageChanged">
            <template v-slot:action>
                <el-button type="primary" @click="onAddRole">添加</el-button>
            </template>
            <template v-slot:content>
                <el-table-column align="center" prop="id" label="角色ID" sortable>
                </el-table-column>
                <el-table-column align="center" prop="role" label="角色名">
                </el-table-column>
                <el-table-column align="center" prop="description" label="描述">
                </el-table-column>
                <el-table-column align="center" label="权限">
                    <template #default="scope">
                        <el-button name="default_all" type="primary" size="small"
                            @click="showRoleDetail(scope.row)">查看</el-button>
                        <el-button name="role_modify" type="primary" size="small" v-if="!(scope.row.role === 'admin')"
                            @click="onEditRolePermission(scope.row)">变更</el-button>
                    </template>
                </el-table-column>
                <el-table-column align="center" label="操作" fixed="right">
                    <template #default="scope">
                        <el-button :disabled="(scope.row.role === 'admin')" name="role_update" size="small" type="primary"
                            @click="onEditRoleInfo(scope.row)">编辑</el-button>
                        <el-popconfirm title="确定删除此角色?" @confirm="onDeleteRole(scope.row)">
                            <template #reference>
                                <el-button :disabled="(scope.row.role === 'admin')" slot="reference" name="role_delete"
                                    type="danger" size="small">
                                    删除
                                </el-button>
                            </template>
                        </el-popconfirm>
                    </template>
                </el-table-column>
            </template>
        </PGTable>

        <el-dialog :title="roleOperateTitle" v-model="showRoleOperate">
            <UpdateRole v-if="operate === 'UpdateRole'" :role="selectedRole" @rolesUpdated="updateRoles" @close="showRoleOperate = false" />
            <AddRole v-if="operate === 'AddRole'" @rolesUpdated="updateRoles" @close="showRoleOperate = false"/>
        </el-dialog>

        <el-drawer :title="roleDetailTitle" v-model="showPermissionDrawler" direction="rtl"
        :destroy-on-close="true">
            <RoleDetail :showEdit="showPermissionEdit" :role="selectedRole" />
        </el-drawer>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';

import PGTable from "@/components/PGTable.vue";
import RoleDetail from "./components/RoleDetail.vue";
import AddRole from "./components/AddRole.vue";
import UpdateRole from "./components/UpdateRole.vue";

import { getRolesPaged, deleteRole } from "@/request/role";
import { RespCodeOK } from "@/request/request";

const roles = ref([])
const total = ref(0)

onMounted(() => {
    updateRoles()
})

function updateRoles(page: number = 1, size: number = 10) {
    getRolesPaged({
        page: page,
        size: size,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            total.value = resp.total
            roles.value = resp.data
        } else {
            ElMessage.error("failed to get role info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get role info:" + err.msg)
    })
}

function onPageChanged(currentPage: number, currentSize: number) {
    updateRoles(currentPage, currentSize)
}

const roleDetailTitle = ref("权限详情")
const showPermissionDrawler = ref(false)
const showPermissionEdit = ref(false)
const selectedRole = ref({})

function showRoleDetail(role: any) {
    selectedRole.value = role
    roleDetailTitle.value = "权限详情"
    showPermissionEdit.value = false
    showPermissionDrawler.value = true
}

function onEditRolePermission(role: any) {
    selectedRole.value = role
    roleDetailTitle.value = "编辑权限"
    showPermissionEdit.value = true
    showPermissionDrawler.value = true
}

const roleOperateTitle = ref("编辑角色")
const showRoleOperate = ref(false)
const operate = ref("")

function onEditRoleInfo(role:any) {
    selectedRole.value = role
    roleOperateTitle.value = "编辑角色"
    showRoleOperate.value = true
    operate.value = "UpdateRole"
}

function onAddRole() {
    roleOperateTitle.value = "添加角色"
    showRoleOperate.value = true
    operate.value = "AddRole"
}

function onDeleteRole(role: any) {
    deleteRole({
        role: role.id
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            ElMessage.success("删除角色成功:" + resp.msg)
            updateRoles()
        } else {
            ElMessage.error("failed to delete role:" + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to delete role:" + err.msg)
    })
}

</script>

<style lang="scss" scoped>
.container {
    height: 100%;
    width: 100%;
}
</style>