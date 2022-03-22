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
 * @Date: 2022-03-10 17:57:39
 * @LastEditTime: 2022-03-18 14:20:15
 */
// 使用 Mock
import Mock from 'mockjs';
import loginAPI from './login'

Mock.mock(/\/api\/login/, 'post', loginAPI.loginByUsername)

Mock.mock(/\/api\/logout/, 'get', loginAPI.logout)

Mock.mock(/\/api\/permission/, 'post', loginAPI.getPermission)

export default Mock