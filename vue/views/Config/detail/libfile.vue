<template>
  <div v-loading="loading" style="height:100%">
    <ky-table
      :getData="libFileList"
      ref="table"
    >
      <template v-slot:table_search>
        <el-input placeholder="请输入关键字进行搜索..." prefix-icon="el-icon-search"
                  clearable
                  style="width: 280px;margin-right: 10px;" v-model="searchInput"
                  @keydown.enter.native="handleSearch"></el-input>
        <el-button icon="el-icon-search" @click="handleSearch">
          搜索
        </el-button>
      </template>
      <template v-slot:table_action>
        <el-popconfirm 
          title="确定删除此文件？"
          cancel-button-type="default"
          confirm-button-type="danger"
          @confirm="handleDelete">
          <auth-button name="user_del" slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 删除 </auth-button>
        </el-popconfirm>
      </template>
      <template v-slot:table>
        <el-table-column  prop="id" label="编号" sortable width="80">
        </el-table-column>
        <el-table-column  prop="path" label="路径">
        </el-table-column>
        <el-table-column  prop="name" label="文件名">
        </el-table-column>
        <el-table-column prop="CreatedAt" label="创建时间" sortable>
          <template slot-scope="scope">
            <span>{{scope.row.CreatedAt | dateFormat}}</span>
          </template>
        </el-table-column>
        <el-table-column  prop="description" label="描述">
        </el-table-column>
      </template>
    </ky-table>
    
    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="display" 
      width="560px"
    >
     <update-form :row="rowData" v-if="type === 'update'" @click="handleClose"></update-form>
    </el-dialog>

  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import { libFileList, delLibFile,libFileSearch } from "@/request/config"
export default {
  components: {
    kyTable,
    AuthButton,
  },
  data() {
    return {
      loading: false,
      display: false,
      searchInput: '',
      title: "",
      type: "",
      rowData: {},
      tableData: []
    }
  },
  methods: {
    libFileList,
    handleClose(type) {
      this.display = false;
      this.title = "";
      this.type = "";
      if(type === 'success') {
        this.refresh();
      }
    }, 
    refresh(){
      this.$refs.table.handleSearch();
    },
    handleSearch() {
      libFileSearch({'search': this.searchInput}).then((res) => {
        if(res.data.code === 200) {
          this.$refs.table.handleLoadSearch(res.data.data);
        }
      })
    },
    handleDelete() {
      delLibFile({ids: this.$refs.table.selectRow.ids}).then(res => {
        if(res.status === 200) {
          this.$message.success(res.data.msg);
          this.refresh();
        } else {
          this.$message.error(res.data.msg);
        }
      })

    },
  },
  filters: {
    dateFormat: function(value) {
      let date = new Date(value);
      let y = date.getFullYear();
      let MM = date.getMonth() + 1;
      MM = MM < 10 ? "0" + MM : MM;
      let d = date.getDate();
      d = d < 10 ? "0" + d : d;
      let h = date.getHours();
      h = h < 10 ? "0" + h : h;
      let m = date.getMinutes();
      m = m < 10 ? "0" + m : m;
      let s = date.getSeconds();
      s = s < 10 ? "0" + s : s;
      return y + "-" + MM + "-" + d + " " + h + ":" + m;
    }
  },
}
</script>

<style scoped>
.search-form{
  /* margin-bottom: 12px; */
}
.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
</style>
