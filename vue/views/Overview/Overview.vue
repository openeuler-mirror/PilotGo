<template>
  <div class="overview">
    <div class="cluster">
      <span class="level-title">集群整体情况</span>
      <el-row class="cluster-row">
        <el-col :span="12">
          <div id="echarts_box1"></div>
        </el-col>
        <el-col :span="12">
          <div id="echarts_box2"></div>
        </el-col>
      </el-row>
    </div>
    <div class="monitor">
      <span class="monitor-title">监控告警情况</span>
      <el-row class="monitor-row" :gutter="30">
        <el-col :span="12">
          <div class="table-title">
            <span class="table-title__left">机器异常</span>
            <div class="table-title__right">
              <span>共{{ machineTableData.length }}条</span>
              <el-pagination
                layout="prev, next"
                :total="machineTableData.length"
                :page-size="pageSize"
                @current-change="handleCurrentChangeMachineTable"
              >
              </el-pagination>
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
              <el-pagination
                layout="prev, next"
                :total="monitorTableData.length"
                :page-size="pageSize"
                @current-change="handleCurrentChangeMonitorTable"
              >
              </el-pagination>
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
import { getOverview } from "@/request/api";
export default {
  name: "Overview",
  props: {},
  data() {
    return {
      char1: "",
      char2: "",
      pageSize: 4, //表格显示数据条数
      machineTable: [], //机器表格显示数据
      monitorTable: [], //监控表格显示数据
      //表格总数据
      machineTableData: [
        {
          ip: "172.17.6.163",
          type: "虚拟机",
          details: "",
        },
      ],
      monitorTableData: [],
      chartData: [],
    };
  },
  computed: {
    option() {
      let _this = this;
      let option = {
        title: {
          text: "",
          left: "center",
          top: 20,
          textStyle: {
            color: "#ccc",
          },
        },
        tooltip: {
          trigger: "item",
        },
        legend: {
          orient: "vertical",
          right: 0,
          data: ["虚拟机", "物理机"],
        },
        series: [
          {
            name: "集群情况",
            type: "pie",
            radius: ["50%", "70%"],
            data: _this.chartData,
          },
        ],
      };
      return option;
    },
  },
  methods: {
    refreshData() {
      let _this = this;
      getOverview().then((res) => {
        let chartData = res.data.data;
        for (let i of chartData) {
          if (i.name == "virtual") {
            i.name = "虚拟机";
          }
          if (i.name == "physics") {
            i.name = "物理机";
          }
        }
        _this.chartData = chartData;
        let chartDiv = document.getElementById("echarts_box1");
        //基于准备好的dom，初始化echarts实例
        _this.chart1 = this.$echarts.init(chartDiv);
        _this.chart2 = this.$echarts.init(
          document.getElementById("echarts_box2")
        );
        _this.chart1.setOption(_this.option);
        _this.chart2.setOption(_this.option);

        _this.machineTable = _this.machineTableData.slice(0, _this.pageSize);
        _this.monitorTable = _this.monitorTableData.slice(0, _this.pageSize);
      });
    },
    currentPageChange(currentPage, tableData, tableName) {
      let _this = this;
      let data = tableData.slice(
        (currentPage - 1) * _this.pageSize,
        currentPage * _this.pageSize
      );
      if (tableName == "machine") {
        _this.machineTable = data;
      } else if (tableName == "monitor") {
        _this.monitorTable = data;
      }
    },
    //机器异常数据分页处理
    handleCurrentChangeMachineTable(currentPage) {
      let _this = this;
      this.currentPageChange(currentPage, _this.machineTableData, "machine");
    },
    //监控异常数据分页处理
    handleCurrentChangeMonitorTable(currentPage) {
      let _this = this;
      this.currentPageChange(currentPage, _this.monitorTableData, "monitor");
    },
  },

  mounted() {
    let _this = this;
    _this.refreshData();
  },
};
</script>

<style scoped lang="scss">
.overview {
  margin: 0 10px;
  .cluster {
    .cluster-row {
      margin-top: 20px;
      padding-top: 10px;
      background-color: #fff;
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
