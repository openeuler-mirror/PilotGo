// store数据持久化
import { createPinia } from 'pinia'
import piniaPersisted from 'pinia-plugin-persistedstate'
const pinia = createPinia();
pinia.use(piniaPersisted);

// 全局数据入口,persist无效
/* import { usePluginStore } from './plugin';
export interface AppStore {
  usePluginStore: ReturnType<typeof usePluginStore>,
}
const appStore: AppStore = {} as AppStore;

// 注册app状态库
export const registerStore = () => {
  appStore.usePluginStore = usePluginStore()
} */

export default pinia;