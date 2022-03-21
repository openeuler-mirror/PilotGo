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
  LastEditTime: 2022-03-17 15:02:50
 -->
<template>
  <div>
    <ky-table
      :getData="getRoles"
      :searchData="searchData"
      ref="table"
    >
      <template v-slot:table_search>
        <div>角色列表</div>
      </template>
      <template v-slot:table>
        <el-table-column prop="id" label="编号">
        </el-table-column>
        <el-table-column  prop="role" label="角色名">
        </el-table-column>
        <el-table-column  prop="username" label="平台自带">
          <template slot-scope="scope">
            {{ [0,1,2].includes(scope.row.type)  ? "是" : "否"}}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button class="editBtn"  type="primary" plain size="mini" @click="handleOpen(scope.row)">查看权限</el-button>
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
      <role-detail :row="row"></role-detail>
    </el-drawer>

  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import RoleDetail from "./detail/detail.vue";
import { getRoles } from "@/request/role"
export default {
  components: {
    kyTable,
    RoleDetail
  },
  data() {
    return {
      row: {},
      display: false,
      destroy: true,
      title: "",
      direction: "rtl",
      searchData: {
        showSelect: false
      },
    }
  },
  methods: {
    getRoles,
    handleClose() {
      this.display = false;
      this.row = {};
      this.title = "";
    },
    handleOpen(row) {
      this.row = row;
      this.display = true;
      this.title = "权限详情";
    },
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
