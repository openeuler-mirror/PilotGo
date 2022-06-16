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
  Date: 2022-02-10 09:37:29
  LastEditTime: 2022-06-16 14:47:27
 -->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="批次:" prop="batches">
        <el-select class="select" v-model="form.batches" collapse-tags multiple placeholder="请选择下发批次">
          <el-option
            v-for="item in batches"
            :key="item.ID"
            :label="item.name"
            :value="item.ID"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="路径:" prop="path">
        <el-input
          type="text"
          size="medium"
          v-model="form.path"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="文件名:" prop="name">
        <el-input
          controls-position="right"
          v-model="form.name"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="描述:" prop="description">
        <el-input
          controls-position="right"
          v-model="form.description"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="内容:" prop="file">
        <el-input
          type="textarea"
          v-model="form.file"
          disabled
        ></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleconfirm">确 定</el-button>
    </div>
  </div>
</template>

<script>
import { getBatches } from "@/request/batch";
import { installFile } from "@/request/config"
export default {
  props: {
    row: {
      type: Object,
      default: {}
    }
  },
  data() {
    return {
      batches: [],
      form: {
        batches: [],
        path: '',
        name: '',
        description: '',
        file: ''
      },
      rules: {
        batches: [
          { 
            required: true, 
            message: "至少选择一个批次",
            trigger: "blur" 
          }],
        path: [
          { 
            required: true, 
            message: "请输入路径",
            trigger: "blur" 
          }],
        name: [{ 
            required: true, 
            message: "请输入文件名",
            trigger: "blur" 
          }],
        description: [{ 
            required: true, 
            message: "请输入具体的描述",
            trigger: "blur" 
          }],
        file: [{
            required: true,
            message: "请输入文件内容",
            trigger: "blur",
          }],
      },
    };
  },
  mounted() {
    this.form.name = this.row.name;
    this.form.path = this.row.path;
    this.form.description = this.row.description;
    this.form.file = this.row.file;
    getBatches().then(res => {
      this.batches = [];
      if(res.data.code === 200) {
        this.batches = res.data.data && res.data.data;
      }
    })
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleconfirm() {
      let params = {
        batches:this.form.batches,
        path: this.form.path,
        name: this.form.name,
        file: this.form.file,
        user: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          installFile(params).then(res => {
            if(res.data.code === 200) {
              this.$message.success(res.data.msg)
              this.$emit("click");
            } else {
              this.$message.error(res.data.msg)
            }
          })
        }
      });
    },
  },
};
</script>
<style lang="scss" scoped>
  .select {
    width: 300px;
  }
</style>