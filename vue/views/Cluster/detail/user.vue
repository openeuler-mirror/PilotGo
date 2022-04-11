<!-- 
  Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
  PilotGo is licensed under the Mulan PSL v2.
  You can use this software accodring to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
      http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
  EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
  Author: zhaozhenfang
  Date: 2022-04-11 12:47:07
  LastEditTime: 2022-04-11 15:00:57
 -->
<template>
 <div class="content" style="width:96%; padding-top:20px; margin: 0 auto">
   <el-descriptions :column="4" size="medium" border>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
    <el-descriptions-item 
      v-for="item in userData" 
      :key="item.UserId" 
      :label="item.Username">
      {{ item.Description }}
    </el-descriptions-item>
   </el-descriptions>
 </div>
</template>
<script>
import {  getUser, getAllUser } from '@/request/cluster';
export default {
  name: "UserInfo",
  data() {
    return {
      userData: [],
    }
  },
  mounted() {
    getAllUser({uuid:this.$route.params.detail}).then((res) => {
      this.userData = res.data.data.user_all;
    })
  }
}
</script>
<style scoped lang="scss">
</style>
