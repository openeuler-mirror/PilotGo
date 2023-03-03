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
  LastEditTime: 2022-06-29 10:11:11
 -->
<template>
  <div class="overview">
    <div class="content">
      <!-- 选择区 -->
      <div class="choice">
        <el-form ref="form" :inline="true" class="flex">
          <el-form-item label="机器 IP:">
            <el-autocomplete
              style="width:100%"
              class="inline-input"
              v-model="macIp"
              :fetch-suggestions="querySearch"
              placeholder="请输入ip关键字"
              @select="handleSelect"
            ></el-autocomplete>
          </el-form-item>
          <el-form-item label="监控时间:">
            <el-date-picker
              v-model="dateRange"
              type="datetimerange"
              prefix-icon="el-icon-date"
              :picker-options="pickerOptions"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              align="center"
              size="large"
              :default-time="['9:00:00', '18:00:00']"
              @change="changeDate">
            </el-date-picker>
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
          <el-form-item>
            <el-button plain type="primary" @click="handleConfirm">确认</el-button>
          </el-form-item>
        </el-form>
      </div>
      <!-- 图表展示区 -->
      <div class="charts">
        <grid-layout :layout.sync="layout"
          :col-num="12"
          :row-height="10"
          :auto-size="true"
          :is-draggable="draggable"
          :is-resizable="resizable"
          :vertical-compact="true"
          :use-css-transforms="true"
          >
          <template>
            <div v-for="(item, i) in layout" :key="i">
            <grid-item
              class="panel"
              :static="item.static"
              :x="item.x"
              :y="item.y"
              :w="item.w"
              :h="item.h"
              :min-w="4"
              :min-h="8"
              :i="item.i"
              :key="item.i"
              :ref="`chartParent${item.i}`"
              drag-allow-from=".drag"
              drag-ignore-from=".noDrag"
              @resize="SizeAutoChange" 
              @resized="SizeAutoChange"
              v-if="item.display"
            >
              <div class="drag" title="可任意拖动">
                <span class="title">{{item.title}}</span>
                <span @click="item.display = false" class="closeChart">✕</span>
              </div>
              <component class="noDrag" style="width:100%;height:100%; padding-top:20px" :is="item.component" :ref="`chart${item.i}`"/>
            </grid-item>
            </div>
          </template>
        </grid-layout>
      </div>
    </div>
  </div>
</template>

