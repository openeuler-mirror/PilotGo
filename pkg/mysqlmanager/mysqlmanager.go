package mysqlmanager

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func Init(ip, username, password, dbname string, port int) (*MysqlManager, error) {
	m := &MysqlManager{
		ip:       ip,
		port:     port,
		userName: username,
		passWord: password,
		dbName:   dbname,
	}
	url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		m.userName,
		m.passWord,
		m.ip,
		m.port,
		m.dbName)
	var err error
	fmt.Println(url)
	m.db, err = gorm.Open("mysql", url)
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

func (m *MysqlManager) Insert(value interface{}, count int64) error {
	db := m.db.Create(value)
	if db == nil {
		return errors.New("insert db == nil")
	}

	if db.Error != nil {
		db.Rollback()
		return db.Error
	}
	if db.RowsAffected != count {
		return errors.New("insert failed:not enough row")
	}
	return db.Error
}

func (m *MysqlManager) Delete(key_name string, keys []string, sqlType interface{}) error {
	str := key_name + " in ("
	for i := 0; i < len(keys)-1; i++ {
		str += "'" + keys[i] + "'" + ","
	}
	str += "'" + keys[len(keys)-1] + "'" + ")"
	db1 := m.db.Where(str).Delete(sqlType)
	if db1 == nil {
		return errors.New("delete db == nil")
	}

	if db1.Error != nil {
		db1.Rollback()
		return db1.Error
	}
	return nil
}

func (m *MysqlManager) Update(value interface{}) error {
	db1 := m.db.Model(value).Update(value)
	if db1 == nil {
		return errors.New("yodate db == nil")
	}

	if db1.Error != nil {
		db1.Rollback()
		return db1.Error
	}
	return nil
}

func GetPluginInfo(m *MysqlManager) (values []PluginInfo, e error) {
	m.db.Find(&values)
	return values, nil
}

func GetMachInfo(m *MysqlManager) (mi []MachInfo, e error) {
	m.db.Find(&mi)
	return mi, nil
}
