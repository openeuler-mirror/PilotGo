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
      <el-form-item label="脚本内容:" prop="content">
        <el-input
          type="textarea"
          disabled
          controls-position="right"
          v-model="form.content"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="目标类型:" :prop="machineType">
        <el-radio-group v-model="machineType" @change="changeMachineType">
          <el-radio value="host">主机</el-radio>
          <el-radio value="batch">批次</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="选择主机:" prop="machine_uuids" v-if="machineType === 'host'">
        <el-select v-model="form.machine_uuids" multiple placeholder="请选择主机IP">
          <el-option v-for="item in hosts" :label="item.ip" :value="item.uuid" :key="item.ip" />
        </el-select>
      </el-form-item>
      <el-form-item label="选择批次:" prop="batch_id" v-if="machineType === 'batch'">
        <el-select v-model="form.batch_id" placeholder="请选择批次">
          <el-option v-for="item in batches" :label="item.name" :value="item.id" :key="item.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="历史版本">
        <el-select v-model="form.version" placeholder="请选择脚本历史版本">
          <el-option v-for="item in historyList" :label="item.version" :value="item.version" :key="item.version" />
        </el-select>
      </el-form-item>
      <el-form-item label="脚本参数" prop="params">
        <el-input-tag
          v-model="form.params"
          placeholder="输入参数后回车"
          aria-label="Please click the Enter key after input"
        />
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button @click="onCancle">取消</el-button>
      <el-button type="primary" @click="onSubmit">立即执行</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, nextTick, watch } from "vue";
import { ElMessage } from "element-plus";

import { RespCodeOK } from "@/request/request";
import { runScript, getScriptHistorys } from "@/request/script";
import { getAllMachines } from "@/request/cluster";
import { getAllBatches } from "@/request/batch";
import type { HistoryItem } from "@/types/script";
import type { MachineInfo } from "@/types/cluster";
import type { BatchItem } from "@/types/batch";

// todo:引入主机和批次不分页接口，下午
const emits = defineEmits(["close"]);
const props = defineProps({
  scriptContent: {
    type: Object,
    defautl: { id: 0, content: "" },
  },
});
interface ScriptInterface {
  batch_id?: number | null;
  script_id: number | null;
  machine_uuids?: number[];
  content: string;
  version?: string;
  params: string[];
  [propName: string]: any;
}
const form = ref<ScriptInterface>({
  script_id: null,
  content: "",
  machine_uuids: [],
  batch_id: null,
  version: "",
  params: [],
});
const formRef = ref();
const script_id = ref(0);
const machineType = ref("host");

watch(
  () => props.scriptContent,
  (newV: any) => {
    if (!newV) return;
    nextTick(() => {
      script_id.value = form.value.script_id = newV.id;
      form.value.content = newV.content as string;
      getHostResource();
    });
  },
  { immediate: true, deep: true }
);

// 变更选择主机类型
const changeMachineType = (value: any) => {
  console.log(value);
};

// 获取form表单列表选项
const hosts = ref<MachineInfo[]>();
const batches = ref<BatchItem[]>();
const historyList = ref<HistoryItem[]>();
const getHostResource = () => {
  // 获取所有机器列表
  getAllMachines()
    .then((res: any) => {
      if (res.code === RespCodeOK) {
        hosts.value = res.data;
      } else {
        ElMessage.error("failed to machine list:" + res.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to machine lists:" + err.msg);
    });

  // 获取所有批次列表
  getAllBatches()
    .then((res: any) => {
      if (res.code === RespCodeOK) {
        batches.value = res.data;
      } else {
        ElMessage.error("failed to batch list:" + res.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to batch lists:" + err.msg);
    });

  // 获取历史版本列表
  getScriptHistorys({ script_id: script_id.value })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        historyList.value = resp.data;
        console.log(historyList.value);
      } else {
        ElMessage.error("failed to history list:" + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to history lists:" + err.msg);
    });
};

const rules = {
  content: [
    {
      required: true,
      message: "请输入脚本内容",
      trigger: "change",
    },
  ],
  machine_uuids: [
    {
      required: true,
      message: "请选择主机ip",
      trigger: "change",
    },
  ],
  batch_id: [
    {
      required: true,
      message: "请选择批次",
      trigger: "change",
    },
  ],
};

const onRunScript = () => {
  runScript({ id: script_id.value, ...form.value })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        emits("close");
        ElMessage.success("脚本执行成功");
      } else {
        ElMessage.error("脚本执行失败: " + resp.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("脚本执行失败: " + err.msg);
    });
};

// 立即执行
const onSubmit = () => {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      onRunScript();
    } else {
      ElMessage.error("内容填写错误");
    }
  });
};

// 取消
const onCancle = () => {
  emits("close");
};
</script>

<style lang="scss" scoped>
.dialog-footer {
  text-align: right;
}
</style>
