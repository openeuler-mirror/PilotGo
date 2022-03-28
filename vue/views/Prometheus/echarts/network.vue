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
  LastEditTime: 2022-03-25 16:25:04
 -->
<template>
  <div class="panel">
    <div id="network"></div>
  </div>
</template>
<script>
import { getData, getCurrData } from "@/request/overview";
export default {
  props: {
    macIp: {
      type: String,
    }
  },
  data() {
    return {
      netChart: {},
      netIn: [],
      netOut: [],
      now: new Date()/1000,
    }
  },
  mounted() {
    this.netChart = this.$echarts.init(document.getElementById('network'))
    this.getStepData(this.netIn,5);
    this.getStepData(this.netOut,6);
  },
  computed: {
    option() {
      return {
        title: {text: '网络平均速率'},
        tooltip: {
          trigger: 'axis',
        },
        legend: {
          data: [],
        },
        grid: [
          {
            bottom: '60%'
          },
          {
            top: '46%'
          }
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
        series: []
      }
    },
  },
  methods: {
    resize(params) {
      this.netChart.resize(params)
      this.netChart.setOption(this.option,true)
    },
    getStepData(thisData,itemIndex) {
      let params= {
        machineip: this.macIp,
        query: itemIndex,
        starttime: parseInt(new Date().getTime()/1000 - 180) + '',
        endtime: parseInt(new Date().getTime()/1000 - 0) + '',
      }
      getData(params).then(res => {
        if(res.data.code === 200) {
          let legend = [];
          let index = 0;
          for(let i of res.data.data) {
            index++;
            let resArr = []; 
            i.label.forEach(item => {
              resArr.push({
                time: item.time,
                value: [item.time, parseInt(item.value).toFixed(2)],
                name: i.device
              })
            })
            legend.push(i.device)
            switch (itemIndex) {
              case 5:
                this.option.series[index] = {
                  name: i.device,
                  smooth: true,
                  type: 'line',
                  showSymbol: false,
                  data: resArr
                }
                break;
              case 6:
                this.option.series[index+4] = {
                  name: i.device,
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
    getPointData(thisData,itemIndex) {
      let params= {
        machineip: this.macIp,
        query: itemIndex,
        time: parseInt(new Date().getTime()/1000 - 0) + '',
      }
      getCurrData(params).then(res => {
        if(res.data.code === 200) {
          res.data.data.forEach( (i,index) => {
            thisData[index].shift();
            thisData[index].push({
              time: i.label.time,
              value: [i.label.time, parseInt(i.label.value).toFixed(2)],
              name: i.device
            })
          })
        }
      })
    },
  }
}
</script>