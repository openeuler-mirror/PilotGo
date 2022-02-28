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
  LastEditTime: 2022-02-28 15:28:25
  Description: provide agent log manager of pilotgo
 -->
<template>
 <div>
   <ky-table
        class="cluster-table"
        ref="table"
        :searchData="searchData"
        :getData="getBatchDetail"
      >
        <template v-slot:table_search>
          <div>机器列表</div>
        </template>
        <template v-slot:table>
         <el-table-column prop="ip" label="IP">
          </el-table-column>
          <el-table-column prop="CPU" label="cpu"> 
          </el-table-column>
          <el-table-column label="状态" width="150">
            <template slot-scope="scope">
              <span v-if="scope.row.state == 1">正常</span>
              <span v-if="scope.row.state == 2">离线</span>
              <span v-if="scope.row.state == 3">空闲</span>
            </template>
          </el-table-column>
           <el-table-column prop="sysinfo" label="系统信息"> 
          </el-table-column>
        </template>
      </ky-table>
 </div>
</template>
<script>
import { getBatchDetail } from "@/request/batch";
import kyTable from "@/components/KyTable";
export default {
  name: "BatchDetail",
  components: {
    kyTable,
  },
  data() {
    return {
      batchTitle: '',
      searchData: {
        ID: 0,
        load: 'false',
        showSelect: false
      }
    }
  },
  mounted() {
    this.searchData.ID = JSON.parse(this.$route.params.id);
    this.searchData.load = 'true';
    this.$refs.table.handleSearch();
  },  
  methods: {
    getBatchDetail,
  }
}
</script>
<style scoped>
</style>
