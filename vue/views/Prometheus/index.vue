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
  LastEditTime: 2022-03-28 15:32:47
 -->
<template>
  <div class="overview">
    <div class="content">
      <!-- 头部标签 -->
      <div class="flag_headedr">
        <span class="iconfont">&#xe663;</span>
        <span class="level-title">机器监控数据</span>
        <button @click="handleAppend"> 显示 </button>
      </div>
      <!-- 选择区 -->
      <div class="choice">
        <el-form ref="form" :model="form">
          <el-form-item label="监控时间:">
            <el-row>
              <el-col :span="7">
              <el-date-picker type="date" placeholder="开始日期" v-model="form.dateSD" style="width: 49%;"></el-date-picker>
              <el-time-picker type="date" placeholder="开始时间" v-model="form.dateST" style="width: 49%;"></el-time-picker>
            </el-col>
            <el-col class="line" :span="2">-</el-col>
            <el-col :span="7">
              <el-date-picker placeholder="结束日期" v-model="form.dateED" style="width: 49%;"></el-date-picker>
              <el-time-picker placeholder="结束时间" v-model="form.dateET" style="width: 49%;"></el-time-picker>
            </el-col>
            </el-row>
          </el-form-item>
          <el-form-item label="新增指标:">
            <el-select v-model="prome" placeholder="请选择">
              <el-option v-for="item in promes"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>

        <!-- <el-button>确定</el-button> -->

      </div>
      <!-- 图表展示区 -->
      <div class="charts flex">
        <cpu-chart v-show="cpuShow" @close="handleClose" ref="cpuchart" :macIp="macIp"></cpu-chart>
        <mem-chart v-show="memShow" @close="handleClose"  ref="memchart" :macIp="macIp"></mem-chart>
        <disk-chart v-show="diskShow" @close="handleClose"  ref="diskchart" :macIp="macIp"></disk-chart>
        <net-chart v-show="netShow" @close="handleClose"  ref="netchart" :macIp="macIp"></net-chart>
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
      prome: '',
      promes: [
        {
          label: 'cpu',
          value: '1'
        },
        {
          label: 'memroy',
          value: '2'
        },
        {
          label: 'disk',
          value: '3'
        },
        {
          label: 'network',
          value: '4'
        }
      ],
      cpuShow: true,
      memShow: true,
      diskShow: false,
      netShow: false,
      macIp: 'localhost:9100',
      label: 'data',
      chartW: 0,
      chartH: 0,
      options: ['123.12.11.23'],
      form: {
        dateSD: '',
        dateST: '',
        dateED: '',
        dateET: '',
      }
    };
  },
  mounted() {
    this.macIp = this.$route.query.ip || '';
    this.chartW = document.getElementsByClassName("charts")[0].clientWidth/2.4;
    this.chartH = document.getElementsByClassName("charts")[0].clientHeight/2;
    this.$refs.cpuchart.resize({width:this.chartW,height: this.chartH});
    this.$refs.memchart.resize({width:this.chartW,height: this.chartH});
    this.$refs.diskchart.resize({width:this.chartW,height: this.chartH});
    this.$refs.netchart.resize({width:this.chartW,height: this.chartH});
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
          break;
        case 4:
          this.netShow = true;
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
      width: 98%;
      height: 14%;
    }
    .charts {
      width: 98%;
      height: 70%;
      flex-flow: wrap;
      flex-direction: column;
      align-content: space-around;
      overflow-y: auto;
    }
  }
}
</style>
