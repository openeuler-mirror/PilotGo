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
  Date: 2022-02-25 16:33:46
  LastEditTime: 2022-04-07 13:24:36
  Description: provide agent log manager of pilotgo
 -->
<template>
 <div class="chartContent">
    <div class="left">
      <h4>基本信息</h4>
      <div class="basic">
        <el-descriptions :column="2" size="medium" border>
          <el-descriptions-item label="平台">{{ basic.macPlatform }}</el-descriptions-item>
          <el-descriptions-item label="架构">{{ basic.mackernel }}</el-descriptions-item>
          <el-descriptions-item label="cpu">{{ basic.macCPU }}</el-descriptions-item>
          <el-descriptions-item label="内存" :span="2">{{basic.macMEM}}</el-descriptions-item>
        </el-descriptions>
      </div>
      <div id="diskChar" :style="{width:'90%',height:'90%'}"></div>
    </div>
    <div class="right">
      <div class="user">
        <h4>用户信息</h4>
        <small-table
          ref="userTab"
          :data="userData"
          :height="tHight"
        >
          <template v-slot:content>
            <el-table-column
              prop="Username"
              label="用户">
            </el-table-column>
            <!-- <el-table-column
              prop="currUser"
              label="当前用户">
            </el-table-column> -->
            <el-table-column
              prop="Description"
              label="备注">
            </el-table-column>
          </template>
        </small-table>
      </div>
      <div class="service">
        <h4>服务信息</h4>
        <small-table
          ref="userTab"
          :data="serviceData"
          :height="tHight"
        >
          <template v-slot:content>
            <el-table-column
              prop="Name"
              label="名称">
            </el-table-column>
            <el-table-column
              prop="Active"
              label="状态">
            </el-table-column>
            <el-table-column label="操作">
              <template slot-scope="scope">
                <el-button v-if="scope.row.Active == 'inactive'" size="mini" type="primary" plain 
                  @click="handleStart(scope.$index,scope.row)"> 
                  启动 </el-button>
                  <el-button v-if="scope.row.Active == 'active'" size="mini" type="primary" plain 
                  @click="handleStop()"> 
                  关闭 </el-button>
              </template>
            </el-table-column>
          </template>
        </small-table>
      </div>
      <div class="kernel">
        <h4>内核信息</h4>
        <small-table
          ref="userTab"
          :data="kernelData"
          :height="tHight"
        >
          <template v-slot:content>
            <el-table-column
              label="内核">
              <template slot-scope="scope">
                {{ Object.keys(scope.row)[0] }}  
              </template> 
            </el-table-column>
            <el-table-column
              label="个数">
              <template slot-scope="scope">
                {{ Object.values(scope.row)[0] }}  
              </template> 
            </el-table-column>
            <el-table-column label="操作">
              <template slot-scope="scope">
                <el-button size="mini" type="primary" plain 
                  @click="handleChange(scope.row)"> 
                  编辑 </el-button>
              </template>
            </el-table-column>
          </template>
        </small-table>
      </div>
    </div>
 </div>
</template>
<script>
import { getCpu, getOS, getMemory, getUser, getAllUser, 
        getserviceList, getSyskernel, getDisk, serviceStart, 
        serviceStop, changeSyskernel } from '@/request/cluster';
import SmallTable from "@/components/SmallTable";
export default {
   /* 某一台机器的详情，左右分布
      左：上 描述列表，下 饼图
      右：上 用户表，中 service,表 下 内核表
   */
  components: {
    SmallTable,
  },
  props: {
    ip: {
      type: String,
      default: '0.0.0.0'
    } 
  },
  data() {
    return {
      userName: '',
      tHight: 200,
      params: {},
      basic: {
        macPlatform: '',
        mackernel: '',
        macCPU: '',
        macMEM: ''
      },
      userData: [],
      serviceData: [],
      kernelData: [],
      diskData: [],
    }
  },
  mounted() {
    this.userName = this.$store.getters.userName;
    let obj = this.params = {uuid:this.$route.params.detail};
     getOS(obj).then((res) => {
      let result = res.data.data.os_info;
      this.basic.macPlatform = result.Platform;
      this.basic.mackernel = result.KernelArch;
    })
    
    getCpu(obj).then((res) => {
      this.basic.macCPU = res.data.data.CPU_info.CpuNum + '核';
    })
    getMemory(obj).then((res) => {
      let memTotal = 0;
      memTotal = res.data.data.memory_info.MemTotal / 1024 / 1024;
      this.basic.macMEM = memTotal.toFixed(2) + 'G';
    })
    getAllUser(obj).then((res) => {
      this.userData = res.data.data.user_all;
    })
    getserviceList(obj).then((res) => {
      this.serviceData = res.data.data.service_list;
    })
    getSyskernel(obj).then((res) => {
      this.kernelData = res.data.data.sysctl_info;
    })
    getDisk(obj).then((res) => {
      let disk = [];
      disk = res.data.data.disk_use;
      disk.filter((item) => {
        return item.usedPercent > 0;
      }).forEach(item => {
        this.diskData.push({value: item.usedPercent.toFixed(2),name: item.device})
      });;
      this.diskEchart(this.diskData);
    })
  },
  methods: {
    diskEchart(disks) {
      let diskChart = this.$echarts.init(document.getElementById('diskChar'))
      let option = {
        title: {					         	
                text: '磁盘信息',                
                textStyle:{					//---主标题内容样式	
                	color:'#000'
                },
                x: 'center',
                y: '16%',
        },
        series: [{
                    name: '磁盘信息',
                    type: 'pie',
                    radius: '55%',
                    data: disks
                  }]};
      diskChart.setOption(option);
    },
    handleStart(index,row) {
      serviceStart({...this.params,userName:this.userName}).then(res => {
        if(res.data.code === 200) {
          this.serviceData[index].Active = 'active';
          this.$message.success(res.data.msg);
        } else {
           this.$message.error(res.data.msg);
        }
      })
    },
    handleStop(index,row) {
      serviceStop({...this.params,userName:this.userName}).then(res => {
        if(res.data.code === 200) {
          this.serviceData[index].Active = 'inactive';
          this.$message.success(res.data.msg);
        } else {
           this.$message.error(res.data.msg);
        }
      })
    },
    handleChange(row) {
      this.$prompt('输入内核数量', '修改内核', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       }).then(({value}) => {
         let name = Object.keys(row)[0];
         let params = {
           args:name + '=' + value, 
           uuid: this.$route.params.uuid, 
           userName: this.userName
          }
         changeSyskernel(params).then(res => {
           if(res.data.code === 200) {
            this.kernelData[index] = {name : value};
            this.$message.success(res.data.msg);
           } else {
            this.$message.error(res.data.msg);
           }
         })
       }) 
      
    }


  },


}
</script>
<style scoped lang="scss">
.chartContent {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  h4 {
    margin: 10px;
    width: 100%;
    text-align: center;
  }
  .left {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    .basic {
      margin-bottom: 20px;
      width: 90%;
    }
  }
  .right {
    flex: 2;  
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    div {
      width: 100%;
    };
  }
}

</style>
