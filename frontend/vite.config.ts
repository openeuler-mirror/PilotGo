/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { createSvgIconsPlugin } from "vite-plugin-svg-icons";

import postCssPxToRem from "postcss-pxtorem";
import { resolve } from "node:path";

export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => /^micro-app/.test(tag),
        },
      },
    }),
    createSvgIconsPlugin({
      iconDirs: [resolve(resolve("src", "assets/icons"))],
      symbolId: "icon-[dir]-[name]",
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 8080,
    https: false,
    cors: true,
    proxy: {
      "/api/v1": {
        target: "https://10.41.107.29:8888", //10.41.107.29:8888 \  10.44.43.181:8888 \ 10.41.161.101:8888"
        secure: false, // 如果是https接口，需要配置这个参数
        changeOrigin: true, // 如果接口跨域，需要进行这个参数配置
        // rewrite: path => path.replace(/^\/demo/, '/demo')
      },
      "/plugin": {
        target: "http://10.44.43.181:8090",
        secure: false, // 如果是https接口，需要配置这个参数
        changeOrigin: true, // 如果接口跨域，需要进行这个参数配置
        // rewrite: path => path.replace(/^\/demo/, '/demo')
      },
    },
  },
  css: {
    postcss: {
      plugins: [
        postCssPxToRem({
          rootValue: 192,
          propList: ["*", "!border"],
          // exclude: /node_modules/i,
        }),
      ],
    },
  },
});
