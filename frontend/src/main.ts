/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import "./styles/main.scss";

import { createApp } from "vue";
import pinia from "@/stores";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

import App from "./App.vue";
import router from "./router";
import microApp from "@micro-zoe/micro-app";

import "amfe-flexible";

// 监听窗口在变化时重新设置跟文件大小
window.onresize = function () {
  console.log(
    "当前字体大小：",
    document.documentElement.clientWidth,
    document.documentElement.style.fontSize,
    window.devicePixelRatio
  );
};
export const app = createApp(App);
app.use(pinia);
app.use(router);
app.use(ElementPlus);

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

app.mount("#app");
microApp.start();
