<template>
  <div class="overview">
    <div class="cluster">
      <span class="level-title">机器监控数据</span>
      <el-select v-model="macIp" filterable placeholder="请选择或输入关键字">
        <el-option
          v-for="item in options"
          :key="item.$index"
          :label="item"
          :value="item">
        </el-option>
      </el-select>
      <el-carousel trigger="click" type="card" :autoplay="false" ref="card" height="500px">
      <el-carousel-item v-for="item in chartData" 
        :key="item.index" 
        :name="item.lable" 
        :label="item.label"
        :id="item.chartId">
      </el-carousel-item>
    </el-carousel>
    </div>
    <div class="monitor">
      <span class="monitor-title">监控告警情况</span>
      <el-row class="monitor-row" :gutter="30">
        <el-col :span="12">
          <div class="table-title">
            <span class="table-title__left">机器异常</span>
            <div class="table-title__right">
              <span>共{{ machineTableData.length }}条</span>>
            </div>
          </div>

          <el-table :data="machineTable" style="width: 100%" size="small">
            <el-table-column prop="ip" label="IP"> </el-table-column>
            <el-table-column prop="type" label="类型"> </el-table-column>
            <el-table-column prop="details" label="详情"> </el-table-column>
          </el-table>
        </el-col>

        <el-col :span="12">
          <div class="table-title">
            <span class="table-title__left">监控异常</span>
            <div class="table-title__right">
              <span>共{{ monitorTableData.length }}条</span>
            </div>
          </div>

          <el-table :data="monitorTable" style="width: 100%" size="small">
            <el-table-column prop="ip" label="IP"></el-table-column>
            <el-table-column prop="type" label="类型"></el-table-column>
            <el-table-column prop="details" label="详情"></el-table-column>
          </el-table>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { getData, getCurrData, getPromeIp } from "@/request/overview";
