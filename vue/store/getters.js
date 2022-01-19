const getters = {
    token: state => state.user.token,
    userName: state => state.user.username,
    roles: state => state.user.roles,
    getUserId: state => state.user.userId,
    activePanel: state => state.permissions.activePanel,
}

export default getters