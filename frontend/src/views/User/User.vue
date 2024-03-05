<template>
    <div class="container">
        <PGTable ref="refTable" :data="users" title="用户列表" :showSelect="true" :total="total"
            v-model:selectedData="selectedUsers" :onPageChanged="onPageChanged">
            <template v-slot:action>
                <div class="search">
                    <el-input v-model.trim="searchInput" placeholder="请输入邮箱名进行搜索..." style="width: 300px;" />
                    <el-button type="primary" @click="onSearchUser">搜索</el-button>
                    <el-divider direction="vertical" style="height: 2.5em;" />
                    <el-button type="primary" style="margin-left: 0px;" @click="onAddUser">添加</el-button>
                    <el-popconfirm title="确定删除此用户？" confirm-button-text="确定" cancel-button-text="取消"
                        @confirm="onDeleteUser">
                        <template #reference>
                            <el-button type="primary">删除</el-button>
                        </template>
                    </el-popconfirm>
                    <el-button type="primary">导出</el-button>
                    <el-button type="primary">批量导入</el-button>
                </div>
            </template>

            <template v-slot:content>
                <el-table-column align="center" prop="username" label="用户名">
                </el-table-column>
                <el-table-column align="center" prop="departName" label="部门">
                </el-table-column>
                <el-table-column align="center" prop="role" label="角色">
                </el-table-column>
                <el-table-column align="center" prop="phone" label="手机号">
                </el-table-column>
                <el-table-column align="center" prop="email" label="邮箱">
                </el-table-column>
                <el-table-column align="center" label="操作" fixed="right" class="operate">
                    <template #default="scope">
                        <el-button type="primary" size="small" @click="onUpdateUser(scope.row)">编辑</el-button>
                        <el-button type="danger" size="small" @click="onResetUserPasswd(scope.row)">重置密码</el-button>
                    </template>
                </el-table-column>
            </template>
        </PGTable>

        <el-dialog :title="title" v-model="display" width="560px">
            <AddUser v-if="displayDialog === 'AddUser'" @userUpdated="updateUsers" @close="display = false" />
            <UpdateUser v-if="displayDialog === 'UpdateUser'" :user="editedUser" @userUpdated="updateUsers"
                @close="display = false" />
        </el-dialog>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';

import PGTable from "@/components/PGTable.vue";
import AddUser from "./components/AddUser.vue";
import UpdateUser from "./components/UpdateUser.vue";

import { getUsers, searchUser, resetUserPasswd, deleteUser } from "@/request/user";
import { RespCodeOK } from "@/request/request";

const refTable = ref()

const users = ref([])
const total = ref(0)

onMounted(() => {
    updateUsers()
})

function updateUsers(page: number = 1, size: number = 10) {
    getUsers({
        page: page,
        size: size,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            total.value = resp.total
            users.value = resp.data
        } else {
            ElMessage.error("failed to get users info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get users info:" + err.msg)
    })
}

const display = ref(false)
const displayDialog = ref("")
const title = ref("")

function onAddUser() {
    title.value = "添加用户"
    displayDialog.value = "AddUser"
    display.value = true
}

const selectedUsers = ref<any[]>([])
function onDeleteUser() {
    let params: string[] = []
    selectedUsers.value.forEach((item: any) => {
        params.push(item.email);
    })

    deleteUser({
        delDatas: params
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            updateUsers()
            ElMessage.success("success to delete users:" + resp.msg)
        } else {
            ElMessage.error("failed to delete users:" + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to delete users:" + err.msg)
    })
}

const editedUser = ref<any>({})
function onUpdateUser(user: any) {
    editedUser.value = user
    title.value = "编辑用户"
    displayDialog.value = "UpdateUser"
    display.value = true
}

const searchInput = ref("")

function onSearchUser() {
    // 重置到table第一页
    refTable.value.resetPage()

    if (searchInput.value === "") {
        tableMode = "list"
        updateUsers()
        return
    }

    tableMode = "search"
    searchUserList()
}

function searchUserList(page: number = 1, size: number = 10) {
    searchUser({
        email: searchInput.value
    }, {
        page: page,
        size: size,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            total.value = resp.total
            users.value = resp.data
        } else {
            ElMessage.error("failed to search users:" + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to search users:" + err.msg)
    })
}

function onResetUserPasswd(user: any) {
    resetUserPasswd({
        email: user.email
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            ElMessage.success("reset user password success:" + resp.msg)
        } else {
            ElMessage.error("failed to reset user password:" + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to reset user password" + err.msg)
    })
}

// list:所有用户清单
// search:搜索用户
let tableMode = "list"
function onPageChanged(currentPage: number, currentSize: number) {
    if (tableMode === "search") {
        searchUserList(currentPage, currentSize)
    } else if (tableMode === "list") {
        updateUsers(currentPage, currentSize)
    } else {
        ElMessage.error("invalid table mode:" + tableMode)
    }
}
</script>

<style lang="scss" scoped>
.container {
    height: 100%;
    width: 100%;

    .search {
        height: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
    }

    .el-button {
        margin-left: 5px;
    }

}
</style>