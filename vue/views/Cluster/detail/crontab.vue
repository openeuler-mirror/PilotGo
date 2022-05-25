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
  Date: 2022-04-28 14:15:30
  LastEditTime: 2022-05-24 17:10:32
 -->
    <!-- 定时任务
      1.增删改
      2.可查看状态和结果
      3.可配置时间
      4.具体的表达式
      表格：任务名、脚本、状态、创建日期、上次启动、上次花费、上次执行结果、操作（删除、开关、编辑）
     -->
<template>
 <div style="height: 100%">
  <!--  <div class="operation">
     <el-button @click="handleShowCreate" type="primary" size="medium" icon="el-icon-plus">新增</el-button>
     <el-button @click="handleDel" type="danger" size="medium" icon="el-icon-delete">删除</el-button>
   </div><br/> -->
   <ky-table
    ref="table"
    :getData="getCronList"
    :searchData="searchData"
    >
      <template v-slot:table_action>
        <el-button @click="handleShowCreate" plain icon="el-icon-plus" title="添加"></el-button>
        <el-popconfirm title="确定删除所选任务吗?" @confirm="handleDel">
          <el-button  icon="el-icon-delete" title="删除"  slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"></el-button>
        </el-popconfirm>
      </template>
    <template v-slot:table>
      <el-table-column
        style="background:rgb(109, 123, 172);"
        label="任务编号"
        prop="ID" width="60">
      </el-table-column>
      <el-table-column
        label="任务名称"
        prop="taskname">
      </el-table-column>
      <el-table-column
        label="执行脚本"
        prop="cmd" width="160">
      </el-table-column>
      <el-table-column
        label="Cron表达式"
        width="160">
        <template slot-scope="scope">
          <span class="cmd">{{scope.row.spec}}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="作业状态"
        prop="status" width="160">
        <template slot-scope="scope">
          <el-switch
            style="display: block"
            :value="scope.row.status"
            active-color="#13ce66"
            inactive-color="#ff4949"
            active-text="开启"
            inactive-text="关闭"
            :loading="true"
            @change="handleChange(scope.row)">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column
        label="创建时间"
        prop="CreatedAt">
        <template slot-scope="scope">
          <span>{{scope.row.CreatedAt | dateFormat}}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="更新时间"
        prop="UpdatedAt">
        <template slot-scope="scope">
          <span>{{scope.row.UpdatedAt | dateFormat}}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="操作" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" size="medium" plain @click="handleEdit(scope.row)">编辑</el-button>
        </template>
      </el-table-column>
    </template>
  </ky-table>
  <el-dialog 
    :title="title"
    top="10vh"
    :before-close="handleClose" 
    :visible.sync="display" 
    :width="dialogWidth"
  >
    <add-form v-if="type === 'create'" @click="handleClose"></add-form>   
    <update-form v-if="type === 'update'" :row="row" @click="handleClose"></update-form>      
  </el-dialog>
</div>
</template>
<script>
import kyTable from "@/components/KyTable";
import AddForm from "../form/cron/addForm";
import UpdateForm from "../form/cron/updateForm";
import { getCronList, changeCStatus, delCron } from '@/request/cluster';
export default {
  name: "CrontabInfo",
  components: {
    kyTable,
    AddForm,
    UpdateForm
  },
  data() {
    return {
      title: '',
      display: false,
      type: '',
      cronName: '',
      uuid: '',
      row: {},
      dialogWidth: '760px',
      filterStatus: true,
      hStyle: {
        background:'rgb(109, 123, 172)',
        color:'#fff',
        textAlign:'center',
        padding:'0',
        height: '46px',
        border: '1px solid #fff'
      },
      tableData: [],
      cronData: [],
      searchData: {
        uuid: this.$route.params.detail || 'test',
      },
    }
  },
  methods: {
    getCronList,
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = "";
      this.$refs.table.handleSearch();
    },
    handleEdit(row) {
      this.row = row;
      this.display = true;
      this.title = "编辑任务";
      this.type = "update";
    },
    handleDel() {
      let ids = this.$refs.table.selectRow.ids;
      delCron({ids:ids}).then(res => {
        if (res.data.code === 200) {
          this.$refs.table.handleSearch();
          this.$message.success("删除成功");
        } else {
          this.$message.success("删除失败");
        }
      })
    },
    handleChange(row) {
      changeCStatus({id: row.ID, uuid: this.$route.params.detail, status: row.status}).then(res => {
        if(res.data.code === 200) {
          this.$refs.table.handleSearch();
          this.$message.success("状态修改成功");
        } else {
          this.$message.success("状态修改失败");
        }
      })
    },
    handleShowCreate() {
      this.display = true;
      this.title = "新增任务";
      this.type = "create";
      this.dialogWidth = '70%';
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
<style scoped lang="scss">
.operation {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
.cmd {
  width: 100%;
  color: #fff;
  display: inline-block;
  background-color: rgb(19, 206, 102);
  border-radius: 16px;
}
</style>