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
