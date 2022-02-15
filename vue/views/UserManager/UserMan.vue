<template>
  <div v-loading="loading">
    <ky-table
      :getData="getUsers"
      ref="table"
      id="exportTab"
    >
      <template v-slot:table_search>
        <el-input placeholder="请输入手机号或邮箱名进行搜索..." prefix-icon="el-icon-search"
                  clearable
                  disabled
                  style="width: 280px;margin-right: 10px;" v-model="keyword"
                  @keydown.enter.native="searchUser"></el-input>
        <el-button icon="el-icon-search" disabled @click="searchUser">
          搜索
        </el-button>
      </template>
      <template v-slot:table_action>
        <el-button @click="handleCreate"> 添加 </el-button>
        <el-popconfirm 
          title="确定删除此用户？"
          cancel-button-type="default"
          confirm-button-type="danger"
          @confirm="handleDelete">
          <el-button slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 删除 </el-button>
        </el-popconfirm>
        <el-button @click="handleExport"> 导出 </el-button>
        <el-upload
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="onSuccess"
          :on-error="onError"
          name="upload"
          accept="xlsx"
          style="display: inline-flex;margin-right: 8px"
          action="/user/import">
          <el-button> 批量导入 </el-button>
        </el-upload>
      </template>
      <template v-slot:table>
        <el-table-column prop="ID" label="编号" width="60">
        </el-table-column>
        <el-table-column  prop="username" label="用户名" width="160">
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="120">
        </el-table-column>
        <el-table-column  prop="email" label="邮箱">
        </el-table-column>
        <el-table-column prop="enable" label="启用" width="80">
          <template slot-scope="scope">
            {{scope.row.enable === true ? "是" : "否"}}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template slot-scope="scope">
            <el-button class="editBtn" type="primary" plain size="mini" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button class="editBtn" type="primary" plain size="mini" @click="handleReset(scope.row.email)">重置密码</el-button>
          </template>
        </el-table-column>
      </template>
    </ky-table>
    
    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="display" 
      width="560px"
    >
     <add-form v-if="type === 'create'" @click="handleClose"></add-form>
     <update-form :row="rowData" v-if="type === 'update'" @click="handleClose"></update-form>
    </el-dialog>

  </div>
</template>

<script>
import FileSaver from 'file-saver'
import XLSX from 'xlsx'
import AddForm from "./form/addForm.vue"
import UpdateForm from "./form/updateForm.vue"
import kyTable from "@/components/KyTable";
import { getUsers, delUser, resetPwd, } from "@/request/user"
export default {
  components: {
    kyTable,
    AddForm,
    UpdateForm,
  },
  data() {
    return {
      loading: false,
      display: false,
      isDelete: true,
      title: "",
      type: "",
      keyword: '',
      rowData: {},
      tableData: []
    }
  },
  methods: {
    getUsers,
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
    searchUser() {
      console.log("待写入按关键子查找用户");
    },
    handleCreate() {
      this.display = true;
      this.title = "新增用户";
      this.type = "create";
    },
    handleEdit(row) {
      this.rowData = row;
      this.display = true;
      this.title = "编辑用户";
      this.type = "update";
    },
    handleReset(email) {
      resetPwd({email: email}).then(res => {
        if(res.code == 200){
          this.$message.success("重置密码成功")
          this.refresh();
        } else {
          this.$message.error(res.msg);
        }
      })
    },
    handleDelete() {
      let delDatas = [];
      this.$refs.table.selectRow.rows.forEach(item => {
        delDatas.push(item.email);
      });
      delUser({email: delDatas[0]}).then(res => {
        if(res.status === 200) {
          this.$message.success(res.data.msg);
        } else {
          this.$message.error(res.data.msg);
        }
        this.refresh();
      })

    },
    handleExport() {
      const xlsxParam = { raw: true }; 
      const t2b = XLSX.utils.table_to_book(document.querySelector('#exportTab'), xlsxParam);
      const userExcel = XLSX.write(t2b, { bookType: 'xlsx', bookSST: true, type: 'array' });
      try {
        FileSaver.saveAs(
          new Blob([userExcel], {
            type: 'application/octet-stream' 
            }), 'userInfo.xlsx');
      } catch (e) {
          console.log(e, userExcel);
          this.$message.error("导出失败")
      }
      return userExcel;
    },
    beforeUpload(file) {
      const fileData = file.name.split(".");
      if (fileData[fileData.length-1] !== 'xlsx') {
        this.$message.error('请上传xlsx格式表格!');
        return false;
      }
      return true;
    },
    onSuccess() {
      this.$message.success("导入成功");
      this.refresh();
    },
    onError() {
      this.$message.success("导入失败");
      this.refresh();
    }
  }
}
</script>

<style scoped>
.search-form{
  margin-bottom: 12px;
}
.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
.editBtn {
  padding: 10px;
}
</style>
