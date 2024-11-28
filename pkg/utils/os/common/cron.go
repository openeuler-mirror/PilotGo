/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 19:05:07 2023 +0800
 */
package common

import (
	"fmt"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	cron "github.com/robfig/cron/v3"
)

var Cron *Crontab

// crontab manager
type Crontab struct {
	Inner    *cron.Cron
	EntryIDs map[int]cron.EntryID
	Mutex    sync.Mutex
}

// new crontab
func NewCrontab() *Crontab {
	return &Crontab{
		Inner:    cron.New(cron.WithSeconds()),
		EntryIDs: make(map[int]cron.EntryID),
	}
}

// IDs ...
func (c *Crontab) IDs() []int {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	validIDs := make([]int, 0, len(c.EntryIDs))
	invalidIDs := make([]int, 0)
	for sid, eid := range c.EntryIDs {
		if e := c.Inner.Entry(eid); e.ID != eid {
			invalidIDs = append(invalidIDs, sid)
			continue
		}
		validIDs = append(validIDs, sid)
	}
	for _, id := range invalidIDs {
		delete(c.EntryIDs, id)
	}
	return validIDs
}

// start the crontab engine
func (c *Crontab) Start() {
	c.Inner.Start()
}

// stop the crontab engine
func (c *Crontab) Stop() {
	c.Inner.Stop()
}

// DeleteByID remove one crontab task
func (c *Crontab) DeleteByID(id int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	eid, ok := c.EntryIDs[id]
	if !ok {
		return fmt.Errorf("crontab id isn't exist")
	}
	c.Inner.Remove(eid)
	delete(c.EntryIDs, id)
	return nil
}

// AddByFunc add function as crontab task
func (c *Crontab) AddByFunc(id int, spec string, f func()) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, ok := c.EntryIDs[id]; ok {
		return fmt.Errorf("crontab id exists")
	}
	eid, err := c.Inner.AddFunc(spec, f)
	if err != nil {
		return err
	}
	c.EntryIDs[id] = eid
	return nil
}

// 创建客户端实例
func CronInit() error {
	crontab := NewCrontab()
	Cron = crontab
	return nil
}

// 开启任务
func CronStart(id int, spec string, command string) error {

	// eg.hello world
	// i := 0
	// taskFunc := func() {
	// 	i++
	// 	fmt.Println("hello world", i)
	// }

	// 添加函数作为定时任务
	taskFunc := func() {
		utils.RunCommand(command)
	}
	if err := Cron.AddByFunc(id, spec, taskFunc); err != nil {
		return fmt.Errorf("error to add crontab task: %s", err)
	}
	Cron.Start()
	time.Sleep(time.Duration(time.Millisecond * 300))
	return nil
}

// 暂停任务
func StopAndDel(id int) error {
	if err := Cron.DeleteByID(id); err != nil {
		return err
	}
	return nil
}
