package model

import "github.com/jinzhu/gorm"

/**
 * @Author: wang hao
 * @Date: 2021/12/23 17:00
 * @Description:批次属性列表
 */

type Batch struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(100)" json:"description"`
	Manager     string `gorm:"type:varchar(100)" json:"manager"`
	Machinelist string `gorm:"type:varchar(100)" json:"machinelist"`
	Depart      string `gorm:"type:varchar(100)"`
}

// type GetBatch struct {
// 	Name        string
// 	Description string
// 	manager     string
// 	Depart      []string
// }
