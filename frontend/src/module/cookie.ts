import Cookies from 'js-cookie';

// cookie定义
// TODO: use simple token name
// const CookieAuthToken = "token"
const CookieAuthToken = "Admin-Token"

export function setToken(token: string) {
    Cookies.set(CookieAuthToken, token)
}
export function getToken() {
  return Cookies.get(CookieAuthToken)
}

export function removeToken() {
    Cookies.remove(CookieAuthToken)
}
