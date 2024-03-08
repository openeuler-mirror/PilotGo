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
  Date: 2022-03-24 15:59:44
  LastEditTime: 2022-06-14 14:23:38
 -->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="告警邮箱:" prop="email">
        <el-select v-model="form.email" multiple placeholder="请选择">
          <el-option
            v-for="item in users"
            :key="item.id"
            :label="item.email"
            :value="item.email"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="告警名称:" prop="alertname">
        <el-input
          class="ipInput"
          type="text"
          size="medium"
          v-model="form.alertname"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="告警IP:" prop="instance">
        <el-input
          class="ipInput"
          controls-position="right"
          :disabled="disabled"
          v-model="form.instance"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="产生时间:" prop="activeAt">
        <el-input
          class="ipInput"
          type="text"
          size="medium"
          :disabled="disabled"
          v-model="form.activeAt"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="告警内容:" prop="annotations">
        <el-input
          class="ipInput"
          controls-position="right"
          v-model="form.annotations"
          autocomplete="off"
        ></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
  </div>
</template>

<script>
import { sendMessage } from "@/request/overview";
import { getUsers } from "@/request/user"
import { checkEmail } from "@/rules/check"
export default {
  props: {
    row: {
      type: Object,
    }
  },
  data() {
    return {
      disabled: true,
      users: [],
      showConfirm: false,
      form: {
        email: '',
        alertname: '',
        instance: '',
        activeAt: '',
        annotations: '',
        startTime: '',
        endTime: '',
      },
      rules: {
        alertname: [
          { 
            required: true, 
            message: "请输入告警名称",
            trigger: "blur" 
          }],
        instance: [{ 
            required: true, 
            message: "请输入IP",
            trigger: "blur" 
          }],
        annotations: [
          {
            required: true,
            message: "请输入告警内容",
            trigger: "blur",
          }],
        email: [
          {
            required: true,
            message: "请选择至少一个邮箱",
            trigger: "blur",
          },],
      },
    };
  },
  mounted() {
    if(this.row) {
      this.form.alertname = this.row.labels.alertname;
      this.form.instance = this.row.labels.instance.split(':')[0];
      this.form.annotations = this.row.annotations.description;
      this.form.activeAt = this.row.activeAt;
    }
    getUsers().then(res => {
      if(res.data.code === 200) {
        this.users = res.data.data;
      }
    })
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleConfirm() {
      let params = {
        "Email": this.form.email,
        "Labels": {
          "alertname": this.form.alertname,
          "IP": this.form.instance,
        },
        "Annotations": {
          "summary": this.form.annotations,
        },
        "StartsAt": new Date().toISOString(),
        "EndsAt": new Date().toISOString(),
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          sendMessage(params)
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click");
                this.$message.success(res.data.msg);
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.msg);
              }
            })
            .catch((res) => {
              this.$message.error("发送失败");
            });
        } else {
          this.$message.error("发送失败");
          return false;
        }
      });
    },
  },
};
</script>