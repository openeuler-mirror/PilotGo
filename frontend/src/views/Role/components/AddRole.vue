<template>
  <div>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="角色名:" prop="rolename">
        <el-input class="ipInput" type="text" v-model="form.rolename" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="描述:" prop="description">
        <el-input class="ipInput" controls-position="right" v-model="form.description" autocomplete="off"></el-input>
      </el-form-item>
    </el-form>

    <div class="footer">
      <el-button type="primary" @click="onAddRole">确 定</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

import { RespCodeOK } from "@/request/request";
import { addRole } from "@/request/role";

const emits = defineEmits(["rolesUpdated", "close"])

const rules = {
  rolename: [
    {
      required: true,
      message: "请输入角色名",
      trigger: "blur"
    }
  ],
}

const formRef = ref()
const form = ref({
  rolename: "",
  description: ""
})

function onAddRole() {
  let params = {
    role: form.value.rolename,
    description: form.value.description
  }
  formRef.value.validate((valid: Boolean) => {
    if (valid) {
      addRole(params)
        .then((resp: any) => {
          if (resp.code === RespCodeOK) {
            emits("rolesUpdated")
            formRef.value.resetFields()
            ElMessage.success(resp.msg);
          } else {
            ElMessage.error(resp.msg);
          }
        })
        .catch((err: any) => {
          ElMessage.error("添加失败,请检查输入内容" + err.msg);
        });
      emits("close")
    } else {
      ElMessage.error("请检查输入内容");
    }
  });
}

</script>

<style lang="scss" scoped>
.footer {
  text-align: right;
}
</style>