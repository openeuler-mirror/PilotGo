<template>
  <div class="ky-table">
    <div class="content">
      <el-table
        :header-cell-style="{ color: 'black', 'background-color': '#f6f8fd' }"
        :cell-style="{ color: 'black' }"
        v-loading="loading"
        ref="multipleTable"
        :data="tableData"
        tooltip-effect="dark"
        :row-class-name="tableRowClassName"
        @select="handleSelectionChange"
        @select-all="handleSelectAll"
      >
        <el-table-column
          type="selection"
          width="55"
          align="center"
        >
        </el-table-column>
        <slot name="table"></slot>
      </el-table>
    </div>
    <div class="pagination">
      <el-pagination
        @current-change="handleCurrentChange"
        :current-page="objSearch.page"
        layout="total, prev, pager, next, jumper"
        :total="total"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    rowClassName: {
      type: Function,
    },
    getData: {
      type: Function,
    },
    searchData: {
      type: Object,
      default: function() {
        return {};
      }
    }
  },
  data() {
    return {
      selectRow: {
        ids: [],
        rows: []
      },
      checked: false,
      displayTip: false,
      total: 0,
      loading: false, 
      tableData: [],
      display: false,
      objSearch: {
        page: 1,
      },
    };
  },
  mounted() {
    this.loadData({ ...this.objSearch, ...this.$route.query });
  },
  methods: {
    loadData(data) {
      this.loading = true;
      this.getData({ ...data, ...this.searchData }).then((response) => {
        const res = response.data;
        if (res.code === 200) {
          this.total = res.total;
          this.objSearch.page = res.page;
          this.loading = false;
          this.tableData = res.data;
        } else {
          this.$message({
            type: "error",
            message: "数据格式错误",
          });
        }
      });
    },
    handleCurrentChange(page) {
      this.objSearch.page = page;
      this.loadData(this.objSearch);
    },
    handleSelectionChange(selection, row) {
        // 1.判断是否已经存在已选中行
        const index = this.selectRow.ids.indexOf(row.ID || row.machineuuid);
        // 2.若是不在就放入选中数组中，在就删掉
        if (index < 0) {
          this.selectRow.ids.push(row.ID || row.machineuuid);
          this.selectRow.rows.push(row);
        } else {
          this.selectRow.ids.splice(index, 1);
          this.selectRow.rows.splice(index, 1);
        }
        console.log(this.selectRow.ids)
    },
    handleSelectAll(selection) {
      if(selection.length === 0) {
        this.selectRow.ids = [];
        this.selectRow.rows = [];
      } else {
        selection.forEach(item => {
          this.selectRow.ids.push(item.ID);
          this.selectRow.rows.push(item);
        });
      }
    },
    handleSearch() {
      this.checked = false;
      this.$refs.multipleTable.clearSelection();
      this.selectRow.ids = [];
      this.selectRow.rows = [];
      this.loadData({...this.objSearch, page: 1});
    },
    refresh() {
      this.checked = false;
      this.$refs.multipleTable.clearSelection();
      this.selectRow.ids = [];
      this.selectRow.rows = [];
      this.loadData({ ...this.objSearch});
    },
    tableRowClassName({ row, rowIndex }) {
      let className = this.rowClassName ? this.rowClassName(row) : "";
      if (rowIndex % 2 == 1) {
        return "line-color " + className;
      }
      return className;
    },
    handleClose() {
      this.display = false;
    },
  },
};
</script>

<style rel="stylesheet/scss" lang="scss">
.ky-table {
  .el-table {
    border: 1px solid #ebeef5;
    .line-color {
      background-color: #f2f7ff;
    }
    th, td {
      text-align: center;
    }
  }
  .pagination {
    .el-pagination {
      text-align: right;
      /* border: 1px solid #ccc; */
      border-top: 0px;
      padding-top: 5px;
      padding-bottom: 5px;
      border-top: 0px;
      .el-pagination__sizes,
      .el-pagination__total {
        float: left;
      }
    }
  }
}
</style>