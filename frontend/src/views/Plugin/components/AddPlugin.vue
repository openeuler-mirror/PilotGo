<template>
  <div>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="插件名称:" prop="custom_name">
        <el-input class="ipInput" controls-position="right" v-model="form.custom_name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="主机地址:" prop="url">
        <el-input class="ipInput" controls-position="right" v-model="form.url" autocomplete="off"></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button type="primary" @click="onAddPlugin">确 定</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage } from 'element-plus';

import { RespCodeOK } from "@/request/request";
import { addPlugin } from '@/request/plugin';

const rules = {
  custom_name: [
    {
      required: true,
      message: '插件名称不能为空',
      trigger: "blur"
    },
  ],
  url: [
    {
      required: true,
      message: 'url不能为空',
      trigger: "blur"
    },
  ]
}

const emits = defineEmits(["pluginUpdated", "close"])

const formRef = ref()
const form = ref<any>({
  custom_name: "",
  url: "",
});

function onDepartSelected(data: any) {
  if (data) {
    form.value.departName = data.label;
    form.value.departId = data.id;
    form.value.departPid = data.pid;
  }
}

function onAddPlugin() {
  let params = {
    custom_name: form.value.custom_name,
    url: form.value.url,
  }
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      addPlugin(params).then((res: any) => {
        if (res.code === RespCodeOK) {
          emits('pluginUpdated')
          ElMessage.success(res.msg);
          formRef.value.resetFields();
        } else {
          ElMessage.error("添加插件失败:" + res.msg);
        }
      }).catch((err: any) => {
        ElMessage.error("添加插件失败2:" + err.msg);
        emits('close')
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