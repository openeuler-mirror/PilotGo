import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store' //导入store (vuex)
import axios from 'axios'
import ElementUI from 'element-ui';  //导入element ui
import 'element-ui/lib/theme-chalk/index.css';
import * as echarts from 'echarts'; //echarts 5.0 导入方式
import './permission';


Vue.prototype.$http = axios  
Vue.prototype.$echarts = echarts

Vue.use(ElementUI);  
Vue.use(echarts);

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
