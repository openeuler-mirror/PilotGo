<template>
  <div class="container">
    <PGTable :data="roles" title="角色列表" :showSelect="true" :total="total" :onPageChanged="onPageChanged">
      <template v-slot:action>
        <auth-button auth="button/role_add" @click="onAddRole">添加</auth-button>
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
            <el-button size="small" @click="handleAuth(scope.row, '权限详情')">查看</el-button>
            <auth-button auth="button/role_modify" size="small" @click="handleAuth(scope.row, '编辑权限')">变更</auth-button>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" fixed="right">
          <template #default="scope">
            <auth-button auth="button/role_update" size="small" @click="onEditRoleInfo(scope.row)">编辑</auth-button>
            <el-popconfirm title="确定删除此角色?" @confirm="onDeleteRole(scope.row)">
              <template #reference>
                <auth-button slot="reference" auth="button/role_delete" plain type="danger" size="small">
                  删除
                </auth-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </template>
    </PGTable>

    <el-dialog :title="roleOperateTitle" v-model="showRoleOperate" :destroy-on-close="true">
      <UpdateRole v-if="operate === 'UpdateRole'" :role="selectedRole" @rolesUpdated="updateRoles"
        @close="showRoleOperate = false" />
      <AddRole v-if="operate === 'AddRole'" @rolesUpdated="updateRoles" @close="showRoleOperate = false" />
      <RoleForm :showEdit="showPermissionEdit" v-if="operate === 'updateAuth'" @rolesUpdated="updateRoles"
        :role="selectedRole" @close="showRoleOperate = false" />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';

import PGTable from "@/components/PGTable.vue";
import AddRole from "./components/AddRole.vue";
import UpdateRole from "./components/UpdateRole.vue";

import { getRolesPaged, deleteRole } from "@/request/role";
import { RespCodeOK } from "@/request/request";
import RoleForm from "./components/RoleForm.vue";
import AuthButton from "@/components/AuthButton.vue";

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

const roleOperateTitle = ref("权限详情")
const showRoleOperate = ref(false)
const showPermissionEdit = ref(false)
const selectedRole = ref({})
const operate = ref("")
const handleAuth = (role: any, type: string) => {
  selectedRole.value = role
  roleOperateTitle.value = type
  showPermissionEdit.value = type === '编辑权限' ? true : false;
  showRoleOperate.value = true
  operate.value = "updateAuth"
}

function onEditRoleInfo(role: any) {
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