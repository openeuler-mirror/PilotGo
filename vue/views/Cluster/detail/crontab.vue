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
  LastEditTime: 2022-04-29 15:48:56
 -->
    <!-- 定时任务
      1.增删改
      2.可查看状态和结果
      3.可配置时间
      4.具体的表达式
      表格：任务名、脚本、状态、创建日期、上次启动、上次花费、上次执行结果、操作（删除、开关、编辑）
     -->
<template>
 <div>
   <div class="filter">
     <el-autocomplete
        style="width:50%"
        class="inline-input"
        v-model="cronName"
        :fetch-suggestions="querySearch"
        placeholder="请输入任务名称"
      >
      <el-switch
          style="display: block"
          active-color="#13ce66"
          inactive-color="#ff4949"
          active-text="开启"
          inactive-text="禁用"
          slot="append"
          @change="handlefilterChange">
        </el-switch>
      </el-autocomplete>
     <el-button @click="handleShowCreate" v-show="showAddBtn">新增</el-button>
   </div>
   <br/>
   <el-form v-show="showForm" :model="form" label-width="80px">
      <el-form-item style="margin-top: -10px; margin-bottom:0px;">
       <span style="color: rgb(241, 139, 14); font-size: 14px;">corn从左到右(用空格隔开):秒 分 小时 月份中的日期 月份 星期中的日期 年份</span>
       <cron v-if="showCronBox" v-model="form.cronExpression"></cron>
     </el-form-item>
     <el-form-item label="Cron:">
       <el-input v-model="form.cronExpression" auto-complete="off" @focus="showCronBox = true" @blur="showCronBox = false">
          <el-button slot="append" @click="handleCreate" title="确定">确定</el-button>
       </el-input>
     </el-form-item>
    </el-form>
   <el-table
    :data="tableData"
    :header-cell-style="hStyle"
    style="width: 100%">
    <el-table-column
      style="background:rgb(109, 123, 172);"
      label="任务编号"
      prop="id">
    </el-table-column>
    <el-table-column
      label="任务名称"
      prop="name">
    </el-table-column>
    <el-table-column
      label="Cron表达式"
      prop="cron">
    </el-table-column>
    <el-table-column
      label="作业状态"
      prop="status">
      <template slot-scope="scope">
        <el-switch
          style="display: block"
          :value="scope.row.status === 1"
          active-color="#13ce66"
          inactive-color="#ff4949"
          active-text="开启"
          inactive-text="禁用"
          @change="handleChange(scope.row)">
        </el-switch>
      </template>
    </el-table-column>
    <el-table-column
      label="创建时间"
      prop="createdAt">
    </el-table-column>
    <el-table-column
      label="更新时间"
      prop="updatedAt">
    </el-table-column>
    <el-table-column
      label="操作">
      <template slot-scope="scope">
        <el-button type="primary" size="medium" icon="el-icon-edit" circle @click="handleEdit(scope.row)"></el-button>
        <el-button type="danger" size="medium" icon="el-icon-delete" circle @click="handleDel(scope.row)"></el-button>
      </template>
    </el-table-column>
  </el-table>
    
  </div>
</template>
<script>
import cron  from '@/components/VueCron';
import { getCronList, createCron, updateCron, delCron } from '@/request/cluster';
export default {
  name: "CrontabInfo",
  components: { 
    cron,
  },
  data() {
    return {
      cronName: '',
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
      showForm: false,
      showAddBtn: true,
      showCronBox: false,
      form: {
        cronExpression: '',
      }
    }
  },
  async mounted() {
    if(this.$route.params.detail != undefined) {
      await this.$http.get('/api/cron_list').then(res => {
        if(res.data.code == 200) {
          this.tableData = res.data.cron_info;
          this.tableData.forEach(item => {
            this.cronData.push({'value': item.name})
          })
        }
      })
    }
  },
  methods: {
    querySearch(queryString, cb) {
      let cronData = this.cronData;
      var results = queryString ? cronData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): cronData;
      cb(results);
    },
    handleEdit(row) {

    },
    handleDel(row) {

    },
    handlefilterChange() {
      // 
      // this.cronName
    },
    handleChange(row) {

    },
    handleCreate() {
      this.showAddBtn = true;
      this.showForm = false;
      console.log(this.cronName)
    },
    handleShowCreate() {
      this.showForm = true;
      this.showAddBtn = false;
    }


  }
}
</script>
<style scoped lang="scss">
.demo-table-expand {
    font-size: 0;
  }
  .demo-table-expand label {
    width: 90px;
    color: #99a9bf;
  }
  .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    width: 50%;
  }
</style>