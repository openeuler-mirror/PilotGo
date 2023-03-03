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
  LastEditTime: 2022-06-20 11:03:53
 -->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="文件名:" prop="name">
        <el-input
          controls-position="right"
          v-model="form.name"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="文件类型:" prop="tyep">
        <el-input
          controls-position="right"
          v-model="form.type"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="文件路径:" prop="path">
        <el-input
          controls-position="right"
          v-model="form.path"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="生效方式:" prop="activeMode">
        <el-input
          controls-position="right"
          v-model="form.activeMode"
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
import { updateLibFile } from "@/request/config";
export default {
  props: {
    row: {
      type: Object,
      default: {}
    },
    uuid: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      id: null,
      form: {
        path: '',
        name: '',
        type: '',
        activeMode: '',
        description: '',
        file: ''
      },
      rules: {
        name: [{ 
            required: true, 
            message: "请输入文件名",
            trigger: "blur" 
          }],
        type: [{ 
            required: true, 
            message: "请输入文件类型",
            trigger: "blur" 
          }],
        path: [{ 
            required: true, 
            message: "请输入文件路径",
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
    this.form.type = this.row.type;
    this.form.activeMode = this.row.activeMode || '无';
    this.form.description = this.row.description;
    this.form.file = this.row.file;
    this.id = this.row.id;
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
          updateLibFile({...this.form,id:this.id})
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