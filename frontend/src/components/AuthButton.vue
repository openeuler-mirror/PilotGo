<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <el-button v-bind="$attrs" v-if="showBtn" :disabled="disabled">
    <slot></slot>
  </el-button>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, watchEffect } from 'vue';
import { hasPermisson } from '@/module/permission';

const props = defineProps({
  auth: {
    type: String,
    required: true,
    default: '',
  },
  // 是否在未授权情况下显示按钮
  show: {
    type: Boolean,
    default: false,
  }
})

const hasAuth = ref(false)
onMounted(() => {
  hasAuth.value = hasPermisson(props.auth)
})

watchEffect(() => {
  hasAuth.value = hasPermisson(props.auth)
})


const showBtn = computed(() => {
  return hasAuth.value || (!hasAuth.value && props.show)
})
// 控制是否使能按钮
const disabled = computed(() => {
  return !hasAuth.value
})

</script>

<style scoped lang="scss"></style>