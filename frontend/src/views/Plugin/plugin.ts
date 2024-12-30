/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
import { ref, triggerRef } from "vue";
import { ElMessage } from "element-plus";
import { RespCodeOK } from "@/request/request";
import { getPlugins } from "@/request/plugin";
import type { Extention } from "@/types/plugin";
// const router = useRouter();
export const iframeComponents = ref<any>([]);

/*
 * @func 更新插件路由和扩展点信息
 */
export function updatePlugins() {
  iframeComponents.value = [];
  getPlugins()
    .then((res: any) => {
      if (res.code === RespCodeOK) {
        let iframes: any = [];
        res.data.forEach((item: any, index: number) => {
          if (item.enabled === 0 || item.status === false) {
            // 0:禁用，1：启用
            return;
          }
          // 创建组件
          let iframeObj: any = {
            path: "/plugin-" + item.name,
            name: "plugin-" + item.name,
            menuName: item.menuName,
            url: item.url,
            plugin_type: item.plugin_type,
            icon: item.icon || "Menu",
          };

          // 筛选插件所有的page页面
          if (item.extentions && item.extentions.length > 0) {
            iframeObj.subMenus = item.extentions.filter((extItem: Extention) => extItem.type === "page");
          }

          iframes.push(iframeObj);
        });
        iframeComponents.value = iframes;
      } else {
        ElMessage.error("查询插件列表错误：" + res.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("处理数据流程出错：", err);
    });
}
