import { ref, onMounted, watchEffect, shallowRef } from "vue";
import { ElMessage } from 'element-plus';

import { RespCodeOK } from "@/request/request";

import { getPlugins } from "@/request/plugin";
import PluginFrame from "@/views/Plugin/PluginFrame.vue";



// const router = useRouter();

export const iframeComponents = ref<any>([])

export function updatePlugins() {
    getPlugins().then((res: any) => {
        if (res.code === RespCodeOK) {
            let iframes: any = []
            res.data.forEach((item: any, index: number) => {
                if (item.enabled === 0) {

                    // 0:禁用，1：启用
                    return;
                }
                // p.push({
                //     path: '/plugin' + index,
                //     name: 'Plugin' + index,
                //     iframeComponent: '',
                //     meta: {
                //         title: 'plugin', header_title: item.name, panel: "plugin" + index, icon_class: 'el-icon-s-ticket', url: item.url,
                //         breadcrumb: [
                //             { name: item.name },
                //         ],
                //     }
                // })

                // 创建组件
                let iframeObj = {
                    path: '/plugin-' + item.name,
                    name: 'plugin-' + item.name,
                    // src: '/plugin/' + item.name,
                    // component: shallowRef(PluginFrame), // 组件文件的引用
                    url: item.url,
                    plugin_type: item.plugin_type
                }
                iframes.push(iframeObj);
            })
            iframeComponents.value = iframes;
        } else {
            ElMessage.error("查询插件列表错误：" + res.msg);
        }
    }).catch((err) => {
        ElMessage.error("查询插件列表错误：" + err.msg);
    })
}
