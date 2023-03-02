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
  LastEditTime: 2022-06-20 13:44:22
 -->
<template>
 <div>
   <br/>
   <el-form ref="form" label-width="100px" class="selectContect">
    <el-form-item label="历史文件1:">
      <el-select class="select" v-model="leftF" placeholder="请选择文件" @change="handleLeft">
        <el-option
          v-for="item in fileOps"
          :key="item.id"
          :label="item.UpdatedAt |dateFormat"
          :value="item.file"
        />
      </el-select>
    </el-form-item>
    <el-form-item label="历史文件2:">
      <el-select class="select" v-model="rightF" placeholder="请选择文件" @change="handleRight">
        <el-option
          v-for="item in fileOps"
          :key="item.id"
          :label="item.UpdatedAt | dateFormat"
          :value="item.file"
        />
      </el-select>
    </el-form-item>
   </el-form>
   <code-diff :old-string="prev" :new-string="curr" file-name="compare.txt" output-format="side-by-side"/>
 </div>
</template>
<script>
import { lastFileList } from "@/request/config"
export default {
  props: {
    leftFile: {
      type: Object,
    },
    rightFile: {
      type: Object,
    },
    id: {
      type: Number,
      default: function() {
        return null
      }
    }
  },
  data() {
    return {
      leftF: '',
      rightF: '',
      prev: '',
      curr: '',
      fileOps: [],
    }
  },
  mounted() {
    this.leftF = this.leftFile;
    this.rightF = this.rightFile;
    this.prev = this.leftFile.file;
    this.curr = this.rightFile.file;
    lastFileList({id: this.id}).then(res => {
      if(res.data.code === 200) {
        this.fileOps = res.data.data && res.data.data;
      }
    })
  },
  methods: {
    handleLeft(value) {
      this.prev = value; 
    },
    handleRight(value) {
      this.curr = value;
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
<style scoped lang="scss">
.selectContect {
  display: flex;
  justify-content: space-around;
  align-items: center;
  .select {
    width: 300px;
  }
}
</style>
