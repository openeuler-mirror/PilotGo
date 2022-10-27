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
 * LastEditTime: 2022-03-10 11:00:43
 * Description: redis初始化
 ******************************************************************************/
package redismanager

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

var (
	EnableRedis bool
	DialTimeout time.Duration
)

func RedisInit(redisConn, redisPwd string, defaultDB int, dialTimeout time.Duration, enableRedis bool) error {
	global.PILOTGO_REDIS = redis.NewClient(&redis.Options{
		Addr:     redisConn,
		Password: redisPwd,
		DB:       defaultDB,
	})
	// 使用超时上下文，验证redis
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), dialTimeout)
	defer cancelFunc()
	_, err := global.PILOTGO_REDIS.Ping(timeoutCtx).Result()
	if err != nil {
		return err
	}
	DialTimeout = dialTimeout
	EnableRedis = enableRedis
	return nil
}
