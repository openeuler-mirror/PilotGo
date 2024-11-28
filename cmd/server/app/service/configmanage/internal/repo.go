/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package internal

import "gitee.com/openeuler/PilotGo/pkg/utils/os/common"

type RepoFile struct {
	RepoSources []*common.RepoSource
}

type RepoConfig struct {
	Repos []*RepoFile
}

func (c *RepoConfig) Record() error {

	return nil
}

func (c *RepoConfig) Load() error {
	return nil
}

func (c *RepoConfig) Apply(uuid string) error {

	return nil
}
