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
  LastEditTime: 2022-06-23 11:13:49
 -->
<template>
  <div id="network" ref="chartDom"></div>
</template>
<script>
import { getData } from "@/request/overview";
import { formatDate } from '@/utils/dateFormat';
export default {
  data() {
    return {
      x: 0,
      y: 0,
      width: 700,
      height: 400,
      netChart: {},
      netIn: [],
      netOut: [],
      now: new Date()/1000,
    }
  },
  mounted() {
    if(this.$store.getters.selectIp) {
      this.getStepData('in',this.netIn,'irate(node_network_receive_bytes_total[5m])',{starttime: parseInt(this.now - 2*60*60) + '',
          endtime: parseInt(this.now - 0) + ''});
      this.getStepData('out',this.netOut,'irate(node_network_transmit_bytes_total[5m])',{starttime: parseInt(this.now - 2*60*60) + '',
          endtime: parseInt(this.now - 0) + ''});
    }
    this.CreateChart();
  },
  computed: {
    option() {
      return {
        tooltip: {
          trigger: 'axis',
        },
        legend: {
          data: [],
        },
        grid: [
          {
            bottom: '60%',
            left: '7%',
            top: '4%',
            right: '3%'
          },
          {
            top: '46%',
            left: '7%',
            right: '3%',
            bottom: '4%'
          },
        ],
        xAxis: [{
          type: 'time',
          splitLine: {
            show: false
          }
        },{
          type: 'time',
          show: false,
          gridIndex: 1,
          position: 'top',
          splitLine: {
            show: false
          }
        },
        ],
        yAxis: [{
            name: 'input(B/s)',
            type: 'value',
          },
          {
            gridIndex: 1,
            name: 'output(B/s)',
            type: 'value',
            inverse: true
          },
        ],
        dataZoom: [
          {
            type: 'inside',
            xAxisIndex:[0,1],
            start: 0,
            end: 100
          },
          {
            xAxisIndex:[0,1],
            start: 0,
            end: 100
          }
        ],
        series: []
      }
    },
  },
  methods: {
    CreateChart(){
      this.netChart = this.$echarts.init(this.$refs.chartDom)
      setTimeout (()=>{
        this.$nextTick(() => {
          this.netChart.resize()
        })
      },0)
    },
    sizechange(params){
      this.netChart.resize(params) 
    },
    getAllData(rangeTime) {
      this.getStepData('in',this.netIn,'irate(node_network_receive_bytes_total[5m])',rangeTime);
      this.getStepData('out',this.netOut,'irate(node_network_transmit_bytes_total[5m])',rangeTime);
    },
    getStepData(type,thisData,querySql,rangeTime) {
      let params= {
        query: querySql,
        start: rangeTime.starttime,
        end: rangeTime.endtime,
        step: '10s'
      }
      getData(params).then(res => {
        if(res.data.status === 'success') {
          let legend = [];
          let index = 0;
          let currNetData = res.data.data.result.filter(item => item.metric.instance === this.$store.getters.selectIp);
          for(let i of currNetData) {
            index++;
            let resArr = []; 
            i.values.forEach(item => {
              resArr.push({
                time: item[0]*1000,
                value: [formatDate(item[0], 'yyyy-MM-dd hh:mm:ss'), parseInt(item[1]).toFixed(2)],
                name: i.metric.device
              })
            })
            legend.push(i.metric.device)
            switch (type) {
              case 'in':
                this.option.series[index] = {
                  name: i.metric.device,
                  smooth: true,
                  type: 'line',
                  showSymbol: false,
                  data: resArr
                }
                break;
              case 'out':
                this.option.series[index+4] = {
                  name: i.metric.device,
                  smooth: true,
                  type: 'line',
                  showSymbol: false,
                  xAxisIndex: 1,
                  yAxisIndex: 1,
                  data: resArr
                }
                break;
              default:
                break;
            }
            thisData.push(resArr);
          }
          this.option.legend.data = legend;
        }
      })
    },
  },
  watch: {
    netIn: function() {
      this.$nextTick(() => {
        this.netChart.setOption(this.option,true)
      })
    },
    netOut: function() {
      this.$nextTick(() => {
        this.netChart.setOption(this.option,true)
      })
    }
  }
}
</script>
<style scoped lang="scss">
  .panel {
    position: relative;
    background: transparent;
    height: 100%;
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
    #network {
      width: 100%;
      height: 100%;
    }
    .closeChart:hover {
      color: rgb(0, 163, 217)
    }
  }
  
</style>