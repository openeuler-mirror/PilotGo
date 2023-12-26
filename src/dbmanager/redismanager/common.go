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

func Set(key string, value interface{}) error {
	var ctx = context.Background()
	if EnableRedis {
		bytes, _ := json.Marshal(value)
		err := Redis().Set(ctx, key, string(bytes), 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func Get(key string, obj interface{}) (interface{}, error) {
	var ctx = context.Background()
	if EnableRedis {
		data, err := Redis().Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(data), obj)
		return obj, nil
	}
	return nil, fmt.Errorf("未启用Redis")
}

func Scan(key string) []string {
	var ctx = context.Background()
	keys := []string{}
	if EnableRedis {
		iterator := Redis().Scan(ctx, 0, key, 0).Iterator()
		for iterator.Next(ctx) {
			key := iterator.Val()
			keys = append(keys, key)
		}
		return keys
	}
	return []string{}
}

func Delete(key string) error {
	var ctx = context.Background()
	if EnableRedis {
		err := Redis().Del(ctx, key).Err()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
