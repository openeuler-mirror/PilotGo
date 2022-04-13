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
  Date: 2022-03-22 11:38:10
  LastEditTime: 2022-04-13 11:36:16
 -->
<template>
  <div class="overview">
    <div class="content">
      <!-- 选择区 -->
      <div class="choice">
        <el-form ref="form" :model="form">
          <el-form-item label="机器 IP:">
            {{macIp}}
            <!-- <el-select v-model="form.IP" multiple placeholder="请选择">
              <el-option
                v-for="item in ips"
                :key="item.ID"
                :label="item.ip"
                :value="item.ip"
              >
              </el-option>
            </el-select> -->
          </el-form-item>
          <el-form-item label="监控时间:">
            <el-row :gutter="20">
              <el-col :span="7">
              <el-date-picker type="date" placeholder="开始日期" v-model="form.dateSD" style="width: 49%;"></el-date-picker>
              <el-time-picker type="date" placeholder="开始时间" v-model="form.dateST" style="width: 49%;"></el-time-picker>
            </el-col>
            <el-col class="line" :span="1">-</el-col>
            <el-col :span="7">
              <el-date-picker placeholder="结束日期" v-model="form.dateED" style="width: 49%;"></el-date-picker>
              <el-time-picker placeholder="结束时间" v-model="form.dateET" style="width: 49%;"></el-time-picker>
            </el-col>
              <el-col :span="4">
                <el-button type="primary" @click="handleConfirm"> 确认 </el-button>
                <el-button @click="resetForm">重置</el-button>
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item label="新增指标:">
            <el-select v-model="prome" placeholder="请选择" @click="handleAppend">
              <el-option v-for="item in promes"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <!-- 图表展示区 -->
      <div class="charts flex">
        <cpu-chart class="space"  v-show="cpuShow" @close="handleClose" ref="cpuchart"></cpu-chart>
        <mem-chart class="space"  v-show="memShow" @close="handleClose"  ref="memchart"></mem-chart>
        <disk-chart class="space"  v-show="diskShow" @close="handleClose"  ref="diskchart"></disk-chart>
        <net-chart class="space"  v-show="netShow" @close="handleClose"  ref="netchart"></net-chart>
      </div>
    </div>
  </div>
</template>

<script>
import CpuChart from './echarts/cpu.vue';
import MemChart from './echarts/memory.vue';
import DiskChart from './echarts/diskIO.vue';
import NetChart from './echarts/network.vue';
export default {
  name: "Prometheus",
  components: {
    CpuChart,
    MemChart,
    DiskChart,
    NetChart,
  },
  data() {
    return {
      macIp: 'localhost:9090',
      ips: [],
      prome: '',
      promes: [
        {
          label: 'cpu',
          value: 1
        },
        {
          label: 'memroy',
          value: 2
        },
        {
          label: 'disk',
          value: 3
        },
        {
          label: 'network',
          value: 4
        }
      ],
      cpuShow: true,
      memShow: true,
      diskShow: false,
      netShow: false,
      label: 'data',
      chartW: 0,
      chartH: 0,
      form: {
        IP: '',
        dateSD: '',
        dateST: '',
        dateED: '',
        dateET: '',
      }
    };
  },
  mounted() {

    this.macIp = this.$store.getters.selectIp && this.$store.getters.selectIp.split(':')[0] || 'localhost:9090';
    this.chartW = document.getElementsByClassName("charts")[0].clientWidth/2.1;
    this.chartH = document.getElementsByClassName("charts")[0].clientHeight/1.6;
    this.$refs.cpuchart.resize({width:this.chartW,height: this.chartH});
    this.$refs.memchart.resize({width:this.chartW,height: this.chartH});   
  },
  methods: {
    handleAppend(key) {
      switch (key) {
        case 1:
          this.cpuShow = true;
          break;
        case 2:
          this.memShow = true;
          break;
        case 3:
          this.diskShow = true;
          this.$refs.diskchart.resize({width:this.chartW,height: this.chartH});
          break;
        case 4:
          this.netShow = true;
          this.$refs.netchart.resize({width:this.chartW,height: this.chartH});
          break;
      
        default:
          break;
      }
    },
    handleClose(value) {
      switch (value) {
        case 1:
          this.cpuShow = false;
          break;
        case 2:
          this.memShow = false;
          break;
        case 3:
          this.diskShow = false;
          break;
        case 4:
          this.netShow = false;
          break;
      
        default:
          break;
      }
    },
    handleConfirm() {
      let sTime = new Date(this.form.dateSD).getFullYear()+'-'+
          (new Date(this.form.dateSD).getMonth()+ 1) +'-'+new Date(this.form.dateSD).getDate()+' '+
          new Date(this.form.dateST).getHours()+':'+new Date(this.form.dateST).getMinutes()+':'+
          new Date(this.form.dateST).getSeconds();
      let eTime = new Date(this.form.dateED).getFullYear()+'-'+
          (new Date(this.form.dateED).getMonth()+ 1) +'-'+new Date(this.form.dateED).getDate()+' '+
          new Date(this.form.dateET).getHours()+':'+new Date(this.form.dateET).getMinutes()+':'+
          new Date(this.form.dateET).getSeconds();
      this.$refs.cpuchart.getCpu({starttime: parseInt(new Date(sTime)/1000)+'', endtime: parseInt(new Date(eTime)/1000)+''})
      this.$refs.memchart.getMem({starttime: parseInt(new Date(sTime)/1000)+'', endtime: parseInt(new Date(eTime)/1000)+''})
      this.$refs.diskchart.getDisk({starttime: parseInt(new Date(sTime)/1000)+'', endtime: parseInt(new Date(eTime)/1000)+''})
      this.$refs.netchart.getNet({starttime: parseInt(new Date(sTime)/1000)+'', endtime: parseInt(new Date(eTime)/1000)+''})
    },
    resetForm() {
        this.form.dateSD = '';
        this.form.dateST = '';
        this.form.dateED = '';
        this.form.dateET = '';
      }
  },
  watch: {
    prome: function(newValue) {
      if(newValue) {
        this.handleAppend(newValue);
      }
    }
  }
};
</script>

<style scoped lang="scss">
.overview {
  width: 100%;
  height: 100%;
  margin: 0 10px;
  .content {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: space-around;
    flex-direction: column;
    .flag_header {
      width: 16%;
      height: 8%;
    }
    .choice {
      width: 92%;
      margin: 0 auto;
      height: 20%;
      input {
        cursor: pointer;
      }
      .line {
        text-align: center;
      }
    }
    .charts {
      width: 100%;
      height: 70%;
      flex-wrap: wrap;
      // flex-direction: column;
      overflow-y: auto;
      .space {
        margin-bottom: 2%;
      }
    }
  }
}
</style>
