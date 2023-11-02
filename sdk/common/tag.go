package common

const (
	TypeOk    = "ok"
	TypeWarn  = "warn"
	TypeError = "error"
)

type TageMessage struct {
	UUID string `json:"machineuuid"`
	Type string `json:"type"`
	Data string `json:"data"`
}
