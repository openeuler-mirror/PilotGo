/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Wed Jan 24 14:54:27 2024 +0800
 */
// store数据持久化
import { createPinia } from 'pinia'
import piniaPersisted from 'pinia-plugin-persistedstate'
const pinia = createPinia();
pinia.use(piniaPersisted);

// 全局数据入口,persist无效
/ * import { usePluginStore } from './plugin';
export interface AppStore {
  usePluginStore: ReturnType<typeof usePluginStore>,
}
const appStore: AppStore = {} as AppStore;

// 注册app状态库
export const registerStore = () => {
  appStore.usePluginStore = usePluginStore()
} */

export default pinia;