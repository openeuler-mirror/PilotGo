<template>
  <div class="ky-table">
    <div class="content">
      <el-table
        cell-class-name="ky-cell-class"
        :header-cell-style="{ color: 'black', 'background-color': '#f6f8fd' }"
        :cell-style="{ color: 'black' }"
        v-loading="loading"
        ref="multipleTable"
        :data="tableData"
        tooltip-effect="dark"
        style="width: 100%"
        size="medium"
        :row-class-name="tableRowClassName"
        @select="handleSelectionChange"
      >
        <slot name="table"></slot>
      </el-table>
    </div>
    <div class="pagination">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="objSearch.page"
        :page-sizes="[20, 25, 50, 75, 100]"
        :page-size="objSearch.perPage"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    row_number: {
        type: Boolean,
        default: false,
    },
    row_id_name: {
      type: String,
      default: 'id',
    },
    rowClassName: {
      type: Function,
    },
    getData: {
      type: Function,
    },
    content: {
      type: String,
      default: "",
    },
    watchs: {
      type: Object,
      default: function () {
        return {};
      },
    },
  },
  data() {
    return {
      checked: false,
      select_id: "0",
      select_number: "0",
      displayTip: false,
      total: 0,
      loading: false, //待修改
      tableData: [],
      display: false,
      objSearch: {
        perPage: 20,
        page: 1,
        search: this.$route.query.search,
      },
    };
  },
  beforeDestroy() {
    /* if (this.timer) {
      clearInterval(this.timer);
    } */
  },
  mounted() {
    this.loadData({ ...this.objSearch, ...this.$route.query });
    for (let key in this.watchs) {
      this.$watch(key, this.watchs[key]);
    }
  },
  methods: {
    loadData(data = { page: 1, perPage: 20 }) {
      if (!this.modalTable) {
        this.$router.replace({
          query: {
            ...this.$route.query,
            page: data.page,
            perPage: data.perPage,
            search: data.search,
          },
        });
      }
      this.getData({ ...data }).then((response) => {
        const res = response.data;
        if (res.code == 0) {
          this.objSearch.perPage = res.data.perPage;
          this.total = res.data.total;
          this.objSearch.page = res.data.page;
          this.loading = false;
          this.tableData = res.data.results;
          this.objSearch.search = data.search;
        } else {
          this.$message({
            type: "error",
            message: res.message,
          });
        }
      });
    },
    handleSizeChange(size) {
      this.objSearch.perPage = size;
      this.loadData({ ...this.$route.query, ...this.objSearch });
    },
    handleCurrentChange(page) {
      this.objSearch.page = page;
      this.loadData(this.objSearch);
    },
    handleSelectionChange(selection, row) {
      if (this.select_id == "2") {
          const index = this.selectRow.excludeIds.indexOf(this.row_number ? row.rowNumber : row[this.row_id_name]);
          if (index < 0) {
          this.selectRow.excludeIds.push(this.row_number ? row.rowNumber : row[this.row_id_name]);
          this.selectRow.excludeRows.push(row);
          this.select_number = this.select_number - 1;
        } else {
          this.selectRow.excludeIds.splice(index, 1);
          this.selectRow.excludeRows.splice(index, 1);
          this.select_number = this.select_number + 1;
        }
        if (this.select_number) {
          this.checked = true;
        } else {
          this.checked = false;
        }
        return;
      }
      if (this.row_number) {
        const index = this.selectRow.rowNumbers.indexOf(row.rowNumber);
        if (index < 0) {
          this.selectRow.rowNumbers.push(row.rowNumber);
          this.selectRow.rows.push(row);
        } else {
          this.selectRow.rowNumbers.splice(index, 1);
          this.selectRow.rows.splice(index, 1);
        }
        this.select_number = this.selectRow.rowNumbers.length;
        if (this.selectRow.rowNumbers.length) {
          this.checked = true;
        } else {
          this.checked = false;
        }
      } else {
        const index = this.selectRow.ids.indexOf(row[this.row_id_name]);
        if (index < 0) {
          this.selectRow.ids.push(row[this.row_id_name]);
          this.selectRow.rows.push(row);
        } else {
          this.selectRow.ids.splice(index, 1);
          this.selectRow.rows.splice(index, 1);
        }
        this.select_number = this.selectRow.ids.length;
        if (this.selectRow.ids.length) {
          this.checked = true;
        } else {
          this.checked = false;
        }
      }
    },
    // 全选禁用
    /* selectable(row) {
      if (this.checked && this.select_id == "2") {
        return false;
      }
      return true;
    }, */
    handleSearch() {
      this.checked = false;
      this.$refs.multipleTable.clearSelection();
      this.selectRow.rowNumbers = [];
      this.selectRow.ids = [];
      this.selectRow.rows = [];
      this.loadData({...this.objSearch, page: 1});
    },
    refresh() {
      this.checked = false;
      this.$refs.multipleTable.clearSelection();
      this.selectRow.rowNumbers = [];
      this.selectRow.ids = [];
      this.selectRow.rows = [];
      this.loadData({ ...this.objSearch});
    },
    handleBlur() {
      this.displayTip = false;
    },
    handleFocus() {
      if (this.content) {
        this.displayTip = true;
      }
    },
    tableRowClassName({ row, rowIndex }) {
      let className = this.rowClassName ? this.rowClassName(row) : "";
      if (rowIndex % 2 == 1) {
        return "line-color " + className;
      }
      return className;
    },
    handleCheck() {
      this.checked = !this.checked;
      if (this.checked) {
        this.select_number = this.total;
        this.select_id = "2";
        this.selectRow.excludeRows = [];
        this.selectRow.excludeIds = [];
      } else {
        this.select_id = "0";
        this.select_number = "0";
      }
    },
    handleSelect(e) {
      this.select_id = e.target.id;
      if (this.select_id == "0") {
        this.checked = false;
      } else if (this.select_id == "1") {
        this.checked = true;
        this.tableData.forEach((row) => {
          this.$refs.multipleTable.toggleRowSelection(row, true);
          if (this.row_number) {
            const index = this.selectRow.rowNumbers.indexOf(row.rowNumber);
            if (index < 0) {
              this.selectRow.rowNumbers.push(row.rowNumber);
              this.selectRow.rows.push(row);
            }
          } else {
            const index = this.selectRow.ids.indexOf(row[this.row_id_name]);
            if (index < 0) {
              this.selectRow.ids.push(row[this.row_id_name]);
              this.selectRow.rows.push(row);
            }
          }
        });
        this.select_number = this.row_number
          ? this.selectRow.rowNumbers.length
          : this.selectRow.ids.length;
      } else {
        this.checked = true;
        this.select_number = this.total;
        this.selectRow.excludeIds = [];
        this.selectRow.excludeRows = [];
      }
    },
    handleClose() {
      this.display = false;
    },
    handleAdvanceSearch(search) {
      this.display = false;
      this.objSearch.search = search;
      this.handleSearch();
    },
    handleSelectRow() {
      if (this.checked) {
        if (this.select_id == "2") {
          this.tableData.forEach((row) => {
            if (!this.selectRow.excludeIds.includes(this.row_number ? row.rowNumber : row[this.row_id_name])) {
                this.$refs.multipleTable.toggleRowSelection(row, true);
            }
            if (this.row_number) {
              const index = this.selectRow.rowNumbers.indexOf(row.rowNumber);
              if (index < 0) {
                this.selectRow.rowNumbers.push(row.rowNumber);
                this.selectRow.rows.push(row);
              }
            } else {
              const index = this.selectRow.ids.indexOf(row[this.row_id_name]);
              if (index < 0) {
                this.selectRow.ids.push(row[this.row_id_name]);
                this.selectRow.rows.push(row);
              }
            }
          });
        } else {
          if (this.row_number) {
            this.tableData.forEach((row) => {
              if (this.selectRow.rowNumbers.indexOf(row.rowNumber ) > -1) {
                this.$refs.multipleTable.toggleRowSelection(row, true);
              }
            });
          } else {
            this.tableData.forEach((row) => {
              if (this.selectRow.ids.indexOf(row[this.row_id_name]) > -1) {
                this.$refs.multipleTable.toggleRowSelection(row, true);
              }
            });
          }
        }
      } else {
        this.$refs.multipleTable.clearSelection();
        this.selectRow.rowNumbers = [];
        this.selectRow.ids = [];
        this.selectRow.rows = [];
        this.select_id = "0";
        this.select_number = "0";
      }
    },
  },
  computed: {
    checkOrSelectId() {
      const { checked, select_id } = this;
      return { checked, select_id };
    },
  },
  watch: {
    tableData: function () {
      this.$nextTick(() => {
        this.handleSelectRow();
      });
    },
    checkOrSelectId: function () {
      this.handleSelectRow();
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