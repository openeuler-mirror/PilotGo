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
  Date: 2022-03-01 09:35:54
  LastEditTime: 2022-03-02 15:25:39
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div class="content">
    <div class="mac">
      <el-card>
        <div slot="header" class="clearfix">
          <span>机器列表</span>
        </div>
          <div v-for="item in macs" :key="item.id" class="text item">
            {{ item.ip }}
          </div>
        </el-card>
    </div>
    <div class="action">
      <el-input class="input"  placeholder="请输入rpm名称" v-model="input" clearable></el-input>
      <el-button :disabled="input.length == 0" type="primary"  @click="handleRpm"> {{ btnName }} </el-button>
      <div class="info">
        <el-descriptions class="margin-top" :title='title + "结果"' size="medium" :column="4">
          <el-descriptions-item label="结果" :span="4">{{ result }}</el-descriptions-item>
          <el-descriptions-item label="详情">具体请查看日志</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
  </div>
</template>
<script>
import { rpmIssue, rpmUnInstall } from '@/request/cluster'

export default {
  props: {
    machines: {
      type: Array
    },
    acType: {
      type: String
    }
  },
  data() {
    return {
      input: '',
      macs: [],
      ids: [],
      btnName: '下发',
      title: '下发结果',
      result: '暂无结果',
    }
  },
  mounted() {
    this.macs = this.machines;
    this.macs.forEach(item => {
      this.ids.push(item.uuid);
    })
    this.btnName = this.acType == '软件包下发' ? '下发' : '卸载';
    this.title = this.acType == '软件包下发' ? '下发' : '卸载';
  },
  methods: {
    handleRpm() {
      let params = {
        uuid: this.ids,
        rpm: this.input,
        userName: this.$store.getters.userName
      }
      this.acType == '软件包下发' ? this.handleIssue(params) : this.handleUnInstall(params);
    },
    handleResult(res) {
      this.result = res.data.data.code === 200 ? '成功' : '失败'
    },
    handleIssue(params) {
      rpmIssue(params).then(res => {
        this.handleResult(res)
      }).catch(error => {
        console.log("api error")
      })
    },
    handleUnInstall(params) {
      rpmUnInstall(params).then(res => {
        this.handleResult(res)
      }).catch(error => {
        console.log("api error")
      })
    }
  }
}
</script>
<style scoped lang="scss">
  .content {
    display: flex;
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
    .mac {
      width: 32%;
      max-height: 400px;
      overflow:auto;
      .clearfix {
        font-size: 16px;
      }
      .text {
        font-size: 16px;
      }
    }
    .action {
      width: 60%;
      .input {
        width: 70%;
      }
      .info {
        margin: 16px 0 0 10px;
      }
    }
  }
</style>
