<template>
  <div v-loading="loading" class="panel" style="height:100%">
    <ky-table :getData="getUsers" ref="table" id="exportTab">
      <template v-slot:table_search>
        <el-input placeholder="请输入邮箱名进行搜索..." prefix-icon="el-icon-search" clearable
          style="width: 280px;margin-right: 10px;" v-model="emailInput" @keydown.enter.native="searchUser"></el-input>
        <el-button icon="el-icon-search" @click="searchUser">
          搜索
        </el-button>
      </template>
      <template v-slot:table_action>
        <auth-button name="user_add" @click="handleCreate"> 添加 </auth-button>
        <el-popconfirm title="确定删除此用户？" cancel-button-type="default" confirm-button-type="danger" @confirm="handleDelete">
          <auth-button name="user_del" slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0">
            删除 </auth-button>
        </el-popconfirm>
        <el-button type="primary" plain @click="handleExport"> 导出 </el-button>
        <el-upload :show-file-list="false" :before-upload="beforeUpload" :on-success="onSuccess" :on-error="onError"
          name="upload" accept="xlsx" style="display: inline-flex;margin-right: 8px" action="/user/import">
          <auth-button name="user_import"> 批量导入 </auth-button>
        </el-upload>
      </template>
      <template v-slot:table>
        <el-table-column prop="username" label="用户名">
        </el-table-column>
        <el-table-column prop="departName" label="部门">
        </el-table-column>
        <el-table-column prop="role" label="角色">
          <template slot-scope="scope">
            {{ handleRoles(scope.row.role) }}
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号">
        </el-table-column>
        <el-table-column prop="email" label="邮箱">
        </el-table-column>
        <el-table-column label="操作" fixed="right" class="operate">
          <template slot-scope="scope">
            <auth-button name="user_edit" class="editBtn" type="primary" plain size="mini"
              @click="handleEdit(scope.row)">编辑</auth-button>
            <auth-button name="user_reset" class="editBtn" type="primary" plain size="mini"
              @click="handleReset(scope.row.email)">重置密码</auth-button>
          </template>
        </el-table-column>
      </template>
    </ky-table>

    <el-dialog :title="title" :before-close="handleClose" :visible.sync="display" width="560px">
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
import { getUsers, delUser, resetPwd, searchUser } from "@/request/user"
export default {
  components: {
    AddForm,
    UpdateForm,
  },
  data() {
    return {
      loading: false,
      display: false,
      isDelete: true,
      emailInput: '',
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
      if (type === 'success') {
        this.refresh();
      }
    },
    handleRoles(roles) {
      if (roles && roles.length > 0) {
        let roleString = "";
        roles.forEach((item, index) => {
          if (index == roles.length - 1) {
            roleString += item;
          } else {
            roleString = item + "," + roleString;
          }
        });
        return roleString;
      } else {
        return "暂无";
      }
    },
    refresh() {
      this.$refs.table.handleSearch();
    },
    searchUser() {
      searchUser({ 'email': this.emailInput }).then((res) => {
        if (res.data.code === 200) {
          this.$refs.table.handleLoadSearch(res.data.data);
        }
      })
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
      resetPwd({ 'email': email }).then((res) => {
        if (res.data.code === 200) {
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
      delUser({ delDatas: delDatas }).then(res => {
        if (res.status === 200) {
          this.$message.success(res.data.msg);
          this.refresh();
        } else {
          this.$message.error(res.data.msg);
        }
      })

    },
    handleExport() {
      const xlsxParam = { raw: true };
      const exportTabElement = document.querySelector('#exportTab');
      const fixed = exportTabElement.querySelector(".el-table__fixed") || exportTabElement.querySelector(".el-table__fixed-right");
      let t2b = null;
      console.log(exportTabElement, 33333)
      if (fixed) {
        const parentNode = fixed.parentNode;
        parentNode.removeChild(fixed);
        t2b = XLSX.utils.table_to_book(exportTabElement, xlsxParam);
        parentNode.appendChild(fixed);
      } else {
        t2b = XLSX.utils.table_to_book(exportTabElement, xlsxParam);
      }
      t2b.Sheets.Sheet1['!cols'][6] = { hidden: true }
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
      if (fileData[fileData.length - 1] !== 'xlsx') {
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
.search-form {
  /* margin-bottom: 12px; */
}

.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
</style>
