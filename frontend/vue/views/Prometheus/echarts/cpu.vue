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
  LastEditTime: 2022-06-24 17:13:47
 -->
<template>
  <div id="cpu" ref="chartDom"></div>
</template>
<script>
import { getData } from "@/request/overview";
import { formatDate } from '@/utils/dateFormat';
export default {
  data() {
    return {
      colorList:['#37A2FF','#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', ],
      querySQL: '100-(avg by(instance)(irate(node_cpu_seconds_total{mode="idle"}[5m]))*100)',
      cpuChart: {},
      cpuData: [],
      now: new Date()/1000,
    }
  },
  mounted() {
    this.CreateChart();
    if(this.$store.getters.selectIp) {
      this.getAllData({starttime: parseInt(this.now - 60*60*2) + '', endtime: parseInt(this.now - 0) + ''});
    }
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
          splitLine: {
            show: true,
          },
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
          name: 'cpu',
          type: 'line',
          smooth: false,
          showSymbol: false,
          lineStyle: {
            // color: 'rgb(92, 123, 217)',
            width: 1
          },
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
    CreateChart(){
      this.cpuChart = this.$echarts.init(this.$refs.chartDom)
      setTimeout (()=>{
        this.$nextTick(() => {
          this.cpuChart.resize()
          // this.cpuChart.resize({width: 300,height:200})
        })
      },0)
    },
    sizechange(params){
      this.cpuChart.resize(params) 
    },
    getAllData(timeRange) {
      let params= {
        query: '100-(avg by(instance)(irate(node_cpu_seconds_total{mode="idle"}[5m]))*100)',
        start: timeRange.starttime,
        end: timeRange.endtime,
        step: '10s'
      }
      getData(params).then(res => {
        this.cpuData = [];
        if(res.data.status === 'success' && res.data.data.result.length > 0) {
          res.data.data.result
            .filter(item => item.metric.instance === this.$store.getters.selectIp)[0]
            .values.forEach(item => {
              let localTime = formatDate(item[0], 'yyyy-MM-dd hh:mm:ss');
              this.cpuData.push({
                time: item[0]*1000,
                value: [localTime, parseInt(item[1]).toFixed(2)]
              })
            })
        }
      })
    },
  },
  watch: {
    cpuData: function() {
      this.$nextTick(() => {
        this.cpuChart.setOption(this.option,true)
      })
    }
  }
}
</script>
<style scoped lang="scss">
  #cpu {
    width: 100%;
    height: 100%;
  }
  
</style>