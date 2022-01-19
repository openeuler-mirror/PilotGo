import Vue from 'vue'
import Vuex from 'vuex'
import user from './modules/user'
import cluster from './modules/cluster'
import permissions from './modules/permissions'
import getters from './getters'

Vue.use(Vuex);

const store =  new Vuex.Store({
    modules: {
      user,
      permissions,
      cluster,
    },
    getters,
})
export default store;