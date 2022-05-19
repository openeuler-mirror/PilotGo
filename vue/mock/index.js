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
 * @LastEditTime: 2022-05-12 16:28:33
 */
// 使用 Mock
import Mock from 'mockjs';
import { cron_list } from './cluster'

//使用mockjs模拟数据
Mock.mock('/api/cron_list','get', {
  code: 200,
  cron_info: [
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
});
Mock.mock('/api/test', 'get', {
    "code":200,
    "data":
      {
        "firewall_config|1-50": " item:test \n",//随机生成日期时间
      }
});

Mock.mock('/macmanager/getIps', 'post', {
  "code":200,
  "data":
    [
      {
        "uuid": Mock.mock('@integer(1, 100)'),
        "departid": '1',
        "departname":'服务器',
        "ip|1": Mock.mock('@ip()'),
      },
      {
        "uuid": Mock.mock('@integer(1, 100)'),
        "departid": '2',
        "departname":'开源社区',
        "ip|1": Mock.mock('@ip()'),
      },
      {
        "uuid": Mock.mock('@integer(1, 100)'),
        "departid": '3',
        "departname":'服务器',
        "ip|1": Mock.mock('@ip()'),
      },
      {
        "uuid": Mock.mock('@integer(1, 100)'),
        "departid": '4',
        "departname":'开源社区',
        "ip|1": Mock.mock('@ip()'),
      }
    ]
    
});

Mock.mock('/config/allRepos', 'get', {
  "code":200,
  "data":
    {
      "firewall_config|1-50": " item:test \n",//随机生成日期时间
    }
});

export default Mock