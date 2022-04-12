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
  Date: 2022-04-08 11:34:55
  LastEditTime: 2022-04-12 13:36:07
 -->
<template>
 <div class="content" style="width:96%; padding-top:20px; margin: 0 auto">
   <el-descriptions :column="2" size="medium" border>
    <el-descriptions-item label="机器IP">{{ basic.IP }}</el-descriptions-item>
    <el-descriptions-item label="所属部门">{{ basic.dept }}</el-descriptions-item>
    <el-descriptions-item label="监控状态">{{ basic.status }}</el-descriptions-item>
    <el-descriptions-item label="系统版本">{{ basic.macPlatform }}</el-descriptions-item>
    <el-descriptions-item label="架构">{{ basic.mackernel }}</el-descriptions-item>
    <el-descriptions-item label="cpu">{{ basic.macCPU }}</el-descriptions-item>
    <el-descriptions-item label="内存">{{basic.macMEM}}</el-descriptions-item>
    <el-descriptions-item label="内核版本">{{basic.osVersion}}</el-descriptions-item>
    <el-descriptions-item :label="item.device" :span="2" v-for="item in diskData" :key="item.$index">
      <span class="diskMount">{{"挂载点："+item.path+"("+item.total+")"}}</span>
      <p class="progress">
        <span :style="{width: item.usedPercent}">{{item.usedPercent}}</span>
      </p>
    </el-descriptions-item>
   </el-descriptions>
 </div>
</template>
<script>
import { getBasicInfo, getCpu, getOS, getMemory,getDisk } from '@/request/cluster';
export default {
  name: "BaseInfo",
  data() {
    return {
      params: {},
      diskData: [
      {
        "device": "/dev/dm-0",
        "fileSystem": "/dev/dm-0(挂载点：/)",
        "fstype": "xfs",
        "path": "/",
        "total": "64G",
        "used": "45G",
        "usedPercent": "71%"
      },
      {
        "device": "/dev/sda1",
        "fileSystem": "/dev/sda1(挂载点：/boot)",
        "fstype": "xfs",
        "path": "/boot",
        "total": "40G",
        "used": "0G",
        "usedPercent": "20%"
      }
    ],
      basic: {
        IP: '127.0.0.1',
        dept: '麒麟',
        manager: 'root',
        status: '在线',
        IP: '192.168.160.128',
        macPlatform: 'Kylin',
        mackernel: 'x86_64',
        macCPU: '4 Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz',
        macMEM: '2957640',
        osVersion: '4.19.90-24.4.v2101.ky10.x86_64'
      },
    }
  },
  mounted() {
    let obj = this.params = {uuid:this.$route.params.detail};
    if(this.$route.params.detail != undefined) {
      getBasicInfo(obj).then(res => {
        this.basic.IP = res.data.data.IP;
        this.basic.dept = res.data.data.depart;
        this.basic.status = res.data.data.state === 1? '在线': res.data.data.state === 2? '离线':'未分配';
      })
      getOS(obj).then((res) => {
        let result = res.data.data.os_info;
        this.basic.macPlatform = result.Platform + ' ' + result.PlatformVersion;
        this.basic.mackernel = result.KernelArch;
      })
      
      getCpu(obj).then((res) => {
        this.basic.macCPU = res.data.data.CPU_info.CpuNum + '核 ' + res.data.data.CPU_info.ModelName;
      })
      getMemory(obj).then((res) => {
        let memTotal = 0;
        memTotal = res.data.data.memory_info.MemTotal / 1024 / 1024;
        this.basic.macMEM = memTotal.toFixed(2) + 'G';
      })
      getDisk(obj).then((res) => {
        this.diskData = res.data.data.disk_use;
      })
    }
  }
}
</script>
<style scoped lang="scss">
.content {
  .diskMount {
    display: inline-block;
    font-size: 12px;
    word-break: break-all;
    width:22%;
    text-align: center;
  }
  .progress {
    display: inline-block;
    width:74%; 
    margin-left: 2%;
    border: 1px solid rgba(11, 35, 117,.5);  
    background: #fff; 
    border-radius: 10px; 
    text-align:left;
    span {
      display: inline-block;
      background: rgba(11, 35, 117,.6);
      text-align:center;
      color: #fff;
      border: 1px solid #fff;
      border-radius: 10px;
    }
  }
}
</style>
