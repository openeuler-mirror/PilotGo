/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package configfile

import "gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"

type ConfigFile = dao.ConfigFile

func AddConfigFile(cf ConfigFile) error {
	return dao.AddConfigFile(cf)
}
