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
  Date: 2022-01-17 09:41:31
  LastEditTime: 2022-05-18 09:33:27
 -->
<template>
 <div class="panel" style="height:100%">
   <transition name="fade-transform" mode="out-in">
    <router-view v-if="$route.name == 'sysctl'"></router-view>
   </transition>
   <div v-if="$route.name == 'repo'" style="height: 100%">
      <ky-table
      :getData="getRepos"
      ref="table"
    >
      <template v-slot:table_search>
        <div>repo列表</div>
      </template>
      <template v-slot:table_action>
        <el-button @click="handleCreate"> 添加 </el-button>
        <el-button :disabled="$refs.table && $refs.table.selectRow.rows.length == 0" @click="handleDelete"> 删除 </el-button>
      </template>
      <template v-slot:table>
        <el-table-column prop="id" label="编号" width="120">
        </el-table-column>
        <el-table-column  prop="repoName" label="repo名称">
        </el-table-column>
        <el-table-column  prop="size" label="文件大小"  width="120">
        </el-table-column>
        <el-table-column  prop="createdAt" label="发布时间">
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template slot-scope="scope">
            <el-button type="primary" plain size="mini" @click="handleEdit(scope.row)">编辑</el-button>   
          </template>
        </el-table-column>
      </template>
    </ky-table>

    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="display" 
      width="560px">
      <add-form v-if="type === 'create'" @click="handleClose"></add-form>
      <update-form :row="rowData" v-if="type === 'update'" @click="handleClose"></update-form>
    </el-dialog>
   </div>
 </div>
</template>
<script>
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import AddForm from "./form/addForm.vue"
import UpdateForm from "./form/updateForm.vue"
import { getRepos, createRepo, delRepos, updateRepo } from "@/request/config"
export default {
  name: "repoConfig",
  components: {
    kyTable,
    AuthButton,
    AddForm,
    UpdateForm,
  },
  data() {
    return {
      title: '',
      type: '',
      display: false,
      rowData: []
    }
  },
  methods: {
    getRepos,
    handleClose() {
      this.display = false;
      this.title = "";
      this.$refs.table.handleSearch();
    },
    handleCreate() {
      this.showRole = true;
      this.type = 'create';
      this.roleTitle = "新增repo";
    },
    handleEdit() {
      this.showRole = true;
      this.type = 'update';
      this.roleTitle = "修改repo";
    },
    handleDelete() {
      let delDatas = [];
      this.$refs.table.selectRow.rows.forEach(item => {
        delDatas.push(item.id);
      });
      delRepos({repos: delDatas}).then(res => {
        if(res.status === 200) {
          this.$message.success(res.data.msg);
          this.$refs.table.handleSearch();
        } else {
          this.$message.error(res.data.msg);
        }
      })
    }

  }
}
</script>
<style scoped lang="scss">
</style>
