import { ref, triggerRef } from "vue";
import { ElMessage } from 'element-plus';
import { RespCodeOK } from "@/request/request";
import { getPlugins } from "@/request/plugin";
import type { Extention } from '@/types/plugin';
// const router = useRouter();
export const iframeComponents = ref<any>([])

/* 
* @func 更新插件路由和扩展点信息
*/
export function updatePlugins() {
  iframeComponents.value = [];
  getPlugins().then((res: any) => {
    if (res.code === RespCodeOK) {
      let iframes: any = []
      res.data.forEach((item: any, index: number) => {
        if (item.enabled === 0 || item.status === false) {
          // 0:禁用，1：启用
          return;
        }
        // 创建组件
        let iframeObj: any = {
          path: '/plugin-' + item.name,
          name: 'plugin-' + item.name,
          url: item.url,
          plugin_type: item.plugin_type,
          subMenus: []
        }
        // 添加扩展点信息
        if (item.extentions && item.extentions.length > 0) {
          // 遍历page页面添加subMenu
          let subMenus = [] as any;
          item.extentions.filter((extItem: Extention) => extItem.type === 'page')
            .forEach((pageItem: Extention) => {
              let subMenuObj = {
                path: '/plugin/' + item.name,
                subRoute: pageItem.url,
                title: pageItem.name,
                hidden: false,
                panel: '/plugin/' + item.name,
                icon: '',
                subMenus: null,
                isPlug: true,
              }
              subMenus.push(subMenuObj);
            })
          iframeObj.subMenus = subMenus;
        }
        iframes.push(iframeObj);
      })
      iframeComponents.value = iframes;
    } else {
      ElMessage.error("查询插件列表错误：" + res.msg);
    }
  }).catch((err) => {
    console.log('处理数据流程出错：', err)
  })
}