<script>
import CpuChart from './echarts/cpu.vue';
import MemChart from './echarts/memory.vue';
import DiskChart from './echarts/diskIO.vue';
import NetChart from './echarts/network.vue';
import merge from 'webpack-merge';
import { pickerOptions } from './datePicker';
import { GridLayout, GridItem } from "vue-grid-layout";
import { getallMacIps } from '@/request/cluster';
export default {
  name: "Prometheus",
  components: {
    CpuChart,
    MemChart,
    DiskChart,
    NetChart,
    GridLayout,
    GridItem,
  },
  data() {
    return {
      layout: [
          {"x":0,"y":0,"w":6,"h":16,"i":"0", static: false, display: true, component:'CpuChart',title:'cpu使用率'},
          {"x":6,"y":0,"w":6,"h":16,"i":"1", static: false, display:false, component:'MemChart',title:'内存使用率'},
          {"x":0,"y":16,"w":6,"h":20,"i":"2", static: false, display:false, component:'DiskChart',title:'磁盘读写速率'},
          {"x":6,"y":16,"w":6,"h":20,"i":"3", static: false, display:false, component:'NetChart',title:'网络平均输入输出速率'},
      ],
      draggable: true,
      resizable: true,
      pickerOptions: pickerOptions,
      macIp: '',
      ips: [],
      ipData: [],
      label: 'data',
      dateRange: [new Date()-2*60*60*1000, new Date()-0],
      timeParams: {
        starttime: (new Date()-2*60*60*1000)/1000+'',
        endtime: (new Date()-0)/1000+''
      },
      prome: '',
      promes: [
        {
          label: 'cpu',
          value: 1,
          querySQL: '100-(avg by(instance)(irate(node_cpu_seconds_total{mode="idle"}[5m]))*100)',
        },
        {
          label: 'memory',
          value: 2,
          querySQL: '(1-(node_memory_MemAvailable_bytes/(node_memory_MemTotal_bytes)))*100',
        },
        {
          label: 'disk',
          value: 3,
          querySQL:''
        },
        {
          label: 'network',
          value: 4,
          querySQL:''
        }
      ],
    };
  },
  activated() {
    window.addEventListener("resize", this.resize);
  },
  mounted() {
    this.macIp = this.$store.getters.selectIp && this.$store.getters.selectIp.split(':')[0];
    getallMacIps().then(res => {
      this.ips = [];
      this.ipData = [];
      if(res.data.code === 200) {
        this.ips = res.data.data && res.data.data;
        this.ips.forEach(item => {
            this.ipData.push({'value':item.ip_dept,'ip':item.ip})
          })
      }
    })
    window.addEventListener("resize", this.resize);
  },
  methods: {
    querySearch(queryString, cb) {
      var ipData = this.ipData;
      var results = queryString ? ipData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): ipData;
      cb(results);
    },
    SizeAutoChange(i,x,y,newH,newW) {
      var EchartId = `chart${i}`; 
      eval(`this.$refs.${EchartId}[0].sizechange({width:${newW},height:${newH*0.92}})`)
    },
    resize(){
      let _this = this;
      this.layout.forEach((i) => {
        if(i.display) {
          let EchartId = `chart${i.i}`;
          let EchartPId = `chartParent${i.i}`;
          let initW = parseInt(eval(`_this.$refs.${EchartPId}[0]`).style.width.split('p')[0]);
          let initH = parseInt(eval(`_this.$refs.${EchartPId}[0]`).style.height.split('p')[0]);
          eval(`_this.$refs.${EchartId}[0]`).sizechange({ width: initW, height: initH*0.9 });
        }
      })
    },  
    handleAppend(key) {
      this.layout[key-1].display = true;
    },
    handleSelect(item) {
      this.macIp = item && item.ip;
    },
    changeDate(timeValue) {
      this.timeParams = {
        starttime: timeValue[0]/1000+'',
        endtime: timeValue[1]/1000+''
      }
    },
    handleConfirm() {
      this.$store.dispatch('setSelectIp', this.macIp);
      this.$router.push({
        query:merge(this.$route.query,{'ip': this.macIp})
      })
      let _this = this;
      this.layout.forEach((i) => {
        if(i.display) {
           var EchartId = `chart${i.i}`; 
          eval(`_this.$refs.${EchartId}[0]`).getAllData(this.timeParams);
        }
      })
    },
  },
  watch: {
    prome: function(newValue) {
      if(newValue) {
        this.handleAppend(newValue);
      }
    },
    '$route': {
      handler() {
        if(this.$route.name) {
          console.log("监控ip",this.$store.getters.selectIp)
        }
      }
    }
  },
  deactivated() {
    window.removeEventListener('resize', this.resize);
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.resize);
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
    .flag_header {
      width: 16%;
      height: 8%;
    }
    .choice {
      width: 96%;
      margin: 0 auto;
      height: 6%;
      background: rgba(255,255,255,0);
      input {
        cursor: pointer;
      }
      .line {
        text-align: center;
      }
    }
    .charts {
      width: 100%;
      height: 94%;
      border-top: 1px dashed #bbb;
      overflow: auto;
      .vue-grid-item:not(.vue-grid-placeholder) {
        background: #fff;
        border-radius: 6px;
        .drag{
          width:100%;
          height: 30px;
          border-radius: 6px 6px 0 0;
          position: absolute;
          z-index: 9999;
          .title {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 100%;
            height: 100%;
            color: rgb(133, 130, 130);
            font-size: 16px;
          }
          .closeChart {
            display: inline-block;
            width: 2%;
            position: absolute;
            top: 2%;
            right: 1%;
            z-index: 1;
            cursor: pointer;
          }
          .closeChart:hover {
            color:#fff;
          }
        }
        .drag:hover {
          background: rgba(242, 150, 38,.6);
        }
      }
    }
  }
}
</style>
