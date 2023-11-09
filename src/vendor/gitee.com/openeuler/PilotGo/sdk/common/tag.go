package common

const (
	TypeOk    = "ok"
	TypeWarn  = "warn"
	TypeError = "error"
)

type Tag struct {
	UUID       string `json:"machineuuid"`
	PluginName string `json:"plugin_name"`
	Type       string `json:"type"`
	Data       string `json:"data"`
}
