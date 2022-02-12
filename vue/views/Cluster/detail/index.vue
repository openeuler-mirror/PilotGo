<template>
 <div>
   <div id="diskChar" :style="{width:'480px',height:'460px'}"></div>
 </div>
</template>
<script>
import { getDeviceInfo } from '@/request/cluster'
export default {
   /* 某一台机器的详情，左右分布
      左：柱状图，磁盘信息，占用多少空闲多少
      右：暂定
   */
  props: {
    ip: {
      type: String
    } 
  },
  data() {
    return {
    }
  },
  mounted() {
    this.drawEchart();
  },
  methods: {
    drawEchart() {
      let myChart = this.$echarts.init(document.getElementById('diskChar'))
      let option = {
        title: {
            text: '',
            top: 0,
            left: 10
        },
        tooltip: {
            trigger: 'item',
            formatter: "{a} <br/>{b} : {c}%"
        },
        toolbox: {
            show : true,
            top:0,
            right:0,
            feature : {
                mark : {show: false},
                magicType : {show: false, type: ['line', 'bar']},
                saveAsImage : {show: true}
            }
        },
        grid: {
          top: 60,
          right: 30,
          bottom: 30,
          left: 60
        },
        legend: {
          top: 0,
          left: 'center',
          data: ['占用','空闲']
        },
        calculable: true,
        xAxis: [
          {
            name: '磁盘',
            nameGap: 6,
            type: 'category',
            data: ['C','D','E','F'],
            axisLabel:{					//---坐标轴 标签
              show:true,					//---是否显示
              inside:false,				//---是否朝内
              rotate:0,					//---旋转角度	
              margin: 10,					//---刻度标签与轴线之间的距离
              color:'#000',
              fontSize: '16'				//---默认取轴线的颜色
            },
          }
        ],
        yAxis : [
          {
            type: 'value',
            name: '比例%',
            axisLine: {
              show: true,
            },
            axisLabel: {
              show:true,
              showMinLabel:true,
              showMaxLabel:true,
              formatter: function (value) {
                  return value;
              }
            }
          },
          {
            type: 'value',
            name: "",
            axisLabel: {
                show:false,
            }
          }
        ],
        series: [
            {
              name:'占用',
              type:'bar',
              yAxisIndex: 0,
              data: [20,30,10,78],
            },
            {
              name:'空闲',
              type:'bar',
              yAxisIndex: 1,
              data:[18,78,28,65],
            }
        ]
    };
      // 通过ip获取磁盘信息
      /* getDeviceInfo({'ip': this.ip}).then((res) => {
        if(res.code == 200) {
          // 设置option中数据的信息
          option.series[0].data = res.data.data[0];
          option.series[1].data = res.data.data[1];
        }
        myChart.setOption(option);
      }) */
      myChart.setOption(option);
    }
  }
}
</script>
<style scoped>
</style>
