package model

import (
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

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

func (b *Batch) ReturnBatch(q *PaginationQ) (list *[]Batch, total uint, err error) {
	list = &[]Batch{}
	tx := mysqlmanager.DB.Find(&list)
	total, err = CrudAll(q, tx, list)
	return
}

// type GetBatch struct {
// 	Name        string
// 	Description string
// 	manager     string
// 	Depart      []string
// }
