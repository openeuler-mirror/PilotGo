import { loginByEmail, logout } from '@/request/user'
import { getToken, setToken, removeToken, getUsername, setUsername, removeUsername, 
    getRoles, setRoles, removeRoles, getUserId, removeUserId, setUserId } from '@/utils/auth'

const user = {
    state: {
        token: getToken(),
        username: getUsername(),
        roles: getRoles() ? JSON.parse(getRoles()) : [],
        userId: getUserId(),
    },
    mutations: {
        SET_TOKEN: (state, token) => {
            state.token = token
        },
        SET_NAME: (state, name) => {
            state.username = name
        },
        SET_ROLES: (state, roles) => {
            state.roles = roles
        },
        SET_USERID: (state, userId) => {
            state.userId = userId
        },
    },
    actions: {
        loginByEmail2({ commit }, userInfo){
            commit('SET_TOKEN', "saajhdjshdjsad12")
            commit('SET_NAME', userInfo.email)
            commit('SET_USERID', "111")
            setToken("saajhdjshdjsad12")
            setUsername(userInfo.email)
            setUserId("111")
        },
        logOut({ commit, dispatch }) {
            commit('SET_TOKEN', '')
            commit('SET_ROLES', [])
            commit('SET_MENUS', [])
            commit('SET_NAME', '')
            removeUsername();
            removeToken();
            removeUserId();
            localStorage.clear()
        },
        loginByEmail({ commit }, userInfo) {
            const username = userInfo.username.trim()
            return new Promise((resolve, reject) => {
                loginByEmail({'email':username, 'password':userInfo.password}).then(response => {
                    const res = response.data;
                    if (res.code != "200") {
                        reject(res)
                    } else {
                        commit('SET_TOKEN', res.data.token)
                        commit('SET_NAME', username)
                        
                        setToken(res.data.token)
                        setUsername(username)
                        resolve()
                    }
                }).catch(error => {
                    reject(error)
                })
            })
        },
        logOut1({ commit }) {
            return new Promise((resolve, reject) => {
                logout().then(() => {
                    commit('SET_TOKEN', '')
                    commit('SET_ROLES', [])
                    commit('SET_MENUS', [])
                    commit('SET_NAME', '')
                    removeRoles();
                    removeUsername();
                    removeToken();
                    removeUserId();
                    localStorage.clear()
                    resolve()
                }).catch(error => {
                    reject(error)
                })
            })
        },
    }
}

export default user