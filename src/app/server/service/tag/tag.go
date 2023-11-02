package tag

type TageMessage struct {
	UUID string `json:"machineuuid"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// 向所有插件发送uuidlist
func RequestTag(UUIDList []string) {
	//获取插件列表
	//向url发送请求
}
