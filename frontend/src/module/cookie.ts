/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
 */
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
