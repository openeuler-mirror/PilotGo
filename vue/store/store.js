import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
    //所有的数据都放在state中
    state: {
        menuIndex: '1',   //导航菜单选中的菜单
        selectedClusterIp: '',
        cockpitPluginServer:'',
        currentUser: JSON.parse(window.sessionStorage.getItem("user")),
    },

    //操作数据，唯一的通道是mutations
    mutations: {
        menuKeySelect: (state, key) => {
            // console.log('key', key)
            state.menuIndex = key
        },
        mutateSelectedClusterIp: (state, ip) => {
            state.selectedClusterIp = ip
        },
        mutateCockpitPluginServer: (state, ip) => {
            state.cockpitPluginServer = ip
        },
    },

    //actions,可以来做异步操作，然后提交给mutations，而后再对state(数据)进行操作
    actions: {

    }
})
