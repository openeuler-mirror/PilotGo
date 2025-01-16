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
      <el-form-item label="脚本名称:" prop="name">
        <el-input type="text" v-model="form.name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="脚本内容:" prop="content">
        <el-input type="textarea" controls-position="right" v-model="form.content" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="脚本描述:" prop="description">
        <el-input type="textarea" controls-position="right" v-model="form.description" autocomplete="off"></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button type="primary" @click="onSubmit">确 定</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, nextTick, watch } from "vue";
import { ElMessage } from "element-plus";

import { RespCodeOK } from "@/request/request";
import { addScript, updateScript } from "@/request/script";

const formRef = ref();
const form = ref<ScriptInterface>({
  name: "",
  content: "",
  description: "",
});
const script_id = ref(0);

const props = defineProps({
  type: {
    type: String,
    default: "add",
  },
  scriptContent: {
    type: Object as () => ScriptInterface,
    defautl: { name: "", content: "", description: "" },
  },
});

interface ScriptInterface {
  name: string;
  content: string;
  description: string;
  [propName: string]: any;
}

watch(
  () => props.scriptContent,
  (newV: any) => {
    nextTick(() => {
      script_id.value = newV.id;
      form.value.name = newV.name as string;
      form.value.content = newV.content as string;
      form.value.description = newV.description as string;
    });
  },
  { immediate: true, deep: true }
);

const rules = {
  name: [
    {
      required: true,
      message: "请输入脚本名称",
      trigger: "change",
    },
  ],
  description: [
    {
      required: true,
      message: "请输入描述内容",
      trigger: "change",
    },
  ],
  content: [
    {
      required: true,
      message: "请输入脚本内容",
      trigger: "change",
    },
  ],
};

const emits = defineEmits(["updateScript", "close"]);

const onUpdateScript = () => {
  updateScript({ id: script_id.value, ...form.value })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        emits("updateScript");
        ElMessage.success("脚本更新成功");
      } else {
        ElMessage.error("脚本更新失败: " + resp.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("脚本更新失败: " + err.msg);
    });
  // emits("close");
};

// 添加脚本
const onAddScript = () => {
  addScript(form.value)
    .then((res: any) => {
      if (res.code === RespCodeOK) {
        emits("updateScript");
        ElMessage.success("脚本添加成功");
      } else {
        ElMessage.error("添加脚本失败:" + res.msg);
      }
    })
    .catch((err: any) => {
      console.log(err);
      ElMessage.error("添加脚本失败:" + err);
    });
  // emits("close");
};

// 确定方法
const onSubmit = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      props.type === "add" ? onAddScript() : onUpdateScript();
    } else {
      ElMessage.error("内容填写错误");
    }
  });
};
</script>

<style lang="scss" scoped>
.dialog-footer {
  text-align: right;
}
</style>
