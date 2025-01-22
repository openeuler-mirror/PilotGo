<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Mar 19 11:41:25 2024 +0800
-->
<template>
  <div class="content">
    <div class="basic">
      <el-descriptions class="margin-top" :column="1">
        <el-descriptions-item label="操作:">{{ audit.action }}</el-descriptions-item>
        <el-descriptions-item label="日志:"><pre>{{ audit.message }}</pre></el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';

const props = defineProps({
  audit: {
    type: Object,
    required: true,
    default: {}
  }
})

let audit = ref({})
onMounted(() => {
  audit.value = props.audit;
})
</script>

<style scoped lang="scss">
.content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;

  .basic {
    width: 100%;
    text-align: center;
  }

  .tab {
    width: 100%;
    height: 90%;
  }
}
</style>