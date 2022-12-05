package model

import "time"

type ConfigFile struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MachineUUID string `json:"uuid"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	UpdatedAt   time.Time
}
