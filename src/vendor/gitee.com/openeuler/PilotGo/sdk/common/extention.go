package common

const (
	// 单机/多机操作扩展，增加选择单机/多机时对机器的操作功能
	ExtentionMachine = "machine"
	// 批处理扩展，增加选择批次时对批次的操作
	ExtentionBatch = "batch"
)

type Extention struct {
	PluginName string `json:"plugin_name"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	URL        string `json:"url"`
}
