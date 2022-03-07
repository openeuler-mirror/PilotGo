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
  Date: 2022-02-25 16:33:45
  LastEditTime: 2022-03-03 17:09:25
  Description: provide agent log manager of pilotgo
 -->
<template>
 <div>
   <ky-table
        class="cluster-table"
        ref="table"
        :getData="getLogs"
      >
        <template v-slot:table_search>
          <div>日志列表</div>
        </template>
        <template v-slot:table_action>
          <el-popconfirm 
          title="确定删除此日志？"
          cancel-button-type="default"
          confirm-button-type="danger"
          @confirm="handleDelete">
          <el-button slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 删除 </el-button>
        </el-popconfirm>
        </template>
        <template v-slot:table>
          <el-table-column prop="type" label="日志名称">
          </el-table-column>
          <el-table-column prop="userName" label="创建者"> 
          </el-table-column>
          <el-table-column prop="status" label="状态"> 
          </el-table-column>
          <el-table-column prop="CreatedAt" label="创建时间">
            <template slot-scope="scope">
              <span>{{scope.row.CreatedAt | dateFormat}}</span>
            </template>
          </el-table-column>
          <el-table-column prop="operation" label="详情">
            <template slot-scope="scope">
              <el-button
                size="mini"
                type="primary"
                plain
                @click="handleDetail(scope.row)">
                查看
              </el-button>
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
        <log-detail v-if="type === 'detail'" :log="log" @click="handleClose"></log-detail>
      </el-dialog>
 </div>
</template>
<script>
import kyTable from "@/components/KyTable";
import LogDetail from "./form/detail.vue"
import { getLogs, deleteLog } from "@/request/log";
export default {
  name: "Log",
  components: {
    kyTable,
    LogDetail,
  },
  data() {
    return {
      display: false,
      title: '',
      type: '',
      log: {},
    }
  },
  methods: {
    getLogs,
    refresh(){
      this.$refs.table.handleSearch();
    },
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = "";
    },
    handleDetail(row) {
      // 查看日志详情
      this.display = true;
      this.title = "日志详情";
      this.type = "detail";
      this.log = row;
    },
    handleDelete() {
      deleteLog({ids: this.$refs.table.selectRow.ids}).then(res => {
        if(res.status === 200) {
          this.$message.success(res.data.msg);
          this.refresh();
        } else {
          this.$message.error(res.data.msg);
        }
      })
    }
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
</style>
