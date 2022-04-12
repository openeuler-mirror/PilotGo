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
  LastEditTime: 2022-04-12 11:06:54
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
      serviceData: [
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-block-8:3.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-cdrom.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-ata\\x2dVMware_Virtual_SATA_CDRW_Drive_01000000000000000001.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-dm\\x2dname\\x2dklas\\x2droot.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-dm\\x2dname\\x2dklas\\x2dswap.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-dm\\x2duuid\\x2dLVM\\x2d7xxzeEXx4pYZGobYEfVfDLo94GpvRW1Do9O9Aopt7HBqBeX36lYS8c5gdWYNnTEf.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-dm\\x2duuid\\x2dLVM\\x2d7xxzeEXx4pYZGobYEfVfDLo94GpvRW1DQxrJ3cZcwYaWM21OZeaAzAcdEp380BQY.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-lvm\\x2dpv\\x2duuid\\x2d0f1bNE\\x2dtMeE\\x2d0dId\\x2dGqIC\\x2d2iw1\\x2duudL\\x2drsnHGQ.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2did-lvm\\x2dpv\\x2duuid\\x2dsmpxkP\\x2dObgV\\x2dwEo1\\x2dUcj8\\x2d5JPC\\x2dzL2y\\x2dNv81Ul.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpartuuid-c2e6d7a8\\x2d01.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpartuuid-c2e6d7a8\\x2d02.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpartuuid-c2e6d7a8\\x2d03.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpath-pci\\x2d0000:00:10.0\\x2dscsi\\x2d0:0:0:0.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpath-pci\\x2d0000:00:10.0\\x2dscsi\\x2d0:0:0:0\\x2dpart1.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpath-pci\\x2d0000:00:10.0\\x2dscsi\\x2d0:0:0:0\\x2dpart2.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpath-pci\\x2d0000:00:10.0\\x2dscsi\\x2d0:0:0:0\\x2dpart3.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2dpath-pci\\x2d0000:02:05.0\\x2data\\x2d2.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2duuid-020f054b\\x2d172b\\x2d49ac\\x2da4a9\\x2d86546a953365.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2duuid-3f250f7e\\x2d202b\\x2d4b1a\\x2da40b\\x2d32772238f5ac.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-disk-by\\x2duuid-d5246dd2\\x2dc66f\\x2d4109\\x2da968\\x2d805b44e7ca72.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-dm\\x2d0.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-dm\\x2d1.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-klas-root.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-klas-swap.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-mapper-klas\\x2droot.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-mapper-klas\\x2dswap.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-rfkill.device",
        "SUB": "plugged"
      },
      {
        "Active": "active",
        "LOAD": "loaded",
        "Name": "dev-sda.device",
        "SUB": "plugged"
      }],
      userName: ''
    }
  },
  mounted() {
    this.userName = this.$store.getters.userName;
    if(this.$route.params.detail != undefined) {
    getserviceList({uuid:this.$route.params.detail}).then((res) => {
      this.serviceData = res.data.data.service_list;
    })
    }
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
