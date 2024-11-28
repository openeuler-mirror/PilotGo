/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
 */
package common

type OSOperator interface {
	SystemOperator
	CpuOperator
	MemoryOperator
	SysctlOperator
	DateTimeOperator
	ServiceOperator
	DiskOperator
	UserOperator
	FirewallOperator
	NetworkOperator
	PackageOperator
}
type SystemOperator interface {
	GetHostInfo() (*SystemInfo, error)
}

type CpuOperator interface {
	GetCPUInfo() (*CPUInfo, error)
}

type MemoryOperator interface {
	GetMemoryConfig() (*MemoryConfig, error)
}

type SysctlOperator interface {
	GetSysctlConfig() (map[string]string, error)
	TempModifyPar(string) (string, error)
	GetVarNameValue(string) (string, error)
}

type DateTimeOperator interface {
	GetTime() (string, error)
}

type ServiceOperator interface {
	GetServiceList() ([]ListService, error)
	GetService(service string) (*ServiceInfo, error)
	RestartService(string) error
	StartService(string) error
	StopService(string) error
}

type DiskOperator interface {
	GetDiskUsageInfo() ([]DiskUsageINfo, error)
	GetDiskInfo() ([]DiskIOInfo, error)
	DiskMount(string, string) (string, error)
	DiskUMount(string) (string, error)
	DiskFormat(string, string) (string, error)
}
type UserOperator interface {
	GetCurrentUserInfo() CurrentUser
	GetAllUserInfo() ([]AllUserInfo, error)
	AddLinuxUser(string, string) error
	DelUser(string) (string, error)
	ChangePermission(string, string) (string, error)
	ChangeFileOwner(string, string) (string, error)
}

type FirewallOperator interface {
	FirewalldSetDefaultZone(string) (string, error)
	FirewalldZoneConfig(string) (*FirewalldCMDList, error)
	FirewalldServiceAdd(string, string) error
	FirewalldServiceRemove(string, string) error
	FirewalldSourceAdd(string, string) error
	FirewalldSourceRemove(string, string) error
	Config() (FireWalldConfig, error)
	Restart() bool
	Stop() bool
	AddZonePort(string, string, string) (string, error)
	DelZonePort(string, string, string) (string, error)
}

type NetworkOperator interface {
	GetHostIp() (string, error)
	GetTCP() ([]NetConnect, error)
	GetUDP() ([]NetConnect, error)
	GetIOCounter() ([]IOCnt, error)
	GetNICConfig() ([]NetInterfaceCard, error)
	ConfigNetworkConnect() ([]map[string]string, error)
	GetNetworkConnInfo() (*NetworkConfig, error)
	GetNICName() (string, error)
	RestartNetwork(string) error
}

type PackageOperator interface {
	InstallRpm(string) error
	RemoveRpm(string) error
	GetAllRpm() ([]string, error)
	GetRpmSource(string) ([]RpmSrc, error)
	GetRpmInfo(string) (*RpmInfo, error)
	GetRepoSource() ([]*RepoSource, error)
}
