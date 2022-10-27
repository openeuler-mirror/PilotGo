package model

import (
	"time"
)

type Script struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description"`
	UpdatedAt   time.Time
	Version     string `gorm:"unique" json:"version"`
	Deleted     int    `json:"deleted"` //deleted为1的时候表示删除，一般表示为0
}
