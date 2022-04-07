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
  Date: 2022-03-22 16:02:18
  LastEditTime: 2022-04-06 18:43:43
 -->
<template>
  <div class="panel">
    <span @click="handleClose" class="closeChart">✕</span>
    <div id="cpu">
    </div>
  </div>
</template>
<script>
import { getData } from "@/request/overview";
export default {
  data() {
    return {
      macIp: '',
      cpuChart: {},
      cpuData: [],
      now: new Date()/1000,
    }
  },
  mounted() {
    this.macIp = this.$store.getters.selectIp;
    this.cpuChart = this.$echarts.init(document.getElementById('cpu'))
    this.getCpu({starttime: parseInt(this.now - 60*60*6) + '', endtime: parseInt(this.now - 0) + ''});
  },
  computed: {
    option() {
      return {
        title: {text:'cpu使用率'},
        tooltip: {
          trigger: 'axis',
          position: [10, 60],
          formatter: function (params) {
            params = params[0];
            return (
              params.value[0] + ' ' + 
              params.value[1] + '%'
            );
          },
          axisPointer: {
            animation: false
          }
        },
        xAxis: {
          type: 'time',
          splitLine: {
            show: false
          }
        },
        yAxis: {
          type: 'value',
          max: 100,
          min: 0,
          boundaryGap: [0, '100%'],
        },
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 20
          },
          {
            start: 0,
            end: 20
          }
        ],
        series: [{
          name: 'cpu',
          type: 'line',
          smooth: false,
          showSymbol: false,
          lineStyle: {width: 1},
          areaStyle: {
            opacity: 0.1,
            color:  '#37A2FF'
          },
          data: this.cpuData,
        }]
      }
    },
  },
  methods: {
    resize(params) {
      this.cpuChart.resize(params)
      this.cpuChart.setOption(this.option,true)
    },
    getCpu(timeRange) {
      let params= {
        machineip: this.macIp,
        query: 1,
      }
      getData({...params, ...timeRange}).then(res => {
        if(res.data.code === 200) {
          res.data.data.forEach(item => {
              this.cpuData.push({
                time: item.time,
                value: [ item.time,parseInt(item.value).toFixed(2)]
              })
            })
        }
      })
    },
    handleClose() {
      this.$emit('close',1);
    }
  },
  watch: {
    cpuData: function() {
      this.cpuChart.setOption(this.option,true)
    }
  }
}
</script>
<style scoped lang="scss">
  .panel {
    position: relative;
    .closeChart {
      display: inline-block;
      width: 4px;
      height: 4px;
      position: absolute;
      top: 2%;
      right: 4%;
      z-index: 1;
      cursor: pointer;
    }
    .closeChart:hover {
      color: rgb(0, 163, 217)
    }
  }
  
</style>