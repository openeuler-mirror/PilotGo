<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div ref="chartContainer" id="chartContainer" class="echarts-container"></div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount, watchEffect, nextTick } from "vue";
import * as echarts from "echarts";
type EChartsOption = echarts.EChartsOption;
const props = defineProps({
  options: Object as () => any,
});

const chartContainer = ref<HTMLDivElement | null>(null);
let chartInstance: echarts.ECharts | null = null;

// 初始化图表
const initChart = () => {
  if (chartContainer.value) {
    chartInstance = echarts.init(chartContainer.value);
  }
};
watchEffect(() => {
  if (props.options) {
    nextTick(() => {
      chartInstance!.setOption(props.options);
    });
  }
});

// 组件挂载时初始化图表
onMounted(() => {
  initChart();
});

// 组件卸载时销毁图表实例
onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.dispose();
  }
});
</script>

<style scoped lang="scss">
.echarts-container {
  width: 100%;
  height: 400px;
}
</style>
