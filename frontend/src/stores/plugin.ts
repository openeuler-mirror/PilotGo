import { defineStore } from 'pinia';
import type { Extention } from '@/types/plugin';
export const usePluginStore = defineStore('plugin', {
  state: () => {
    return {
      extention: [] as Extention[],
    };
  },
  actions: {
    addExtention(pluginExt: Extention[]) {
      // 添加扩展点
      Array.prototype.push.apply(this.extention, pluginExt);
    },
    delExtention(pluginExt: Extention[]) {
      // 删除扩展点
      let arr = [] as Extention[];
      pluginExt.map(item => {
        this.extention.map((extItem, index) => {
          if (extItem.name === item.name) {
            this.extention.splice(index, 1);
          }
        })
      })
      this.extention = arr;
    },
  },
  persist: true
});