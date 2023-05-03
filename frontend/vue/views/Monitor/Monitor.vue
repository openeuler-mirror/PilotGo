<template>
  <div class="monitor">
    <div class="monitor-conditions">
      <span>监控参数选择</span>
      <div class="conditions-select">
        <el-row>
          <el-col :span="8">
            <span class="text">主机(可多选)</span>
            <el-select
              size="small"
              v-model="host"
              multiple>
              <el-option
                v-for="item in hostData"
                :key="item"
                :label="item"
                :value="item"
              ></el-option>
            </el-select>
          </el-col>
          <el-col :span="15">
            <span class="text">时间选择</span>
            <el-date-picker
              v-model="dataTime"
              type="datetimerange"
              :picker-options="pickeroptions"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small">
            </el-date-picker>
          </el-col>

        </el-row>

        <el-row>
          <el-col :span="16">
            <span class="text">指标</span>
            <el-select size="small" v-model="indicators" style="width: 550px">
              <el-option
                v-for="item in inputData"
                :key="item"
                :label="item"
                :value="item"
              ></el-option>
            </el-select>
          </el-col>

          <el-col :span="4">
            <el-button
              type="primary"
              size="mini"
              :disabled="!(indicators.length != 0 && dataTime.length != 0 && host.length != 0)"
              @click="handleClick"
              >确定</el-button>
            <el-button type="primary" size="mini" @click="handleClickReset">重置</el-button>
          </el-col>
        </el-row>
      </div>
    </div>

    <div class="monitor-warnning">
      <span>监控告警图</span>
      <div id="monitor-echarts_box"></div>
    </div>

  </div>
</template>

<script>
import { getChart, getChartName } from "@/request/api";

export default {
  name: "Monitor",
  data() {
    return {
      host: [], // 主机
      hostData: ['localhost:9080', 'localhost:9100'],
      inputData: [],  //select的option
      indicators: "", //选中的指标
      dataTime: "", // 选择的时间
      pickeroptions: {
        disabledDate : time => {
          let myDate = new Date()
          return time.getTime() > myDate.getTime()// 返回 所有时间 大于 当前时间
        }
      },

      myChart: "",
      chartOption: {
        //图表数据
        title: { text: "" },
        tooltip: {},
        grid: {
          left: "5%",
          bottom: "30%"
        },
        legend: {
          data: [],
          orient: 'vertical',
          x: "left",
          y: "bottom",
          left: "1%",
          backgroundColor: "#222",
          textStyle: {
            color: "#fff", // 图例文字颜色
            fontSize: 10
          }
        },
        xAxis: {
          data: []
        },
        yAxis: { type: "value" },
        series: [] // {name, type, areaStyle{data[]}}
      }
    };
  },
  methods: {
    handleClickReset() {
      let _this = this
      _this.indicators = ''
      _this.dataTime = ''
      _this.host = ''
    },
    handleClick() {
      let _this = this;

      if (_this.indicators.length != 0 && _this.host.length != 0 && _this.dataTime != 0) {
        let count = 0; //计数器
        _this.chartOption.series = [];
        _this.chartOption.legend.data = [];
        _this.chartOption.xAxis.data = [];

        let systime = new Date();
        let startTime = _this.dataTime[0].getTime() / 1000
        let endTime = _this.dataTime[1].getTime() / 1000
        let hours = endTime - startTime == 0 ? 1 : Math.ceil(( endTime - startTime ) / 3600 )
        // 当开始时间和结束时间相等时， 将开始时间后退一分钟，避免url出错获取不到数据
        if(endTime - startTime == 0) {
          startTime = startTime - 60;
        }
        let step = hours * 14  //向上取整获得小时数*每小时14step

        let url = `/query_range?query=${_this.indicators}&start=${startTime}&end=${endTime}&step=${step}&_=${systime}`

        getChart(url).then(res => {
          let chartData = res.data.data.result;
          console.log('res', chartData)

          chartData.filter(item => {
            // 筛选出符合主机的数据
            if (_this.host.indexOf(item.metric.instance) > -1) {
              return item
            }
          }).forEach(item => {

            let name = _this.indicators + "{";
            for (const metricKey in item.metric) {
              let oldKey = "";
              if (oldKey !== metricKey && metricKey !== "__name__") {
                name += `${metricKey}="${item.metric[metricKey]}",`
              }
              oldKey = metricKey;
            }
            name = name.substring(0, name.length - 1); //去除最后一个逗号
            name += "}";

            //图例组件  系列

            _this.chartOption.legend.data.push(name);
            _this.chartOption.series.push({
              name: name,
              type: "line",
              data: []
            });

            //图表 系列列表
            item.values.forEach(value => {
              if (count == 0) {
                // 只循环一次
                let xAxisData = new Date(value[0] * 1000);
                let xYear = xAxisData.getFullYear(),
                    xMonth = xAxisData.getMonth(),
                    xDate = xAxisData.getDate(),
                    xHours = xAxisData.getHours(),
                    xMin = xAxisData.getMinutes(),
                    xSeconds = xAxisData.getSeconds();


                // 根据时间的不同x轴的刻度名称不同
                if((hours / 24) - 365 >= 0) {
                  xAxisData =`${xYear}年${xMonth}月${xDate}日`
                } else if((hours / 24) >= 28) {
                  xAxisData =`${xMonth}月${xDate}日`
                } else if (hours > 24) {
                  xAxisData =`${xDate}号${xHours}点`
                } else if (hours <= 24 && hours > 1) {
                  xAxisData =`${xHours}:${xMin}`
                } else {
                  xAxisData =`${xMin}m${xSeconds}s`
                }

                _this.chartOption.xAxis.data.push(xAxisData);
              }
              _this.chartOption.series[count].data.push(Number(value[1]));
            });

            count++;
          });
          // console.log('chartOption', _this.chartOption)

          // 获取dom
          let chartDiv = document.getElementById("monitor-echarts_box");

          //基于准备好的dom，初始化echarts实例
          if (!_this.myChart) {
            _this.myChart = this.$echarts.init(chartDiv);
          }

          // 使用刚指定的配置项和数据，显示图表
          _this.myChart.setOption(_this.chartOption, true);
        });
      }
    },
  },
  mounted() {
    let _this = this;
    let time = Date.now(); //获取当前时间 转换为毫秒

    getChartName(time).then(res => {
      _this.inputData = res.data.data;
    });
  },
};
</script>

<style lang="scss" scope>
.monitor {
  .monitor-conditions {
    .conditions-select {
      background: #fff;
      margin-top: 20px;
      height: 100px;
      line-height: 50px;
      // text-align: center;
      .el-row {
        margin-left: 30px;
      }

      .el-col {
        .text {
          font-size: 14px;
        }
      }
    }
  }

  .monitor-warnning {
    margin-top: 30px;
    #monitor-echarts_box {
      margin-top: 20px;
      background: #fff;
      width: 100%;
      height: 500px;
    }
  }
}
</style>
