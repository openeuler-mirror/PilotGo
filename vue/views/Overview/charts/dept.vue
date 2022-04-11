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
  Date: 2022-03-22 18:26:25
  LastEditTime: 2022-04-08 11:30:49
 -->
<template>
  <div id="dept" style="width:100%;height:100%">
    
  </div>
</template>
<script>
import { getDeptDatas } from '@/request/overview'
export default {
  name: 'DepartChart',
  props: {
    
  },
  data() {
    return {
      deptChart: {},
      xAxis: [],
      normal: [],
      offline: [],
      free: []
    }
  },
  computed: {
    option() {
      return {
        title: {
          text: '各部门机器数量'
        },
        tooltip: {
          show: true,
        },
        legend: {
          show: true,
          data: ["在线","离线","未分配"]
        },
        xAxis: {
          type: 'category',
          data: this.xAxis,
        },
        yAxis: {
          type: 'value'
        },
        grid: {
          left:50,
          right:40,
          bottom: 20
        },
        series: [
          {
            name: '在线',
            data: this.normal,
            type: 'bar',
            itemStyle: {
              color: 'rgb(11, 35, 117)',
              borderRadius: [9, 9, 0, 0]
            },
            label: {
              show: true,
              position: 'inside'
            },
          },
          {
            name: '离线',
            data: this.offline,
            type: 'bar',
            itemStyle: {
              color: 'rgb(202, 205, 210)',
              borderRadius: [9, 9, 0, 0]
            },
            label: {
              show: true,
              position: 'inside'
            },
          },
          {
            name: '未分配',
            data: this.free,
            type: 'bar',
            itemStyle: {
              color: 'rgb(92, 123, 207)',
              borderRadius: [9, 9, 0, 0]
            },
            label: {
              show: true,
              position: 'inside'
            },
          },
        ]
      };
    }
  },
  mounted() {
    this.deptChart = this.$echarts.init(document.getElementById('dept'))
    getDeptDatas().then(res => {
      if(res.data.code === 200) {
        let data = res.data.data.data;
        data.forEach(item => {
          this.xAxis.push(item.depart);
          this.normal.push(item.normal);
          this.offline.push(item.offline);
          this.free.push(item.free);
        })
      }
    })
  },
  methods: {
    resize(params) {
      this.deptChart.resize(params)
      this.deptChart.setOption(this.option,true)
    },
    getDept() {
      
  },
  
},
watch: {
    option: function() {
      this.deptChart.setOption(this.option,true)
    }
  }
}
</script>