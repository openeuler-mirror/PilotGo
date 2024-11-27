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
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
	"k8s.io/klog/v2"
)

var (
	EnableRedis bool
	DialTimeout time.Duration

	global_redis *redis.Client
)

func RedisInit(redisConn, redisPwd string, defaultDB int, dialTimeout time.Duration, enableRedis bool, stopCh <-chan struct{}, useTLS bool) error {
	var cfg *redis.Options
	if useTLS {
		cfg = &redis.Options{
			Addr:     redisConn,
			Password: redisPwd,
			DB:       defaultDB,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	} else {
		cfg = &redis.Options{
			Addr:     redisConn,
			Password: redisPwd,
			DB:       defaultDB,
		}
	}

	global_redis = redis.NewClient(cfg)
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), dialTimeout)
	// 使用超时上下文，验证redis
	go func() {
		<-stopCh
		global_redis.Close()
		cancelFunc()
		klog.Warning("global_redis success exit")

	}()
	defer cancelFunc()
	_, err := global_redis.Ping(timeoutCtx).Result()
	if err != nil {
		return err
	}
	DialTimeout = dialTimeout
	EnableRedis = enableRedis
	return nil
}

func Redis() *redis.Client {
	return global_redis
}
