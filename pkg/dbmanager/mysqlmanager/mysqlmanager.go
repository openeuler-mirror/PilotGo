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
 * LastEditTime: 2022-04-25 16:00:43
 * Description: mysql初始化
 ******************************************************************************/
package mysqlmanager

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB  *gorm.DB
	Url string
)

func MysqlInit(ip, username, password, dbname string, port int) (*MysqlManager, error) {
	m := &MysqlManager{
		ip:       ip,
		port:     port,
		userName: username,
		passWord: password,
		dbName:   dbname,
	}
	Url = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		m.userName,
		m.passWord,
		m.ip,
		m.port,
		m.dbName)
	var err error
	m.db, err = gorm.Open("mysql", Url)
	if err != nil {
		return nil, err
	}
	DB = m.db
	m.db.DB().SetMaxIdleConns(10)
	m.db.DB().SetMaxOpenConns(100)
	//禁用复数
	m.db.SingularTable(true)
	return m, nil
}
