package common

const (
	// 单机/多机操作扩展，增加选择单机/多机时对机器的操作功能
	ExtentionMachine = "machine"
	// 批处理扩展，增加选择批次时对批次的操作
	ExtentionBatch = "batch"
	// 主页面扩展，增加侧边栏入口及主页面
	ExtentionPage = "page"
)

type Extention interface {
	Clone() Extention
}

type MachineExtention struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Permission string `json:"permission"`
}

type BatchExtention struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Permission string `json:"permission"`
}

type PageExtention struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	IsIndex    bool   `json:"is_index"`
	URL        string `json:"url"`
	Permission string `json:"permission"`
}

func (me *MachineExtention) Clone() Extention {
	result := *me
	return &result
}

func (be *BatchExtention) Clone() Extention {
	result := *be
	return &result
}

func (pe *PageExtention) Clone() Extention {
	result := *pe
	return &result
}

// 解析extentions参数
func ParseParameters(data []map[string]interface{}) []Extention {
	var extentions []Extention
	for _, v := range data {
		switch v["type"] {
		case ExtentionMachine:
			me := &MachineExtention{
				Name:       v["name"].(string),
				URL:        v["url"].(string),
				Permission: v["permission"].(string),
				Type:       v["type"].(string),
			}
			extentions = append(extentions, me)

		case ExtentionBatch:
			be := &BatchExtention{
				Name:       v["name"].(string),
				URL:        v["url"].(string),
				Permission: v["permission"].(string),
				Type:       v["type"].(string),
			}
			extentions = append(extentions, be)

		case ExtentionPage:
			pe := &PageExtention{
				Name:       v["name"].(string),
				URL:        v["url"].(string),
				Permission: v["permission"].(string),
				Type:       v["type"].(string),
				IsIndex:    v["is_index"].(bool),
			}
			extentions = append(extentions, pe)
		}
	}
	return extentions
}
