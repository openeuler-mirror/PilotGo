/*
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: zhaozhenfang
 * @Date: 2022-03-18 14:07:35
 * @LastEditTime: 2022-03-18 14:17:53
 */
const userMap = {
  "code": 200,
  "msg": "成功",
  "data": {
      "id": 1,
      "username": "admin",
      "departName":"麒麟",
      "userType": 0,
      "departId": 1,
      "roleId": "1",
      "token": 'sadasdafaccadsf'
  }
}

export default {
  loginByUsername: config => {
      return userMap
  },
  logout: () => {
      return {
          "timestamp": "2021-07-27T09:45:48+0800",
          "code": "0",
          "message": "成功",
          "data": null
      }
  },
  getPermission: () => {
    return {
      "code": 200,
      "msg": "成功",
      "data": {
          "menu": [
              "overview","cluster","batch","usermanage","rolemanage","firewaall","log"
          ],
          "button":[
            "user_add","user_del","create_batch","rpm_install"
          ]
  }
    }
  }
}