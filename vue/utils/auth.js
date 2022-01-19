import Cookies from 'js-cookie'

const TokenKey = 'Admin-Token'

const Username = 'Username'

const Roles = 'Roles'

const UserId = "userId"


const whileList = [
    "/login",
    "/401",
    "/404"
]

export function getToken() {
    return Cookies.get(TokenKey)
}

export function setToken(token) {
    if (token) {
        return Cookies.set(TokenKey, token)
    }
}

export function removeToken() {
    return Cookies.remove(TokenKey)
}

export function getUsername() {
    return Cookies.get(Username)
}

export function setUsername(username) {
    if (username) {
        return Cookies.set(Username, username)
    }
}

export function removeUsername() {
    return Cookies.remove(Username)
}

export function getRoles() {
    return Cookies.get(Roles)
}

export function removeRoles() {
    return Cookies.remove(Roles)
}
export function setRoles(roles) {
    if (roles) {
        return Cookies.set(Roles, roles)
    }
}

export function getUserId() {
    return Cookies.get(UserId)
}

export function setUserId(userId) {
    if (userId) {
        return Cookies.set(UserId, userId)
    }
}

export function removeUserId() {
    return Cookies.remove(UserId)
}

export function hasPermission(menus, to) {
    if (whileList.includes(to.path)) return true
    if (!to.meta) return true
    return Array.isArray(menus) && menus.includes(to.meta.panel)
}