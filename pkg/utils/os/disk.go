package os

import (
	"github.com/shirou/gopsutil/disk"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
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
		logger.Error("get Partitions failed, err:%v\n", err)
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
