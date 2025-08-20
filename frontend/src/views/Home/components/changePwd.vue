<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Wed Mar 13 11:12:44 2024 +0800
-->
<template>
  <el-dialog v-model="showDialog" title="修改密码" width="500" @close="handleClose">
    <el-form :model="form" ref="pwdForm">
      <el-form-item label="用户：" :label-width="formLabelWidth">
        <el-input v-model="form.email" autocomplete="off" disabled />
      </el-form-item>
      <el-form-item label="新密码：" :label-width="formLabelWidth" :rules="rules">
        <el-input v-model="form.password" autocomplete="off" @change="handleConfirm(pwdForm)" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleConfirm(pwdForm)"> 确认 </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { userStore } from "@/stores/user";
import { changeUserPwd } from "@/request/user";
import { RespCodeOK } from "@/request/request";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
const formLabelWidth = "80px";
const showDialog = ref(true);
const pwdForm = ref<FormInstance>();
const form = reactive({
  email: "",
  password: "",
});
let rules = {
  required: true,
  message: "please input new password",
  trigger: "blur",
};
onMounted(() => {
  form.email = userStore().user.name as string;
});

let emit = defineEmits(["close"]);
// 修改密码
const handleConfirm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid) => {
    if (valid) {
      showDialog.value = false;
      changeUserPwd(form).then((res: any) => {
        if (res.code === RespCodeOK) {
          ElMessage.success(res.msg);
          handleClose();
        } else {
          ElMessage.error(res.msg);
        }
      });
    }
  });
};

const handleClose = () => {
  showDialog.value = false;
  emit("close");
};
</script>

<style scoped></style>
