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
  Date: 2022-02-28 14:22:36
  LastEditTime: 2022-03-03 16:53:24
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div class="content">
    <div class="basic">
      <el-descriptions class="margin-top" :column="2">
        <el-descriptions-item label="类型">{{ type }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ status }}</el-descriptions-item>
        <el-descriptions-item label="进度">
        </el-descriptions-item>
      </el-descriptions>
           <el-progress type="line" :percentage="percent" :status="statusType"></el-progress>
    </div>
      <small-table
        ref="stable"
        :data="result"
        :height='theight'>
        <template v-slot:content>
        <el-table-column prop="ip" label="ip"></el-table-column>
        <el-table-column label="状态">
          <template slot-scope="scope">
            {{ scope.row.code == 200 ? '成功' : '失败' }}
          </template>
        </el-table-column>
        <el-table-column
          prop="message"
          label="原因">
        </el-table-column>
        </template>
      </small-table>
  </div>
</template>

<script>
import { getLogDetail } from '@/request/log'
import SmallTable from "@/components/SmallTable";
export default {
  name: 'LogDetail',
  components: {
    SmallTable,
  },
  props: {
    log: {
      type: Object
    }
  },
  data() {
    return {
      result: [],
      type: '',
      status: '',
      theight: 260,
      statusType: '',
      percent: 10,
    }
  },
  mounted() {
    this.status = this.log.status;
    this.statusType = this.log.status == '成功' ? 'success' : 'exception';
    getLogDetail({id: this.log.id}).then(res => {
      this.result = res.data.data;
      this.type = res.data.data[0].action;
      let errMac = this.result.filter(item => item.code == 400).length;
      if(this.log.status == '成功') {
        this.percent = 100;
      } else {
        this.percent = (errMac / this.result.length).toFixed(2) * 100;
      }
    })
  },
  methods: {
    handleCancel() {
      this.$refs.addIPForm.resetFields();
      this.$emit("click");
    },
  }
}
</script>
<style scoped lang="scss">
.content {
  display: flex;
  flex-direction: column;
  .basic{
    width: 100%;
    height: 140px;
    overflow: auto;
  }
}
</style>