<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->

<template>
  <div class="frameC">
    <div class="frameC" v-if="showFrame">
      <iframe :src="props.url" class="frameC"></iframe>
    </div>
    <div class="frameC" v-show="!showFrame">
      <micro-app
        class="frameC"
        :baseroute="'/' + props.name"
        :name="props.name"
        :url="props.url"
        inline
        keep-alive
        fiber
        iframe
        disable-memory-router
      ></micro-app>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watchEffect } from "vue";
import { useRoute } from "vue-router";
import microApp from "@micro-zoe/micro-app";

const showFrame = ref(false);
const route: any = useRoute();
let props = reactive({
  url: "",
  plugin_type: "",
  name: "",
});
watchEffect(() => {
  if (route.meta.plugin_type) {
    // 如果是插件
    if (route.meta.subRoute) {
      // 插件的子路由
      let { plugin_type, url, subRoute, name } = route.meta;
      props.plugin_type = plugin_type;
      props.name = name;

      showFrame.value = props.plugin_type === "iframe" ? true : false;
      if (props.plugin_type === "micro-app") {
        props.url = window.location.origin + subRoute;
      } else {
        props.url = url;
      }
    } else {
      // 插件只有一级路由
    }
  }
});
</script>

<style lang="scss" scoped>
.frameC {
  width: 100%;
  height: 100%;
}
</style>
