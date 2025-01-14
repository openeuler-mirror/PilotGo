<template>
  <svg aria-hidden="true" class="svgSize">
    <use :xlink:href="symbolId" :fill="color" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from "vue";
/**
 * name: assets/icon 路径下svg图标name不带.svg部分
 * color: 非必须,若需要填充纯色，先清空svg中的fill属性，改为""
 */
const props = defineProps({
  name: {
    type: String,
    required: true,
    default: "",
  },
  color: {
    type: String,
    required: false,
    default: "#000",
  },
});

const symbolId = computed(() => `#icon-${props.name}`);
</script>

<style scoped lang="scss">
.svgSize {
  width: 100%;
  height: 100%;
}
</style>
