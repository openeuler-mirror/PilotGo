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
  LastEditTime: 2022-06-23 11:13:21
 -->
<template>
  <div id="io" ref="chartDom"></div>
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
      diskChart: {},
      inData: [],
      outData: [],
      now: new Date()/1000,
    }
  },
  mounted() {
    if(this.$store.getters.selectIp){
      this.getStepData('write',this.inData,'irate(node_disk_writes_completed_total[1m])',{starttime: parseInt(this.now - 60*60*2) + '',
          endtime: parseInt(this.now - 0) + '',});
      this.getStepData('read',this.outData,'irate(node_disk_reads_completed_total[1m])',{starttime: parseInt(this.now - 60*60*2) + '',
          endtime: parseInt(this.now - 0) + '',});
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
          data: []
        },
        grid: [
          {
            bottom: '60%',
            left: '5%',
            top: '4%',
            right: '3%'
          },
          {
            top: '46%',
            left: '5%',
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
            name: 'input(K/s)',
            type: 'value',
          },
          {
            gridIndex: 1,
            name: 'output(K/s)',
            type: 'value',
            inverse: true
          },
        ],
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 100,
            xAxisIndex:[0,1],
          },
          {
            start: 0,
            end: 100,
            xAxisIndex:[0,1],
          }
        ],
        series: []
      }
    },
  },
  methods: {
    CreateChart(){
      this.diskChart = this.$echarts.init(this.$refs.chartDom)
      setTimeout (()=>{
        this.$nextTick(() => {
          this.diskChart.resize()
        })
      },0)
    },
    sizechange(params){
      this.diskChart.resize(params) 
    },
    getAllData(rangeTime) {
      this.getStepData('write',this.inData,'irate(node_disk_writes_completed_total[1m])',rangeTime);
      this.getStepData('read',this.outData,'irate(node_disk_reads_completed_total[1m])',rangeTime);
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
          let currDiskData = res.data.data.result.filter(item => item.metric.instance === this.$store.getters.selectIp);
          for(let i of currDiskData) {
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
              case 'write':
                this.option.series[index]= {
                  neme: i.metric.device,
                  smooth: true,
                  type: 'line',
                  showSymbol: false,
                  data: resArr,
                }
                break;
              case 'read':
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
    }
  },
  watch: {
    inData: function() {
      this.$nextTick(() => {
        this.diskChart.setOption(this.option,true)
      })
    },
    outData: function() {
      this.$nextTick(() => {
        this.diskChart.setOption(this.option,true)
      })
    }
  }
}
</script>
<style scoped lang="scss">
  #io {
    width: 100%;
    height: 100%;
    color: rgb(38, 51, 173)
  }
  
</style>