package common

type AsyncCmdResult struct {
	TaskID string       `json:"task_id"`
	Result []*CmdResult `json:"result"`
}

type CmdResult struct {
	MachineUUID string `json:"machine_uuid"`
	MachineIP   string `json:"machine_ip"`
	RetCode     int    `json:"retcode"`
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
}

type CmdStruct struct {
	Batch   *Batch `json:"batch"`
	Command string `json:"command"`
}
