/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Dec 09 13:56:05 2024 +0800
 */
package client

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	client *clientv3.Client
	ctx    context.Context
}

// NewClient creates a new etcd client
func NewClient(endpoints []string, dialTimeout time.Duration) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: cli,
		ctx:    context.Background(),
	}, nil
}

// Close closes the client
func (c *Client) Close() error {
	return c.client.Close()
}

// Get gets the value for a key
func (c *Client) Get(key string) (*clientv3.GetResponse, error) {
	return c.client.Get(c.ctx, key)
}

// Put puts a key-value pair
func (c *Client) Put(key, value string) (*clientv3.PutResponse, error) {
	return c.client.Put(c.ctx, key, value)
}

// Delete deletes a key
func (c *Client) Delete(key string) (*clientv3.DeleteResponse, error) {
	return c.client.Delete(c.ctx, key)
}

// Watch watches for changes on a key
func (c *Client) Watch(key string) clientv3.WatchChan {
	return c.client.Watch(c.ctx, key)
}
