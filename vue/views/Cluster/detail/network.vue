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
  LastEditTime: 2022-06-20 16:45:39
 -->
<template>
  <div class="content">
    <div class="operation">
     <el-button plain  type="primary" size="medium" @click="handleEdit">编辑</el-button>
   </div><br/>
   <h3>IPv4:</h3><br/>
   <el-descriptions :column="2" size="medium" border>
    <el-descriptions-item label="IP分配">{{ net.IPdist }}</el-descriptions-item>
    <el-descriptions-item label="IPv4地址">{{ net.IPv4Addr }}</el-descriptions-item>
    <el-descriptions-item label="IPv4子网前缀长度">{{ net.IPv4SPL }}</el-descriptions-item>
    <el-descriptions-item label="IPv4网关">{{ net.IPv4gateway }}</el-descriptions-item>
    <el-descriptions-item label="IPv4首选DNS">{{net.IPv4DNS[0]}}</el-descriptions-item>
    <el-descriptions-item label="IPv4备选DNS">{{net.IPv4DNS[1]}}</el-descriptions-item>
    </el-descriptions><br/><br/><br/>
    <h3>IPv6:</h3><br/>
    <el-descriptions :column="2" size="medium" border>
      <el-descriptions-item label="IP分配">{{ net.IPdist }}</el-descriptions-item>
    <el-descriptions-item label="IPv6地址">{{ net.IPv6Addr }}</el-descriptions-item>
    <el-descriptions-item label="IPv6子网前缀长度">{{ net.IPv6SPL }}</el-descriptions-item>
    <el-descriptions-item label="IPv6网关">{{ net.IPv6gateway }}</el-descriptions-item>
    <el-descriptions-item label="IPv6首选DNS">{{net.IPv6DNS[0]}}</el-descriptions-item>
    <el-descriptions-item label="IPv6备选DNS">{{net.IPv6DNS[1]}}</el-descriptions-item>
   </el-descriptions>
   
   <el-dialog 
    :title="title"
    top="10vh"
    :before-close="handleClose" 
    :visible.sync="display" 
    width="70%"
  >  
    <update-form v-if="type === 'update'" :row="net" @click="handleClose"></update-form>      
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
        IPdist: '手动',
        IPv4Addr: '172.17.127.29',
        IPv4SPL: '24',
        IPv4gateway: '172.17.127.252',
        IPv4DNS: ['123.150.150.150','223.5.5.5'],
        IPv6Addr: '172.17.127.291',
        IPv6SPL: '4',
        IPv6gateway: '172.17.127.222',
        IPv6DNS: ['122.150.150.150','222.5.5.5'],
      },
      netData: [],
      tableData: []
    }
  },
  mounted() {
    if(this.$route.params.detail != undefined) {
    getNetwork({uuid:this.$route.params.detail}).then(res => {
      if(res.data.code === 200) {
        res.data.data.net_io.forEach(item => item.nic = []);
        this.tableData = res.data.data.net_io;
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