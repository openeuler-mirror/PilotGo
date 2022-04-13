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
  LastEditTime: 2022-04-13 14:29:59
 -->
<template>
 <div>
   <el-table
    :data="tableData"
    :header-cell-style="hStyle"
    style="width: 100%">
    <el-table-column type="expand">
        <template slot-scope="props">
          <el-table :data="props.row.nic" :header-cell-style="childHstyle">
            <el-table-column label="IP地址" prop="IPAddr"></el-table-column>
            <el-table-column label="Mac地址" prop="MacAddr"></el-table-column>
          </el-table>  
        </template>
    </el-table-column>
    <el-table-column
      style="background:rgb(109, 123, 172);"
      label="网卡名称"
      prop="Name">
    </el-table-column>
    <el-table-column
      label="接收字节"
      prop="BytesRecv">
    </el-table-column>
    <el-table-column
      label="发送字节"
      prop="BytesSent">
    </el-table-column>
    <el-table-column
      label="接收包"
      prop="PacketsRecv">
    </el-table-column>
    <el-table-column
      label="发送包"
      prop="PacketsSent">
    </el-table-column>
  </el-table>
 </div>
</template>
<script>
import { getNetwork, getNetNic, getNetTcp, getNetUdp } from '@/request/cluster';
export default {
  name: "NetworkInfo",
  data() {
    return {
      hStyle: {
        background:'rgb(109, 123, 172)',
        color:'#fff',
        textAlign:'center',
        padding:'0',
        height: '46px',
        border: '1px solid #fff'
      },
      childHstyle: {
        background:'rgba(109, 123, 172,.6)',
        color:'#fff',
        textAlign:'center',
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
    getNetNic({uuid:this.$route.params.detail}).then(res =>{
      if(res.data.code === 200) {
        res.data.data.net_nic.forEach(item => {
          this.tableData.forEach(net => {
            if(net.Name === item.Name) {
              net.nic.push(item)
            }
          })
        });
      }
    })
    }
  }
}
</script>
<style scoped lang="scss">
.demo-table-expand {
    font-size: 0;
  }
  .demo-table-expand label {
    width: 90px;
    color: #99a9bf;
  }
  .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    width: 50%;
  }
</style>