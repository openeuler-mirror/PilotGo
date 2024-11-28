/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package localstorage

import (
	"encoding/json"
	"flag"
	"sync"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"github.com/google/uuid"
)

type localData struct {
	AgentUUID string `json:"agent_uuid"`
}

var globalLocalData *localData
var globalLock sync.Mutex
var LocalStorageFile string

// init local storage, if file not found, then init new one
func Init() error {
	flag.StringVar(&LocalStorageFile, "uuid", "./.pilotgo-agent.data", "pilotgo-agent uuid data")
	flag.Parse()
	if fok, _ := utils.IsFileExist(LocalStorageFile); !fok {
		if err := reset(); err != nil {
			return errors.WithMessage(err, "init local storage failed")
		}
	}

	if err := load(); err != nil {
		return errors.WithMessage(err, "load local storage failed")
	}

	return nil
}

func AgentUUID() string {
	globalLock.Lock()
	uuid := globalLocalData.AgentUUID
	globalLock.Unlock()
	return uuid
}

// load local data from file
func load() error {
	s, err := utils.FileReadString(LocalStorageFile)
	if err != nil {
		return errors.Wrap(err, "failed to read local storage file")
	}

	data := &localData{}
	if err := json.Unmarshal([]byte(s), data); err != nil {
		return errors.Wrap(err, "failed to unmarshal local storage data")
	}

	globalLock.Lock()
	globalLocalData = data
	globalLock.Unlock()

	return nil
}

// init new local data and save to file
func reset() error {
	data := &localData{
		AgentUUID: uuid.New().String(),
	}

	bs, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal local data")
	}

	err = utils.FileSaveString(LocalStorageFile, string(bs))
	if err != nil {
		return errors.Wrap(err, "failed to save local storage file")
	}

	globalLock.Lock()
	globalLocalData = data
	globalLock.Unlock()
	return nil
}
