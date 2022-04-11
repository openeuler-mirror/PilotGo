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
  Date: 2022-04-11 12:47:34
  LastEditTime: 2022-04-11 16:13:19
 -->
<template>
 <div class="content" style="width:96%; padding-top:20px; margin: 0 auto">
   <el-descriptions :column="3" size="medium" border>
     <el-descriptions-item labelStyle="width:20%;background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px;" label="服务名">状态</el-descriptions-item>
     <el-descriptions-item labelStyle="width:20%;background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px;" label="服务名">状态</el-descriptions-item>
     <el-descriptions-item labelStyle="width:20%;background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px;" label="服务名">状态</el-descriptions-item>
    <el-descriptions-item 
      v-for="item in serviceData" 
      :key="item.$index" 
      :label="item.Name">  
      <span class="statusBtn">{{ item.Active === 'active'? '正在运行': '已停止' }}</span><br/><br/>
        <el-button class="smallBtn" size="mini" plain type="primary" @click="handleStart(item.Name)">启动</el-button>
        <el-button class="smallBtn" size="mini" plain type="primary" @click="handleStop(item.Name)">停止</el-button>
        <el-button class="smallBtn" size="mini" plain type="primary" @click="handleRestart(item.Name)">重启</el-button>
    </el-descriptions-item>
   </el-descriptions>
 </div>
</template>
<script>
import { getserviceList,serviceStart, serviceStop, serviceRestart } from '@/request/cluster';
export default {
  name: "ServiceInfo",
  data() {
    return {
      serviceData: [],
      userName: ''
    }
  },
  mounted() {
    this.userName = this.$store.getters.userName;
    getserviceList({uuid:this.$route.params.detail}).then((res) => {
      this.serviceData = res.data.data.service_list;
    })
  },
  methods: {
    handleStart(sericeName) {
      serviceStart({uuid:this.$route.params.detail, userName:this.userName, service: sericeName}).then(res => {
        if(res.data.code === 200) {
          this.$message.success(res.data.msg);
        } else {
           this.$message.error(res.data.msg);
        }
      })
    },
    handleStop(sericeName) {
      serviceStop({uuid:this.$route.params.detail,userName:this.userName, service: sericeName}).then(res => {
        if(res.data.code === 200) {
          this.$message.success(res.data.msg);
        } else {
           this.$message.error(res.data.msg);
        }
      })
    },
    handleRestart(sericeName) {
      serviceRestart({uuid:this.$route.params.detail,userName:this.userName, service: sericeName}).then(res => {
        if(res.data.code === 200) {
          this.$message.success(res.data.msg);
        } else {
           this.$message.error(res.data.msg);
        }
      })
    },
  }
}
</script>
<style scoped lang="scss">
.statusBtn {
  display: inline-block;
  width:70px; 
  font-size: 12px;
  border-radius:11px; 
  background: rgb(109, 123, 172);
  color:#fff;
}
.smallBtn {
  padding: 6px;
  margin-left: 1%;
}
</style>
