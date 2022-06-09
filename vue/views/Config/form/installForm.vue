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
  LastEditTime: 2022-06-09 17:31:51
 -->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="机器:" prop="macIp">
        <el-select class="select" v-model="form.macIp" collapse-tags multiple placeholder="请选择机器ip">
          <el-option
            v-for="item in macs"
            :key="item.uuid"
            :label="item.ip_dept"
            :value="item.ip"
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
import { getallMacIps } from '@/request/cluster'
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
      ips: [],
      uuids: [],
      macs: [],
      form: {
        macIp: [],
        path: '',
        name: '',
        description: '',
        file: ''
      },
      rules: {
        macIp: [
          {
            required: true,
            message: '请选择至少一个ip',
            trigger: "blur"
          }
        ],
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
    getallMacIps().then(res => {
      this.macs = [];
      if(res.data.code === 200) {
        this.macs = res.data.data && res.data.data;
      }
    })
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleconfirm() {
      let selectMacs = [];
      this.uuids = [];
      this.form.macIp.forEach(item => {
        selectMacs.push(...this.macs.filter(mac => mac.ip === item))
      })
      selectMacs.forEach(item => {
        this.uuids.push(item.uuid)
      })
      let params = {
        uuids: this.uuids,
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