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
  Date: 2022-02-22 16:43:19
  LastEditTime: 2022-06-29 10:28:09
  Description: 'Components Table'
 -->
<template>
  <div class="ky-table">
    <div class="header">
      <div class="header_content">
        <div class="table_search">
          <slot name="table_search"></slot>
        </div>
        <div class="table_action">
          <slot name="table_action"></slot>
        </div>
      </div>
    </div>
    <div class="content" v-loading="loading" element-loading-text="数据加载中" element-loading-spinner="el-icon-loading">
      <el-table height="100%" :header-cell-style="{ color: 'black', 'background-color': '#f6f8fd' }"
        :cell-style="{ color: 'black' }" ref="multipleTable" :data="tableData" tooltip-effect="dark"
        :row-class-name="tableRowClassName" @select="handleSelectionChange" @select-all="handleSelectAll">
        <el-table-column type="selection" v-if="showSelect" width="55" align="center" :selectable="checkSelectTable">
        </el-table-column>
        <slot name="table"></slot>
      </el-table>
    </div>
    <div class="pagination">
      <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="objSearch.page"
        :page-sizes="[10, 20, 25, 50, 75, 100]" :page-size="objSearch.size"
        layout="total, sizes, prev, pager, next, jumper" :total="total">
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    isLoadTable: {
      type: Boolean,
      default: function () {
        return true;
      }
    },
    isSource: {
      type: Number,
      default: function () {
        return 1;
      }
    },
    getData: {
      type: Function,
    },
    getSourceData: {
      type: Function,
    },
    showSelect: {
      type: Boolean,
      default: function () {
        return true;
      }
    },
    searchData: {
      type: Object,
      default: function () {
        return {

        };
      }
    },
    treeNodes: {
      type: Array,
      default: function () {
        return []
      }
    },
    isRowClick: {
      type: Boolean,
      default: function () {
        return false
      }
    }
  },
  data() {
    return {
      selectRow: {
        ids: [],
        rows: []
      },
      expands: [],
      checked: false,
      displayTip: false,
      total: 0,
      loading: false,
      tableData: [],
      display: false,
      objSearch: {
        page: 1,
        size: 10
      },
    };
  },
  mounted() {
    this.loadData({ ...this.objSearch });
  },
  watch: {
    isSource: function (newV, oldV) {
      this.$nextTick(() => {
        if (newV) {
          this.getSData({ ...this.objSearch, page: 1 })
        }
      })
    },
  },
  methods: {
    loadData(pageParams) {
      this.tableData = [];
      this.loading = true;
      this.getData({ ...pageParams, ...this.searchData }).then((response) => {
        const res = response.data;
        if (res.code === 200) {
          this.loading = false;
          this.total = res.total;
          this.objSearch.page = res.page;
          this.loading = false;
          this.tableData = res.data;
        } else {
          /* this.$message({
            type: "error",
            message: "数据格式错误",
          }); */
          this.loading = false;
        }
      }).catch((err) => {
        console.log("get data error:", err)
      })
    },
    getSData(data) {
      this.tableData = [];
      this.getSourceData(data).then((response) => {
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
    handleSetTableData(datas) {
      this.loading = false;
      this.tableData = datas;
    },
    handleSearch() {
      this.checked = false;
      this.$refs.multipleTable.clearSelection();
      this.selectRow.ids = [];
      this.selectRow.rows = [];
      this.loadData({ ...this.objSearch, page: 1 });
    },
    changeExpandRow(row) {
      this.$refs.multipleTable.toggleRowExpansion(row);
    },
    expandSelect(row, expandedRows) {
      this.expands = []
      if (expandedRows.length > 0) {
        row ? this.expands.push(row.name) : ''
      }
    },
    handleLoadSearch(data) {
      // 渲染高级搜索后的数据
      this.tableData = data;
    },
    handleSizeChange(size) {
      // 修改每页显示个数
      this.objSearch.size = size;
      this.loadData({ ...this.objSearch });
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
        this.selectRow.ids.push(row.ID || row.machineuuid || row.id);
        this.selectRow.rows.push(row);
      } else {
        this.selectRow.ids.splice(index, 1);
        this.selectRow.rows.splice(index, 1);
      }
    },
    handleSelectAll(selection) {
      if (selection.length === 0) {
        this.selectRow.ids = [];
        this.selectRow.rows = [];
      } else {
        selection.forEach(item => {
          this.selectRow.ids.push(item.ID);
          this.selectRow.rows.push(item);
        });
      }
    },
    refresh() {
      this.checked = false;
      this.$refs.multipleTable.clearSelection();
      this.selectRow.ids = [];
      this.selectRow.rows = [];
      this.loadData({ ...this.objSearch });
    },
    tableRowClassName({ row, rowIndex }) {
      // rowIndex备用，隔行显示颜色不同的时候用
      if (this.isRowClick) {
        return ['row-expand'];
      }

    },
    handleClose() {
      this.display = false;
    },
    checkSelectTable() {
      return this.treeNodes.length === 0;
    },
  },
};
</script>

<style rel="stylesheet/scss" lang="scss">
.ky-table {
  height: 96%;

  .header {
    width: 100%;
    height: 6%;
    border-radius: 6px 6px 0 0;
    background: linear-gradient(to right, rgb(11, 35, 117) 0%, rgb(96, 122, 207) 100%, );

    .header_content {
      height: 100%;
      margin: 0 10px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      color: #fff;

      .el-button {
        font-size: 12px;
        padding: 10px;
      }
    }
  }

  .content {
    height: 92%;
    overflow-y: auto;
  }

  .el-table {
    .line-color {
      background-color: #fff;
    }

    th,
    td {
      text-align: center;
    }

    .el-checkbox__input.is-checked .el-checkbox__inner,
    .el-checkbox__input.is-indeterminate .el-checkbox__inner {
      background-color: rgb(82, 108, 193);
      border-color: rgb(82, 108, 193)
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