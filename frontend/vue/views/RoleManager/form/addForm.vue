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
  LastEditTime: 2022-04-06 15:11:15
 -->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="角色名:" prop="rolename">
        <el-input
          class="ipInput"
          type="text"
          size="medium"
          v-model="form.rolename"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="描述:" prop="description">
        <el-input
          class="ipInput"
          controls-position="right"
          v-model="form.description"
          autocomplete="off"
        ></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleAdd">确 定</el-button>
    </div>
  </div>
</template>

<script>
import { addRole } from "@/request/role";
export default {
  data() {
    return {
      form: {
        rolename: '',
        description: "",
      },
      rules: {
        rolename: [
          { 
            required: true, 
            message: "请输入角色名",
            trigger: "blur" 
          }],
        
      },
    };
  },
  methods: {
    handleCancel() {
      this.refresh();
      this.$emit("click");
    },
    refresh() {
      this.$refs.form.resetFields();
    },
    handleAdd() {
      let params = {
        role: this.form.rolename,
        description: this.form.description
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          addRole(params)
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click","success");
                this.$message.success(res.data.msg);
                this.refresh();
              } else {
                this.$message.error(res.data.msg);
              }
            })
            .catch((res) => {
              this.$message.error("添加失败,请检查输入内容");
            });
        }
      });
    },
  },
};
</script>