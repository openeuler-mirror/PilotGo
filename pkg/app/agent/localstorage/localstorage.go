/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: yangzhao1
 * Date: 2022-04-13 15:11:03
 * LastEditTime: 2022-04-18 16:02:57
 * Description: provide agent log manager of pilotgo
 ******************************************************************************/

package localstorage

import (
	"encoding/json"
	"flag"
	"sync"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/sdk/logger"
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
			logger.Error("init local storage failed:%s", err.Error())
			return err
		}
	}

	if err := load(); err != nil {
		logger.Error("load local storage failed:%s", err.Error())
		return err
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
		return err
	}

	data := &localData{}
	if err := json.Unmarshal([]byte(s), data); err != nil {
		return err
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
		return err
	}

	err = utils.FileSaveString(LocalStorageFile, string(bs))
	if err != nil {
		return err
	}

	globalLock.Lock()
	globalLocalData = data
	globalLock.Unlock()
	return nil
}
