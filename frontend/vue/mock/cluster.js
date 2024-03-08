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
 * @Date: 2022-04-29 11:02:09
 * @LastEditTime: 2022-04-29 11:18:53
 */
import Mock from 'mockjs';
export let cron_list = {
  code: 200,
  corn_info: [
    {
      id: 1,
      name: 'test1',
      cron: '0 * * * *',
      createdAt: Mock.mock('@date()'),
      updatedAt: Mock.mock('@date()'),
      status: 1,
      description: '测似例1'
    },
    {
      id: 2,
      name: 'test2',
      cron: '0 0 * * *',
      createdAt: Mock.mock('@date()'),
      updatedAt: Mock.mock('@date()'),
      status: 0,
      description: '测似例2'
    },

  ]
}