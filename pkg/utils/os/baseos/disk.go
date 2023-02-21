package baseos

import (
	"fmt"
	"math"
	"strconv"

	"github.com/shirou/gopsutil/disk"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type DiskIOInfo struct {
	PartitionName string
	Label         string
	ReadCount     uint64
	WriteCount    uint64
	ReadBytes     uint64
	WriteBytes    uint64
	IOTime        uint64
}

type DiskUsageINfo struct {
	Device      string `json:"device"`
	Path        string `json:"path"`
	Fstype      string `json:"fstype"`
	Total       string `json:"total"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedPercent"`
}

const (
	CapacitySize = 1
)

// 获取磁盘的使用情况
func (b *BaseOS) GetDiskUsageInfo() []DiskUsageINfo {
	diskusage := make([]DiskUsageINfo, 0)
	parts, err := disk.Partitions(false)
	if err != nil {
		logger.Error("get Partitions failed, err:%v\n", err.Error())
		return nil
	}
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		device := part.Device
		path := diskInfo.Path
		fstype := diskInfo.Fstype
		t := diskInfo.Total / (1024 * 1024 * 1024)
		var total string
		if t > CapacitySize { //判断磁盘容量是否大于1G
			total = strconv.FormatUint(diskInfo.Total/(1024*1024*1024), 10) + "G"
		} else {
			total = strconv.FormatUint(diskInfo.Total/(1024*1024), 10) + "M"
		}
		u := diskInfo.Used / (1024 * 1024 * 1024)
		var used string
		if u > CapacitySize { //判断磁盘已用容量是否大于1G
			used = strconv.FormatUint(diskInfo.Used/(1024*1024*1024), 10) + "G"
		} else {
			used = strconv.FormatUint(diskInfo.Used/(1024*1024), 10) + "M"
		}
		usedPercent := int(math.Floor(diskInfo.UsedPercent))
		tmp := DiskUsageINfo{
			Device:      device,
			Path:        path,
			Fstype:      fstype,
			Total:       total,
			Used:        used,
			UsedPercent: strconv.Itoa(usedPercent) + "%",
		}
		diskusage = append(diskusage, tmp)
	}
	return diskusage
}

// 获取磁盘的IO信息
func (b *BaseOS) GetDiskInfo() []DiskIOInfo {

	diskinfo := make([]DiskIOInfo, 0)
	ioStat, _ := disk.IOCounters()
	for k, v := range ioStat {
		label := v.Label
		readCount := v.ReadCount
		writeCount := v.WriteCount
		readBytes := v.ReadBytes
		writeBytes := v.WriteBytes
		ioTime := v.IoTime

		tmp := DiskIOInfo{
			PartitionName: k,
			Label:         label,
			ReadCount:     readCount,
			WriteCount:    writeCount,
			ReadBytes:     readBytes,
			WriteBytes:    writeBytes,
			IOTime:        ioTime,
		}
		diskinfo = append(diskinfo, tmp)
	}
	return diskinfo
}

/*
挂载磁盘
1.创建挂载磁盘的目录
2.挂载磁盘
*/
func (b *BaseOS) CreateDiskPath(mountpath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("mkdir %s", mountpath))
	if err != nil {
		logger.Error("创建挂载目录失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("创建挂载目录成功!%s", tmp)
	return tmp
}

func (b *BaseOS) DiskMount(sourceDisk, destPath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("mount %s %s", sourceDisk, destPath))
	if err != nil {
		logger.Error("挂载磁盘失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("挂载磁盘成功!%s", tmp)
	return tmp
}

// 卸载磁盘
func (b *BaseOS) DiskUMount(diskPath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("umount %s", diskPath))
	if err != nil {
		logger.Error("卸载磁盘失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("卸载磁盘成功!%s", tmp)
	return tmp
}

// 磁盘格式化
func (b *BaseOS) DiskFormat(fileType, diskPath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("mkfs.%s %s", fileType, diskPath))
	if err != nil {
		logger.Error("格式化磁盘失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("格式化磁盘成功!%s", tmp)
	return tmp
}
