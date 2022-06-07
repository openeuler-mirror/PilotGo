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
  LastEditTime: 2022-06-07 10:46:53
 -->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
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
          placeholder="请按格式输入"
          v-model="form.file"
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
import { saveRepo, getRepoDetail } from "@/request/config";
export default {
  props: {
    row: {
      type: Object,
      default: {}
    },
    uuid: {
      type: String,
      default: ''
    },
    macIp: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      form: {
        path: '',
        name: '',
        description: '',
        file: ''
      },
      rules: {
        path: [
          { 
            required: true, 
            message: "请输入路径",
            trigger: "blur" 
          }],
        name: [{ 
            required: true, 
            message: "请输入repo名",
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
    getRepoDetail({uuid: this.uuid,file: this.row.path+'/'+this.row.name}).then(res => {
      if(res.data.code === 200) {
        this.form.file = res.data && res.data.data.file;
      }
    })
    this.form.name = this.row.name;
    this.form.path = this.row.path;
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSet() {
      
    },
    handleconfirm() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          saveRepo({...this.form,ip:this.macIp})
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click","success");
                this.$message.success(res.data.msg);
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("添加失败, 请检查输入内容");
            });
        }
      });
    },
  },
};
</script>