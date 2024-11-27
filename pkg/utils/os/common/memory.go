/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 19:05:07 2023 +0800
 */
package common

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type MemoryConfig struct {
	MemTotal     int64
	MemFree      int64
	MemAvailable int64
	Buffers      int64
	Cached       int64
	SwapCached   int64
}

func (conf *MemoryConfig) String() string {
	b, err := json.Marshal(*conf)
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	return out.String()
}
