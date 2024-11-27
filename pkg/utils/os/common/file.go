/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: wubijie <wubijie@kylinos.cn>
 * Date: Mon Apr 24 14:39:39 2023 +0800
 */
package common

type UpdateFile struct {
	Path            string `json:"path"`
	Name            string `json:"name"`
	Text            string `json:"text"`
	FileLastVersion string `json:"filelastversion"`
}
