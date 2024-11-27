/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 19:05:07 2023 +0800
 */
package common

// 形如	openssl-1:1.1.1f-4.oe1.x86_64
//
//	OS
//	openssl=1:1.1.1f-4.oe1
type RpmSrc struct {
	Name     string
	Repo     string
	Provides string
}

type RpmInfo struct {
	Name         string
	Version      string
	Release      string
	Architecture string
	InstallDate  string
	Size         string
	License      string
	Signature    string
	Packager     string
	Vendor       string
	URL          string
	Summary      string
}

type RepoSource struct {
	File           string
	ID             string
	Name           string
	MirrorList     string
	BaseURL        string
	MetaLink       string
	MetadataExpire string
	GPGCheck       int
	GPGKey         string
	Enabled        int
}
