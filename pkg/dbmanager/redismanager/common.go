/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-10-26 09:05:39
 * LastEditTime: 2022-04-25 16:00:43
 * Description: redis客户端结构体
 ******************************************************************************/
package redismanager

import (
	"context"
	"encoding/json"
	"fmt"
)

func Set(key string, value interface{}) {
	var ctx = context.Background()
	if EnableRedis {
		bytes, _ := json.Marshal(value)
		err := Redis.Set(ctx, key, string(bytes), 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func Get(key string, obj interface{}) (interface{}, error) {
	var ctx = context.Background()
	if EnableRedis {
		data, err := Redis.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(data), obj)
		return obj, nil
	}
	return nil, fmt.Errorf("未启用Redis")
}
