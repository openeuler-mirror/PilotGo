package common

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
