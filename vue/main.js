import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store/store' //导入store (vuex)
import axios from 'axios'
import ElementUI from 'element-ui';  //导入element ui
import 'element-ui/lib/theme-chalk/index.css';
import * as echarts from 'echarts'; //echarts 5.0 导入方式


Vue.prototype.$http = axios       //注册axios
Vue.prototype.$echarts = echarts

Vue.use(ElementUI);  //Vue引用element ui
Vue.use(echarts);

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,   //在new Vue 中使用
  components: { App },
  template: '<App/>'
})
