<template>
  <div class="content">
    <div class="basic">
      <el-descriptions class="margin-top" :column="1">
        <el-descriptions-item label="操作:">{{ audit.action }}</el-descriptions-item>
        <el-descriptions-item label="日志:">{{ audit.message }}</el-descriptions-item>
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