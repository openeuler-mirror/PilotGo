<template>
    <div class="content">
        <div class="users">
            <div class="current">
                当前用户： {{ currentUser.Username }}
            </div>
            <div class="search">
                <el-autocomplete style="width:50%" class="inline-input" v-model="userName" placeholder="请输入用户名"
                    :fetch-suggestions="querySuggestions" @select="onSelectUser"></el-autocomplete>
                <el-button plain type="primary">搜索</el-button>
            </div>
        </div>
        <div class="info">
            <p class="title">用户信息详情：</p>
            <el-descriptions :column="3" border>
                <el-descriptions-item label="用户名"> {{ userInfo.Username }} </el-descriptions-item>
                <el-descriptions-item label="用户ID"> {{ userInfo.UserId }} </el-descriptions-item>
                <el-descriptions-item label="用户组ID"> {{ userInfo.GroupId }} </el-descriptions-item>
                <el-descriptions-item label="家目录"> {{ userInfo.HomeDir }} </el-descriptions-item>
                <el-descriptions-item label="shell类型"> {{ userInfo.ShellType }} </el-descriptions-item>
                <el-descriptions-item label="描述"> {{ userInfo.Description }} </el-descriptions-item>
            </el-descriptions>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus';

import { getCurrentUser, getMachineAllUser } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

// 机器UUID
const machineID = ref(route.params.uuid)


const userName = ref("")
const allUser = ref<any[]>([])
const currentUser = ref<any>({})
const userInfo = ref<any>({})

onMounted(() => {
    getMachineAllUser({ uuid: machineID.value }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            allUser.value = resp.data.user_all

            // 嵌套调用，避免两者请求不同步
            getCurrentUser({ uuid: machineID.value }).then((resp: any) => {
                if (resp.code === RespCodeOK) {
                    currentUser.value = resp.data.user_info

                    userInfo.value = allUser.value.filter((item: any) => item.Username === currentUser.value.Username)[0];
                } else {
                    ElMessage.error("failed to get current machine user info: " + resp.msg)
                }
            }).catch((err: any) => {
                ElMessage.error("failed to get current machine user info:" + err.msg)
            })

        } else {
            ElMessage.error("failed to get machine users info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machine users info:" + err.msg)
    })
})

function querySuggestions(query: string, callback: Function) {
    let result: any[] = []

    allUser.value.forEach((item: any) => {
        if (item.Username.indexOf(query) === 0) {
            result.push({ "value": item.Username })
        }
    })
    callback(result)
}

function onSelectUser(name: any) {
    allUser.value.forEach((item: any) => {
        if (item.Username === name.value) {
            userInfo.value = item
        }
    })
}

</script>
<style lang="scss" scoped>
.content {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;

    .users {
        width: 100%;
        height: 20%;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .current {
            width: 20%;
            height: 100%;
            font-weight: bold;
            color: rgb(92, 85, 85);
            border: 1px solid rgb(236, 235, 255);
            background: rgb(236, 235, 255);
            border-radius: 10px;
            display: flex;
            align-items: center;

            svg {
                width: 36%;
                height: 100%;
            }
        }

        .search {
            width: 70%;
        }
    }

    .info {
        width: 100%;
        height: 80%;
        display: flex;
        flex-direction: column;

        .title {
            width: 30%;
            margin: 2% 0;
        }
    }
}
</style>