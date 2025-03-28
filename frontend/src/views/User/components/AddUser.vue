<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="用户名:" prop="userName">
        <el-input class="ipInput" type="text" v-model="form.userName" autocomplete="off" show-word-limit
          maxlength="10"></el-input>
      </el-form-item>
      <el-form-item label="密码:" prop="password">
        <el-input type="password" class="ipInput" controls-position="right" v-model="form.password"
          autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="部门:" prop="departName">
        <el-input class="ipInput" controls-position="right" v-model="form.departName" autocomplete="off"></el-input>
        <PGTree style="width: 98%;" :showHeader="false" @onNodeClicked="onDepartSelected">
        </PGTree>
      </el-form-item>
      <el-form-item label="用户角色:" prop="role">
        <el-select v-model="form.role" multiple placeholder="请选择">
          <el-option v-for="item in roles" :key="item.id" :label="item.role" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="手机号:" prop="phone">
        <el-input class="ipInput" controls-position="right" v-model="form.phone" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="邮箱:" prop="email">
        <el-input class="ipInput" controls-position="right" v-model="form.email" autocomplete="off"></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button type="primary" @click="onAddUser">确 定</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';

import PGTree from "@/components/PGTree.vue";

import { checkEmail, checkPhone } from "./logic";
import { RespCodeOK } from "@/request/request";
import { getRoles } from '@/request/role';
import { addUser } from '@/request/user';

const rules = {
  userName: [{
    required: true,
    message: "请输入用户名",
    trigger: "change"
  }],
  password: [{
    required: true,
    message: "请输入密码",
    trigger: "change"
  }],
  departName: [{
    required: true,
    message: "请选择部门",
    trigger: "change"
  }],
  role: [{
    required: true,
    message: "请选择角色",
    trigger: "change"
  }],
  phone: [{
    required: false,
    message: "请输入手机号",
    trigger: "change",
  }, {
    validator: checkPhone,
    message: "请输入正确的手机号格式",
    trigger: "change",
  }],
  email: [{
    required: true,
    message: "请输入邮箱",
    trigger: "change",
  },
  {
    validator: checkEmail,
    message: "请输入正确的邮箱格式",
    trigger: "change",
  }],
}

const emits = defineEmits(["userUpdated", "close"])

const roles = ref<any[]>();

function updateRoles() {
  getRoles().then((resp: any) => {
    if (resp.code === RespCodeOK) {
      roles.value = resp.data
    } else {
      ElMessage.error("failed to get roles info: " + resp.msg)
    }
  }).catch((err) => {
    ElMessage.error("failed to get roles info: " + err.msg)
  })
}

onMounted(() => {
  updateRoles()
})

const formRef = ref()
const form = ref<any>({
  userName: "",
  password: "",
  phone: "",
  email: "",
  departName: "",
  departId: "",
  departPid: "",
  role: "",
});

function onDepartSelected(data: any) {
  if (data) {
    form.value.departName = data.label;
    form.value.departId = data.id;
    form.value.departPid = data.pid;
  }
}

function onAddUser() {
  let params = {
    userName: form.value.userName,
    password: form.value.password,
    phone: form.value.phone,
    email: form.value.email,
    departName: form.value.departName,
    departId: form.value.departId,
    departPid: form.value.departPid,
    // TODO: fix this
    roleid: form.value.role.toString(),
  }
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      addUser(params).then((res: any) => {
        if (res.code === RespCodeOK) {
          emits('userUpdated')
          formRef.value.resetFields();
          ElMessage.success(res.msg);
        } else {
          ElMessage.error("添加用户失败:" + res.msg);
        }
      }).catch((err: any) => {
        ElMessage.error("添加用户失败:" + err.msg);
      });
      emits('close')
    } else {
      ElMessage.error("内容填写错误");
    }
  });
}


</script>

<style lang="scss" scoped>
.dialog-footer {
  text-align: right;
}
</style>