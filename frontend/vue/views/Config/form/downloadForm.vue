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
  LastEditTime: 2022-07-01 14:23:34
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
      <el-form-item label="文件类型:" prop="type">
        <el-select class="select" v-model="form.type" placeholder="请选择文件类型">
          <el-option
            v-for="item in types"
            :key="item.$index"
            :label="item"
            :value="item"
          >
          </el-option>
        </el-select>
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
      default: null
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
      types: [],
      form: {
        name: '',
        path: '',
        type: [],
        activeMode: '',
        description: '',
        file: ''
      },
      rules: {
        name: [{ 
            required: true, 
            message: "请输入repo名",
            trigger: "blur" 
          }],
        path: [{ 
            required: true, 
            message: "请输入文件路径",
            trigger: "blur" 
          }],
        type: [{ 
            required: true, 
            message: "请选择文件类型",
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
    this.types = this.$store.getters.configType;
    if(this.row) {
      getRepoDetail({uuid: this.uuid,file: this.row.path+'/'+this.row.name}).then(res => {
        if(res.data.code === 200) {
          this.form.file = res.data && res.data.data.file;
        }
      })
      this.form.name = this.row.name;
      this.form.path = this.row.path;
    }
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSet() {
      
    },
    handleconfirm() {
      let params = {
        user: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          saveRepo({...this.form,...params})
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click","success");
                this.$message.success(res.data.msg);
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.msg);
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