package baseos

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils"
	"gitee.com/PilotGo/PilotGo/utils/os/common"
	"github.com/shirou/gopsutil/disk"
)

const (
	CapacitySize = 1
)

// 获取磁盘的使用情况
func (b *BaseOS) GetDiskUsageInfo() ([]common.DiskUsageINfo, error) {
	diskusage := make([]common.DiskUsageINfo, 0)
	parts, err := disk.Partitions(false)
	if err != nil {
		logger.Error("get Partitions failed, err: %s", err.Error())
		return nil, err
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
		tmp := common.DiskUsageINfo{
			Device:      device,
			Path:        path,
			Fstype:      fstype,
			Total:       total,
			Used:        used,
			UsedPercent: strconv.Itoa(usedPercent) + "%",
		}
		diskusage = append(diskusage, tmp)
	}
	return diskusage, nil
}

// 获取磁盘的IO信息
func (b *BaseOS) GetDiskInfo() ([]common.DiskIOInfo, error) {

	diskinfo := make([]common.DiskIOInfo, 0)
	ioStat, err := disk.IOCounters()
	if err != nil {
		logger.Error("get disk IO failed, err: %s", err.Error())
		return nil, err
	}
	for k, v := range ioStat {
		label := v.Label
		readCount := v.ReadCount
		writeCount := v.WriteCount
		readBytes := v.ReadBytes
		writeBytes := v.WriteBytes
		ioTime := v.IoTime

		tmp := common.DiskIOInfo{
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
	return diskinfo, nil
}

/*
挂载磁盘
1.创建挂载磁盘的目录
2.挂载磁盘
*/

func (b *BaseOS) DiskMount(sourceDisk, mountPath string) (string, error) {
	// 创建挂载目录
	err := os.MkdirAll(mountPath, 0644)
	if err != nil {
		logger.Error("failed to create a mounted directory: %s", err.Error())
		return err.Error(), fmt.Errorf("failed to create a mounted directory: %s", err.Error())
	}
	logger.Info("successfully created a mounted directory: %s", mountPath)

	exitc, stdo, stde, err := utils.RunCommand(fmt.Sprintf("mount %s %s", sourceDisk, mountPath))
	fmt.Printf("[diskmount]%v, %v, %v, %v\n", exitc, stdo, stde, err)
	if exitc == 0 && stdo == "" && stde == "" && err == nil {
		logger.Info("successfully mounted disk: %s", stdo)
		return stdo, nil
	}
	logger.Error("failed to mount disk")
	return "failed to mount disk", fmt.Errorf("failed to mount disk: %d, %s, %s, %v", exitc, stdo, stde, err)
}

// 卸载磁盘
func (b *BaseOS) DiskUMount(diskPath string) (string, error) {
	exitc, stdo, stde, err := utils.RunCommand(fmt.Sprintf("umount %s", diskPath))
	if exitc == 0 && stdo == "" && stde == "" && err == nil {
		logger.Info("successfully unmounted the disk: %s", stdo)
		return stdo, nil
	}
	logger.Error("failed to unmount the disk")
	return "failed to unmount the disk", fmt.Errorf("failed to unmount the disk: %d, %s, %s, %v", exitc, stdo, stde, err)
}

// 磁盘格式化
func (b *BaseOS) DiskFormat(fileType, diskPath string) (string, error) {
	exitc, stdo, stde, err := utils.RunCommand(fmt.Sprintf("mkfs.%s -F %s", fileType, diskPath))
	if exitc == 0 && stdo != "" && stde != "" && err == nil {
		logger.Info("successfully formatted the disk: %s", stdo)
		return stdo, nil
	}
	logger.Error("failed to format the disk: %d, %s, %s, %v", exitc, stdo, stde, err)
	return "", fmt.Errorf("failed to format the disk: %d, %s, %s, %v", exitc, stdo, stde, err)
}
