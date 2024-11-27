/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Wed Apr 27 00:34:37 2022 +0800
 */
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