export default {
  name: "Overview",
  props: {},
  data() {
    return {
      label: 'data',
      cpuChart: {},
      memChart: {},
      ioChart: {},
      netChart: {},
      cpuData: [],
      memData: [],
      netIn: [],
      netOut: [],
      inData: [],
      outData: [],
      netName: [],
      chartData: [
        {
          label: 'CPU',
          chartId: 'cpu',
          index: 1
        },
        {
          label: '内存',
          chartId: 'memory',
          index: 2
        },
        {
          label: 'I/O',
          chartId: 'io',
          index: 3
        },
        {
          label: '网络',
          chartId: 'network',
          index: 4
        }
      ],
      options: [],
      macIp: 'localhost:9100',
      chartW: 0,
      chartH: 0,
      now: new Date().getTime()/1000,
      currIndex: 0,
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
        // max: 100,
        min: 0,
        boundaryGap: [0, '100%'],
      },
      timer: [],
      pageSize: 4, //表格显示数据条数
      machineTable: [], //机器表格显示数据
      monitorTable: [], //监控表格显示数据 
      machineTableData: [
        {
          ip: "172.17.6.163",
          type: "虚拟机",
          details: "",
        },
      ],//表格总数据
      monitorTableData: [],
    };
  },
  computed: {
    ioData() {
      const {inData,outData} = this;
      return {inData,outData}
    },
    netData() {
      const {netIn,netOut} = this;
      return {netIn,netOut}
    },
    cpuOption() {
      return {
        title: {text:'cpu使用率'},
        tooltip: this.tooltip,
        xAxis: this.xAxis,
        yAxis: {...this.yAxis,max:100},
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
    memOption() {
      return {
        title: {text: '内存使用率'},
        tooltip: this.tooltip,
        xAxis: this.xAxis,
        yAxis: {...this.yAxis,max:100},
        series: [{
          name: 'memory',
          type: 'line',
          smooth: true,
          showSymbol: false,
          lineStyle: {width: 1},
          areaStyle: {
            opacity: 0.3,
            color:  '#37A2FF'
          },
          data: this.memData,
        }]
      }
    },
    ioOption() {
      return {
        title: {text: '磁盘读写速率'},
        tooltip: {
          trigger: 'axis',
        },
        legend: {
          data: []
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
        series: []
      }
    },
    netOption() {
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
    }


  },
  mounted() {
    getPromeIp({departid: this.$store.getters.UserDepartId}).then(res => {
      if(res.data.code == 200) {
        this.options = res.data.data;
      }
    })
    this.chartW = document.getElementsByClassName("cluster")[0].clientWidth/2+20;
    this.chartH = document.getElementsByClassName("cluster")[0].clientHeight;
    
    this.getStepCpu();
    this.getStepMem();
    this.getStepNet();
    this.getStepIO();
  },
  methods: {
    getStepData(thisData,itemIndex) {
      let params= {
        machineip: this.macIp,
        query: itemIndex,
        starttime: parseInt(this.now - 180) + '',
        endtime: parseInt(this.now - 0) + '',
      }
      getData(params).then(res => {
        if(res.data.code === 200) {
          if(itemIndex == 1 || itemIndex == 2) {
              res.data.data.forEach(item => {
              thisData.push({
                time: item.time,
                value: [ item.time,parseInt(item.value).toFixed(2)]
              })
            })
          } else {
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
                case 3:
                  this.ioOption.series[index]= {
                    neme: i.device,
                    smooth: true,
                    type: 'line',
                    showSymbol: false,
                    data: resArr,
                  }
                  break;
                case 4:
                  this.ioOption.series[index+4] = {
                    name: i.device,
                    smooth: true,
                    type: 'line',
                    showSymbol: false,
                    xAxisIndex: 1,
                    yAxisIndex: 1,
                    data: resArr
                  }
                  break;
                case 5:
                  this.netOption.series[index] = {
                    name: i.device,
                    smooth: true,
                    type: 'line',
                    showSymbol: false,
                    data: resArr
                  }
                  break;
                case 6:
                  this.netOption.series[index+4] = {
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
            if(itemIndex == 3 || itemIndex == 4) {
              this.ioOption.legend.data = legend;
            } else {
              this.netOption.legend.data = legend;
            }
          }
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
          if(itemIndex == 0 || itemIndex == 1) {
            thisData.shift();
            let item = res.data.data;
            thisData.push({
              time: item.time,
              value: [ item.time,parseInt(item.value).toFixed(2)]
            });
          } else {
            res.data.data.forEach( (i,index) => {
              thisData[index].shift();
              thisData[index].push({
                time: i.label.time,
                value: [i.label.time, parseInt(i.label.value).toFixed(2)],
                name: i.device
              })
            })
          }
        }
      })
    },
    getCurrCpu() {
      this.getPointData(this.cpuData,1);
    },
    getStepCpu() {
      this.cpuChart = this.$echarts.init(document.getElementById('cpu'))
      this.cpuChart.resize({width: this.chartW,height: this.chartH})
      this.getStepData(this.cpuData, 1);
      this.timer[0] = setInterval(this.getCurrCpu, 10000); 
    },
    getCurrMem() {
      this.getPointData(this.memData,2);
    },
    getStepMem() {
      this.memChart = this.$echarts.init(document.getElementById('memory'))
      this.memChart.resize({width: this.chartW,height: this.chartH})
      this.getStepData(this.memData,2);
      this.timer[1] = setInterval(this.getCurrMem, 10000); 
    },
    getCurrIO() {
      this.getPointData(this.inData,3);
      this.getPointData(this.outData,4);
    },
    getStepIO() {
      this.ioChart = this.$echarts.init(document.getElementById('io'))
      this.ioChart.resize({width: this.chartW,height: this.chartH})
      this.getStepData(this.inData,3);
      this.getStepData(this.outData,4);
      this.timer[2] = setInterval(this.getCurrIO, 10000); 
    },
    getCurrNet() {
      this.getPointData(this.netIn,5);
      this.getPointData(this.netOut,6);
    },
    getStepNet() {
      this.netChart = this.$echarts.init(document.getElementById('network'))
      this.netChart.resize({width: this.chartW,height: this.chartH})
      this.getStepData(this.netIn,5);
      this.getStepData(this.netOut,6);
      this.timer[3] = setInterval(this.getCurrNet, 10000); 
    },

  },
  watch: {
    cpuData: function(newData) {
        this.cpuChart.setOption(this.cpuOption,true)
    },
    memData: function(newData) {
      this.memChart.setOption(this.memOption,true)
    },
    netData: function(newData) {
      this.netChart.setOption(this.netOption,true)
    },
    ioData: function(newData) {
      this.ioChart.setOption(this.ioOption,true)
    },
    macIp: function(newIp) {
      this.macIp = newIp;
      clearInterval(this.timer[0]);
      clearInterval(this.timer[1]);
      clearInterval(this.timer[2]);
      clearInterval(this.timer[3]);
      this.getStepCpu();
      this.getStepMem();
    this.getStepNet();
    this.getStepIO();
      
    }
  },
};
</script>

<style scoped lang="scss">
.overview {
  margin: 0 10px;
  .cluster {
    .el-carousel__item {
      background-color: #d3dce6;
    }
  }

  //监控告警情况
  .monitor {
    .monitor-row {
      margin-top: 20px;
    }

    .table-title {
      width: 100%;
      height: 40px;
      border: 1px solid #f4f4f4;
      background-color: #fff;
      font-size: 14px;
      line-height: 40px;

      .table-title__left {
        float: left;
        margin-left: 6px;
      }
      .table-title__right {
        float: right;
        .el-pagination {
          float: right;
          padding: 0;
          height: 30px;
          margin-top: 5px;
        }
      }
    }
  }

  #echarts_box1,
  #echarts_box2 {
    width: 400px;
    height: 200px;
  }
}
</style>
