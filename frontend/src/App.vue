<template>
  <el-config-provider :locale="zhCn">
    <RouterView />
  </el-config-provider>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import locale from 'element-plus/es/locale/lang/zh-cn';
const zhCn = ref(locale); 
</script>

<style scoped lang="scss"></style>
