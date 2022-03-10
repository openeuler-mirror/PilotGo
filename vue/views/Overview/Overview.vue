<template>
  <div class="overview">
    <div class="cluster">
      <span class="level-title">机器监控数据</span>
      <el-select v-model="macIp" filterable placeholder="请选择或输入关键字">
        <el-option
          v-for="item in options"
          :key="item.label"
          :label="item.label"
          :value="item.label">
        </el-option>
      </el-select>
      <el-carousel trigger="click" type="card" :autoplay="false" ref="card" height="360px">
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
import { getData, getCurrData } from "@/request/overview";
import { getClusters } from "@/request/cluster";
export default {
  name: "Overview",
  props: {},
  data() {
    return {
      label: 'data',
      cpuChart: {},
      memChart: {},
      cpuData: [],
      memData: [],
      netData: [],
      inData: [],
      outData: [],
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
      options: [{
          value: '0',
          label: 'localhost:9100'
        },{
          value: '1',
          label: '172.17.127.20'
        }, {
          value: '2',
          label: '172.17.127.21'
        }, {
          value: '3',
          label: '172.17.127.22'
        }, ],
      macIp: 'localhost:9100',
      chartW: 0,
      chartH: 0,
      now: new Date().getTime()/1000,
      currIndex: 0,
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
    cpuOption() {
      return {
        title: {
          text: 'cpu使用率'
        },
        tooltip: {
          trigger: 'axis',
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
          min: 0,
          boundaryGap: [0, '100%'],
        },
        series: [
          {
            name: '',
            type: 'line',
            smooth: true,
            showSymbol: false,
            lineStyle: {
              width: 1
            },
            areaStyle: {
              opacity: 0.3,
              color:  '#37A2FF'
            },
            data: this.cpuData,
          }
        ]
      }
    },
    memOption() {
      return {
        title: {
          text: '内存使用率'
        },
        tooltip: {
          trigger: 'axis',
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
          min: 0,
          boundaryGap: [0, '100%'],
        },
        series: [
          {
            name: '',
            type: 'line',
            smooth: true,
            showSymbol: false,
            lineStyle: {
              width: 1
            },
            areaStyle: {
              opacity: 0.3,
              color:  '#37A2FF'
            },
            data: this.memData,
          }
        ]
      }
    }
  },
  mounted() {
    this.chartW = document.getElementsByClassName("cluster")[0].clientWidth/2+20;
    this.chartH = document.getElementsByClassName("cluster")[0].clientHeight;
    
    this.getStepCpu();
    this.getStepMem();
    // this.getStepNet();
    // this.getStepIO();
  },
  methods: {
    getStepData(thisData,itemIndex) {
      let params= {
        machineip: this.macIp,
        query: itemIndex,
        starttime: parseInt(this.now - 180) + '',
        endtime: parseInt(this.now - 0) + '',
        step: 10
      }
      getData(params).then(res => {
        if(res.data.code === 200) {
          res.data.data.forEach(item => {
            thisData.push({
              time: item.time,
              value: [ item.time,item.value]
            })
          })
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
          thisData.shift();
          let item = res.data.data;
          thisData.push({
            time: item.time,
            value: [ item.time,item.value]
          });
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
      setInterval(this.getCurrCpu, 10000); 
    },
    getCurrMem() {
      this.getPointData(this.memData,2);
    },
    getStepMem() {
      this.memChart = this.$echarts.init(document.getElementById('memory'))
      this.memChart.resize({width: this.chartW,height: this.chartH})
      this.getStepData(this.memData,2);
      setInterval(this.getCurrCpu, 10000); 
    },
    getCurrNet() {
      this.getPointData(this.netData,4);
    },
    getStepNet() {
      this.getStepData(this.netData,4);
      this.option.title.text = "平均网络输入输出";
      setInterval(this.getCurrCpu, 10000); 
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
      this.resizeChart('network',{width: this.chartW,height: this.chartH})
      this.myChart.setOption(this.option,true)
    },
    inData: function(newData) {
      this.resizeChart('io',{width: this.chartW,height: this.chartH})
      this.myChart.setOption(this.option,true)
    },
    macIp: function(newIp) {
      clearInterval();
    }
  },
  beforeRouteLeave(to, form, next) {
    clearInterval();
    next()
  }
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
    margin-top: 40px;

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
