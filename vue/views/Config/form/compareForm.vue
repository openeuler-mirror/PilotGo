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
  LastEditTime: 2022-06-08 14:35:05
 -->
<template>
 <div>
   <br/>
   <el-form ref="form" label-width="100px">
    <el-form-item label="历史文件:" prop="role">
      <el-select class="select" v-model="repoName" placeholder="请选择" @change="handleSelect">
        <el-option
          v-for="item in fileOps"
          :key="item.id"
          :label="item.name"
          :value="item.file"
        />
      </el-select>
    </el-form-item>
   </el-form>
   <code-diff :old-string="prev" :new-string="curr" file-name="compare.txt" output-format="side-by-side"/>
 </div>
</template>
<script>

export default {
  props: {
    files: {
      type: Array,
    },
    detail: {
      type: String,
    }
  },
  data() {
    return {
      prev: 'file empty',
      curr: '',
      fileOps: [],
      repoName: '',
    }
  },
  mounted() {
    this.curr = this.detail;
    this.fileOps = this.files;
  },
  methods: {
    handleSelect(value) {
      this.prev = value; 
    }
  }
}
</script>
<style scoped lang="scss">
.select {
  width: 300px;
}
</style>
