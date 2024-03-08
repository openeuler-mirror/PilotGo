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
  LastEditTime: 2022-06-23 11:12:57
 -->
<template>
  <div id="memory" ref="chartDom"></div>
</template>
<script>
import { getData, getCurrData } from "@/request/overview";
import { formatDate } from '@/utils/dateFormat';
export default {
  data() {
    return {
      x: 0,
      y: 0,
      width: 700,
      height: 400,
      memChart: {},
      memData: [],
      now: new Date()/1000,
    }
  },
  mounted() {
    if(this.$store.getters.selectIp){
      this.getAllData({starttime: parseInt(this.now - 2*60*60) + '', endtime: parseInt(this.now - 0) + ''})
    }
    this.CreateChart();
  },
  computed: {
    option() {
      return {
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
        grid:{
          top:'4%',
          left: '5%',
          right: '3%',
        },
        xAxis: {
          type: 'time',
          splitLine: {
            show: false
          }
        },
        yAxis: {
          type: 'value',
          min: 0,
          boundaryGap: [0, '100%'],
        },
        dataZoom: [
          {
            start: 0,
            end: 100
          },
          {
            type: 'inside',
            start: 0,
            end: 100
          }
        ],
        series: [{
          name: 'memory',
          type: 'line',
          smooth: false,
          showSymbol: false,
          lineStyle: {width: 1},
          areaStyle: {
            opacity: 0.1,
            color:  '#37A2FF'
          },
          data: this.memData,
        }]
      }
    },
  },
  methods: {
    CreateChart(){
        this.memChart = this.$echarts.init(this.$refs.chartDom)
        setTimeout (()=>{
            this.$nextTick(() => {
                this.memChart.resize()
            })
        },0)
    },
    sizechange(params){
      this.memChart.resize(params) 
    },
    getAllData(timeRange) {
      let params= {
        query: '(1-(node_memory_MemAvailable_bytes/(node_memory_MemTotal_bytes)))*100',
        start: timeRange.starttime,
        end: timeRange.endtime,
        step: '10s'
      }
      getData(params).then(res => {
        this.memData = [];
        if(res.data.status === 'success' && res.data.data.result.length > 0) {
          res.data.data.result
            .filter(item => item.metric.instance === this.$store.getters.selectIp)[0]
            .values.forEach(item => {
              let localTime = formatDate(item[0], 'yyyy-MM-dd hh:mm:ss');
              this.memData.push({
                time: item[0]*1000,
                value: [localTime, parseInt(item[1]).toFixed(2)]
              })
            })
        }
      })
    },
  },
  watch: {
    memData: function() {
      this.$nextTick(() => {
        this.memChart.setOption(this.option,true)
      })
    }
  }
}
</script>
<style scoped lang="scss">
  #memory {
    width: 100%;
    height: 100%;
  }
  
</style>