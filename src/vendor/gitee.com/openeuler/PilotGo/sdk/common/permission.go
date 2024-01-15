package common

type Permission struct {
	Resource string `json:"resource"` //字符串中不允许包含"/"
	Operate  string `json:"operate"`
}
