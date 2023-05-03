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
  Date: 2022-04-11 13:07:46
 LastEditTime: 2023-04-18 09:46:14
 -->
<template>
  <div class="content">
    <div class="operation">
     <el-button plain  type="primary" size="medium" @click="handleEdit">编辑</el-button>
   </div><br/>
   <h3>IPv4:</h3><br/>
   <el-descriptions :column="2" size="medium" border>
    <el-descriptions-item label="IP分配">{{ net.BOOTPROTO === 'dhcp'? '动态DHCP' : '手动' }}</el-descriptions-item>
    <el-descriptions-item label="IPv4地址">{{ net.IPADDR }}</el-descriptions-item>
    <el-descriptions-item label="IPv4子网掩码">{{ net.NETMASK || '无'}}</el-descriptions-item>
    <el-descriptions-item label="IPv4网关">{{ net.GATEWAY }}</el-descriptions-item>
    <el-descriptions-item label="IPv4首选DNS">{{net.DNS1}}</el-descriptions-item>
    <el-descriptions-item label="IPv4备选DNS">{{net.DNS2 || '无'}}</el-descriptions-item>
    </el-descriptions>

   <el-dialog 
    :title="title"
    top="10vh"
    :before-close="handleClose" 
    :visible.sync="display" 
    width="70%"
  >  
    <update-form v-if="type === 'update'" :net="net" @click="handleClose"></update-form>      
  </el-dialog>
 </div>
</template>
<script>
import { getNetwork } from '@/request/cluster';
import UpdateForm from "../form/network/updateForm";
export default {
  name: "NetworkInfo",
  components: {
    UpdateForm
  },
  data() {
    return {
      title: '',
      display: false,
      type: '',
      net: {
        // BOOTPROTO: '手动',
        // IPADDR: '172.17.127.29',
        // NETMASK: '24',
        // GATEWAY: '172.17.127.252',
        // DNS1: '123.150.150.150',
        // DNS2: '255.5.5.5',
      },
      netData: [],
      tableData: []
    }
  },
  mounted() {
    if(this.$route.params.detail != undefined) {
    getNetwork({uuid:this.$route.params.detail}).then(res => {
      if(res.data.code === 200) {
        this.net = res.data.data;
      } else {
        console.log(res.data.msg)
      }
    })
    }
  },
  methods: {
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = "";
    },
    handleEdit() {
      this.display = true;
      this.title = "修改网络配置";
      this.type = "update";
      this.row = this.net;
    }
  }
}
</script>
<style scoped lang="scss">
.content {
  width:96%; 
  padding-top:20px; 
  margin: 0 auto;
  .operation {
    width: 100%;
    text-align: right;
  }
  h3 {
    color: rgb(145, 139, 139);
  }
}
</style>