
<template>
  <div class="iframe-content" v-if="showFrame">
    <iframe :src="props.url" class="iframe"></iframe>
  </div>
  <div class="micro_content" v-show="!showFrame">
    <!-- 当开启fiber后，micro-app会降低子应用的优先级，通过异步执行子应用的js文件 -->
    <micro-app class="micro" :name="props.name" :url="microUrl" inline keep-alive fiber iframe></micro-app>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watchEffect } from "vue";
import { useRoute } from "vue-router";
import type { Extention } from '@/types/plugin';


const showFrame = ref(false);
const microUrl = ref('');
const route: any = useRoute();
let props = reactive({
  url: '',
  plugin_type: '',
  name: '',
  subMenus: [] as any


})

watchEffect(() => {
  props.url = route.meta.props.url;
  props.plugin_type = route.meta.props.plugin_type;
  props.name = route.name;
  props.subMenus = route.meta.subMenus.filter((item: Extention) => item.type === 'page');

  // 判断插件采用哪种方式展示
  showFrame.value = props.plugin_type === 'iframe' ? true : false;
  if (props.plugin_type === 'micro-app') {
    // 采用micro微前端
    microUrl.value = window.location.origin + '/plugin/atune'
  }

})
</script>

<style lang="scss" scoped>
.iframe-content,
.micro_content {
  box-sizing: border-box;
  height: 100%;
  width: 100%;

  .iframe,
  .micro {
    width: 100%;
    height: 100%;
    border: 0;
  }
}
</style>