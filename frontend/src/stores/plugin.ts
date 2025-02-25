/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Feb 20 11:33:42 2024 +0800
 */
import { defineStore } from "pinia";
import type { Extention } from "@/types/plugin";
export const usePluginStore = defineStore("plugin", {
  state: () => {
    return {
      extention: [] as any[],
    };
  },
  actions: {
    addExtention(pluginExt: Extention[]) {
      // 添加扩展点
      const map = new Map();
      [...this.extention, ...pluginExt].forEach((item) => {
        if (!map.has(item.name)) {
          map.set(item.name, item);
        }
      });
      let result = Array.from(map.values());

      this.extention = result;
    },
    delExtention(pluginExt: Extention[]) {
      // 删除扩展点
      let arr = [] as Extention[];
      pluginExt.map((item) => {
        this.extention.map((extItem, index) => {
          if (extItem.name === item.name) {
            this.extention.splice(index, 1);
          }
        });
      });
      this.extention = arr;
    },
  },
  persist: true,
});
