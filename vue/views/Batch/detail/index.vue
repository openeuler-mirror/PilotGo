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
  LastEditTime: 2022-05-19 14:54:49
  Description: provide agent log manager of pilotgo
 -->
<template>
 <div style="height: 100%">
   <ky-table
        class="cluster-table"
        ref="table"
        :showSelect="showSelect"
        :searchData="searchData"
        :getData="getBatchDetail"
      >
        <template v-slot:table_search>
          <div>机器列表</div>
        </template>
        <template v-slot:table_action>
          <auth-button name="rpm_install"  @click="handleIssue()"> rpm下发 </auth-button>
          <auth-button name="rpm_uninstall" @click="handleUnInstall()"> rpm卸载</auth-button>
        </template>
        <template v-slot:table>
         <el-table-column prop="ip" label="IP">
          </el-table-column>
          <el-table-column prop="CPU" label="cpu"> 
          </el-table-column>
          <el-table-column label="状态" width="150">
            <template slot-scope="scope">
              <state-dot :state="scope.row.state"></state-dot>
            </template>
          </el-table-column>
           <el-table-column prop="sysinfo" label="系统信息"> 
          </el-table-column>
        </template>
      </ky-table>
      <el-dialog 
        :title="title"
        :before-close="handleClose" 
        :visible.sync="display" 
        width="560px"
      >
        <rpm-issue v-if="type === 'issue'" :acType='title' :machines='machines' @click="handleClose"></rpm-issue>
      </el-dialog>
 </div>
</template>
<script>
import { getBatchDetail } from "@/request/batch";
import RpmIssue from "../form/rpmIssue";
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import StateDot from "@/components/StateDot";
export default {
  name: "BatchDetail",
  components: {
    kyTable,
    AuthButton,
    RpmIssue,
    StateDot,
  },
  data() {
    return {
      display: false,
      title: '',
      type: '',
      batchTitle: '',
      machines: [],
      searchData: {
        ID: JSON.parse(this.$route.params.id) || 0
      },
      showSelect: false,
    }
  },
  mounted() {
    getBatchDetail({ ID: parseInt(this.$route.params.id) }).then(res => {
      this.machines = [];
      if(res.data.code === 200) {
        this.machines = res.data.data;
      }
    })
  },  
  methods: {
    getBatchDetail,
    handleClose(type) {
      this.display = false;
      this.title = "";
      this.type = "";
    },
    handleIssue() {
      this.display = true;
      this.title = "软件包下发";
      this.type = "issue"; 
    },
    handleUnInstall() {
      this.display = true;
      this.title = "软件包卸载";
      this.type = "issue"; 
    },
  }
}
</script>
<style scoped>
</style>
