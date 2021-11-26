package model

/**
 * @Author: zhang han
 * @Date: 2021/11/18 14:00
 * @Description:防火墙
 */

type ZonePort struct {
	Zone string `gorm:"type:varchar(25);not null" json:"zone,omitempty" form:"zone"`
	Port int    `gorm:"type:varchar(25);not null" json:"port,omitempty" form:"port"`
}
