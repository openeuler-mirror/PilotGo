package auditlog

import (
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/app/server/dao"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 日志执行操作状态
const (
	StatusSuccess     = "成功"
	StatusPartSuccess = "部分成功"
	StatusRunning     = "运行中"
	StatusFail        = "失败"
)

// 日志记录归属模块
const (
	LogTypeUser       = "用户" // 登录 注销(父日志没有创建者和部门信息) 添加 删除 修改密码 重置密码 修改用户信息
	LogTypePermission = "权限" // 角色权限 编辑角色 删除角色 添加角色
	LogTypePlugin     = "插件" // null
	LogTypeBatch      = "批次" // 添加批次 删除批次 编辑批次
	LogTypeOrganize   = "组织"
	LogTypeMachine    = "机器"       // null
	LogTypeRPM        = "软件包安装/卸载" // rpm安装 rpm卸载
	LogTypeService    = "运行服务"     // null
	LogTypeSysctl     = "配置内核参数"   // null
	LogTypeBroadcast  = "配置文件下发"   // 配置文件下发
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

// deprecated
func New(module, action string, userid uint) *dao.AuditLog {
	return &dao.AuditLog{
		LogUUID: uuid.New().String(),
		Module:  module,
		Status:  "",
		UserID:  userid,
		Action:  action,
	}
}

func New_sub(module, action, agentUUID, message, status string, userid uint) *dao.AuditLog {
	return &dao.AuditLog{
		LogUUID:   uuid.New().String(),
		AgentUUID: agentUUID,
		Module:    module,
		Status:    status,
		UserID:    userid,
		Action:    action,
		Message:   message,
	}
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

// 查询单条日志
func GetById(logUUId string) (dao.AuditLog, error) {
	return dao.GetAuditLogById(logUUId)
}

func GetByModule(name string) ([]dao.AuditLog, error) {
	return dao.GetAuditLogByModule(name)
}