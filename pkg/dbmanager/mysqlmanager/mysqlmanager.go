/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-10-26 09:05:39
 * LastEditTime: 2023-06-28 15:57:07
 * Description: mysql初始化
 ******************************************************************************/
package mysqlmanager

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var global_db *gorm.DB

func ensureDatabase(m *MysqlManager) error {
	Url := fmt.Sprintf("%s:%s@(%s:%d)/?charset=utf8mb4&parseTime=true",
		m.userName,
		m.passWord,
		m.ip,
		m.port)
	db, err := gorm.Open(mysql.Open(Url))
	if err != nil {
		return err
	}

	creatDataBase := "CREATE DATABASE IF NOT EXISTS " + m.dbName + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci"
	db.Exec(creatDataBase)

	d, err := db.DB()
	if err != nil {
		return err
	}
	if err = d.Close(); err != nil {
		return err
	}
	
	return nil
}

func MysqlInit(ip, username, password, dbname string, port int) (*MysqlManager, error) {
	m := &MysqlManager{
		ip:       ip,
		port:     port,
		userName: username,
		passWord: password,
		dbName:   dbname,
	}
	var err error
	err = ensureDatabase(m)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		m.userName,
		m.passWord,
		m.ip,
		m.port,
		m.dbName)
	// gorm db does not need to close
	m.db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	global_db = m.db

	var db *sql.DB
	if db, err = m.db.DB(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return m, nil
}

type MysqlManager struct {
	ip       string
	port     int
	userName string
	passWord string
	dbName   string
	db       *gorm.DB
}

func MySQL() *gorm.DB {
	return global_db
}
