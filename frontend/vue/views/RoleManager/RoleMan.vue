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
  Date: 2022-01-19 17:30:12
 LastEditTime: 2023-09-04 13:57:49
 -->
<template>
  <div style="width: 100%; height: 100%" class="panel">
    <ky-table
      :getData="getRolesPaged"
      :showSelect="showSelect"
      ref="table"
    >
      <template v-slot:table_search>
        <div>角色列表</div>
      </template>
      <template v-slot:table_action>
        <auth-button name="role_add"  @click="handleCreate"> 添加 </auth-button>
      </template>
      <template v-slot:table>
        <el-table-column prop="id" label="角色ID" sortable>
        </el-table-column>
        <el-table-column  prop="role" label="角色名">
        </el-table-column>
        <!-- <el-table-column  prop="username" label="平台自带">
          <template slot-scope="scope">
            {{ [0].includes(scope.row.type)  ? "是" : "否"}}
          </template>
        </el-table-column> -->
        <el-table-column  prop="description" label="描述">
        </el-table-column>
        <el-table-column label="权限">
          <template slot-scope="scope">
            <auth-button v-if="![0].includes($store.getters.userType) && ![1].includes(scope.row.id)" name="default_all" size="mini" @click="handleOpen(scope.row)">查看</auth-button>
            <auth-button v-if="[1].includes(scope.row.id)" name="default_all" size="mini" @click="handleOpen(scope.row)">查看</auth-button>
            <auth-button v-if="![1].includes(scope.row.id)" name="role_modify" size="mini" @click="handleChange(scope.row)">变更</auth-button>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right">
          <template slot-scope="scope">
            <auth-button :disabled="[1].includes(scope.row.id)" name="role_update" size="mini" @click="handleEdit(scope.row)">编辑</auth-button>
            <el-popconfirm 
              title="确定删除此角色?"
              cancel-button-type="default"
              confirm-button-type="danger"
              @confirm="handleDelete(scope.row)">
              <auth-button :disabled="[1].includes(scope.row.id)" slot="reference" name="role_delete" size="mini"> 删除 </auth-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </template>
    </ky-table>
    
    <el-drawer
      :title="title"
      :visible.sync="display"
      :direction="direction"
      :destroy-on-close="destroy"
      :before-close="handleClose">
      <role-detail @click="handleClose" :showPop="showPop" :row="row"></role-detail>
    </el-drawer>

    <el-dialog 
      :title="roleTitle"
      :before-close="handleRoleClose" 
      :visible.sync="showRole" 
      width="560px"
    >
      <add-form v-if="type === 'create'"  @click="handleRoleClose"></add-form>
      <update-form ref="updateRole" v-if="type === 'update'" :row="row"  @click="handleRoleClose"></update-form>
    </el-dialog>

  </div>
</template>

<script>
import RoleDetail from "./detail/detail.vue";
import AddForm from "./form/addForm.vue";
import UpdateForm from "./form/updateForm.vue";
import { getRolesPaged, delRole } from "@/request/role"
export default {
  components: {
    AddForm,
    UpdateForm,
    RoleDetail
  },
  data() {
    return {
      row: {},
      showRole: false,
      roleTitle: "",
      display: false,
      destroy: true,
      title: "",
      type: "",
      showPop: true,
      direction: "rtl",
      showSelect: false
    }
  },
  methods: {
    getRolesPaged,
    handleClose() {
      this.display = false;
      this.title = "";
      this.$refs.table.handleSearch();
    },
    handleRoleClose(type) {
       this.showRole = false;
       this.roleTitle = "";
       this.type = "";
       if(type === 'success') {
        this.$refs.table.handleSearch();
      }
    },
    handleCreate() {
      this.showRole = true;
      this.type = 'create';
      this.roleTitle = "新增角色";
    },
    handleEdit(row) {
      this.row = row;
      this.showRole = true;
      this.type = 'update';
      this.roleTitle = "编辑角色";
    },
    handleOpen(row) {
      this.row = row;
      this.display = true;
      this.showPop = true;
      this.title = "权限详情";
    },
    handleChange(row) {
      this.row = row;
      this.display = true;
      this.showPop = false;
      this.title = "变更权限";
    },
    handleDelete(role) {
      let params = {
          roleId: role.id,
      };
      delRole(params).then(res => {
        if(res.status === 200) {
          this.$message.success(res.data.msg);
          this.$refs.table.handleSearch();
        } else {
          this.$message.error(res.data.msg);
        }
      })
    }
  },
  
}
</script>

<style scoped>
.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
</style>
