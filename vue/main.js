/*
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: zhaozhenfang
 * @Date: 2022-02-12 16:12:15
 * @LastEditTime: 2022-05-31 18:12:27
 */
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store' //导入store (vuex)
import axios from 'axios'
import ElementUI from 'element-ui';  //导入element ui
import CodeDiff from 'v-code-diff';
import 'element-ui/lib/theme-chalk/index.css';
import * as echarts from 'echarts'; //echarts 5.0 导入方式
import './permission';
import './styles/index.scss'
import './iconfont/iconfont.js'
import './iconfont/iconfont.css'
import './mock/index.js' //引入mockjs,上线后注掉




Vue.prototype.$http = axios  
Vue.prototype.$echarts = echarts

Vue.use(ElementUI);  
Vue.use(echarts);
Vue.use(CodeDiff);

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
