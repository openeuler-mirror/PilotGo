<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="content">
    <div class="users">
      <div class="current">
        <img style="width:66px;" src="../../../assets/user.png" alt="user.png">
        当前用户： <span style="font-weight: bold;">{{ currentUser.Username }}</span>
      </div>
      <div class="search">
        <el-autocomplete style="width:50%" class="inline-input" v-model="userName" placeholder="请输入用户名"
          :fetch-suggestions="querySuggestions" @select="onSelectUser"></el-autocomplete>
        &nbsp;
        <el-button plain type="primary">搜索</el-button>
      </div>
    </div>
    <div style="width:30%;">
      <el-divider content-position="left">用户信息详情</el-divider>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户名："> {{ userInfo.Username }} </el-descriptions-item>
        <el-descriptions-item label="用户ID："> {{ userInfo.UserId }} </el-descriptions-item>
        <el-descriptions-item label="用户组ID："> {{ userInfo.GroupId }} </el-descriptions-item>
        <el-descriptions-item label="家目录："> {{ userInfo.HomeDir }} </el-descriptions-item>
        <el-descriptions-item label="shell类型："> {{ userInfo.ShellType }} </el-descriptions-item>
        <el-descriptions-item label="描述："> {{ userInfo.Description }} </el-descriptions-item>
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
    width: 50%;
    height: 140px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;

    .current {
      display: flex;
      align-items: center;
    }
  }
}
</style>