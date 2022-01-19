const cluster = {
    state: {
        selectIp: '',
    },
    mutations: {
        SET_SELECTIP(state, ip) {
            state.selectIp = ip;
        },
    },
}

export default cluster;