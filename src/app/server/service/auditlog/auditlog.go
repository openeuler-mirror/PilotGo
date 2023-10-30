package auditlog

import (
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gorm.io/gorm"
)

// 日志执行操作状态
const (
	StatusOK     = "OK"
	StatusFailed = "failed"
)

// 日志记录归属模块
const (
	ModuleUser    = "user"    // 登录 注销(父日志没有创建者和部门信息) 添加 删除 修改密码 重置密码 修改用户信息
	ModuleRole    = "role"    // 角色权限 编辑角色 删除角色 添加角色
	ModulePlugin  = "plugin"  // null
	ModuleBatch   = "batch"   // 添加批次 删除批次 编辑批次
	ModuleMachine = "machine" // null
	ModuleDepart  = "depart"
	//LogTypeRPM       = "软件包安装/卸载" // rpm安装 rpm卸载
	//LogTypeService   = "运行服务"     // null
	//LogTypeSysctl    = "配置内核参数"   // null
	//LogTypeBroadcast = "配置文件下发"   // 配置文件下发
)

type AuditLog = dao.AuditLog

// 单机操作成功状态:是否成功，机器数量，比率
const (
	ActionOK    = "1,1,1.00"
	ActionFalse = "0,1,0.00"
)

// 计算批量机器操作的状态：成功数，总数目，比率
func BatchActionStatus(StatusCodes []string) (status string) {
	var StatusOKCounts int
	for _, success := range StatusCodes {
		if success == strconv.Itoa(http.StatusOK) {
			StatusOKCounts++
		}
	}
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(StatusOKCounts)/float64(len(StatusCodes))), 64)
	rate := strconv.FormatFloat(num, 'f', 2, 64)
	status = strconv.Itoa(StatusOKCounts) + "," + strconv.Itoa(len(StatusCodes)) + "," + rate
	return
}

// 计算json返回状态
func ActionStatus(StatusCodes []string) (ok bool) {
	for _, code := range StatusCodes {
		if code == strconv.Itoa(http.StatusBadRequest) {
			return false
		} else {
			continue
		}
	}
	return true
}

func Add(log *dao.AuditLog) error {
	return log.Record()
}

// 修改日志的操作状态
func UpdateStatus(log *dao.AuditLog, status string) error {
	return log.UpdateStatus(status)
}

// 添加message信息
func UpdateMessage(log *dao.AuditLog, message string) error {
	return log.UpdateMessage(message)
}

// 查询所有日志
func Get() (*[]dao.AuditLog, *gorm.DB, error) {
	return dao.GetAuditLog()
}

// 查询子日志
func GetAuditLogById(logUUId string) (*[]dao.AuditLog, *gorm.DB, error) {
	return dao.GetAuditLogById(logUUId)
}

// 查询父日志为空的记录
func GetParentLog() (*[]AuditLog, *gorm.DB, error) {
	return dao.GetParentLog()
}

func GetByModule(name string) ([]dao.AuditLog, error) {
	return dao.GetAuditLogByModule(name)
}
