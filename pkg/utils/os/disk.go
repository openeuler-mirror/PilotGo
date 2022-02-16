package os

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
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
	Device      string  `json:"device"`
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// 获取磁盘的使用情况
func GetDiskUsageInfo() []DiskUsageINfo {
	diskusage := make([]DiskUsageINfo, 0)
	parts, err := disk.Partitions(true)
	if err != nil {
		logger.Error("get Partitions failed, err:%v\n", err.Error())
		return nil
	}
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		device := part.Device
		path := diskInfo.Path
		fstype := diskInfo.Fstype
		total := diskInfo.Total
		free := diskInfo.Free
		used := diskInfo.Used
		usedPercent := diskInfo.UsedPercent
		tmp := DiskUsageINfo{
			Device:      device,
			Path:        path,
			Fstype:      fstype,
			Total:       total,
			Free:        free,
			Used:        used,
			UsedPercent: usedPercent,
		}
		diskusage = append(diskusage, tmp)
	}
	return diskusage
}

// 获取磁盘的IO信息
func GetDiskInfo() []DiskIOInfo {

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

/*挂载磁盘
1.创建挂载磁盘的目录
2.挂载磁盘*/
func CreateDiskPath(mountpath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("mkdir %s", mountpath))
	if err != nil {
		logger.Error("创建挂载目录失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("创建挂载目录成功!%s", tmp)
	return tmp
}
func DiskMount(sourceDisk, destPath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("mount %s %s", sourceDisk, destPath))
	if err != nil {
		logger.Error("挂载磁盘失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("挂载磁盘成功!%s", tmp)
	return tmp
}

// 卸载磁盘
func DiskUMount(diskPath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("umount %s", diskPath))
	if err != nil {
		logger.Error("卸载磁盘失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("卸载磁盘成功!%s", tmp)
	return tmp
}

// 磁盘格式化
func DiskFormat(fileType, diskPath string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("mkfs.%s %s", fileType, diskPath))
	if err != nil {
		logger.Error("格式化磁盘失败!%s", err.Error())
		return err.Error()
	}
	logger.Info("格式化磁盘成功!%s", tmp)
	return tmp
}
