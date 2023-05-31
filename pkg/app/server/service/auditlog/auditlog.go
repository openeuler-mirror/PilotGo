package auditlog

import (
	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
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
	LogTypeUser       = "用户"
	LogTypePermission = "权限"
	LogTypePlugin     = "插件"
	LogTypeBatch      = "批次"
	LogTypeOrganize   = "组织"
	LogTypeMachine    = "机器"
)

type AuditLog = dao.AuditLog

func NewAuditLog(module, action, msg string, u service.User) *AuditLog {
	return &AuditLog{
		LogUUID:    uuid.New().String(),
		Module:     module,
		Status:     StatusRunning,
		OperatorID: u.ID,
		Action:     action,
		Message:    msg,
	}
}

func AddAuditLog(log *dao.AuditLog) error {
	return log.Record()
}

// 修改日志的操作状态
func UpdateStatus(log *dao.AuditLog, status string) error {
	return log.UpdateStatus(status)
}

// 查询所有日志
func GetAuditLog() ([]dao.AuditLog, error) {
	return dao.GetAuditLog()
}

// 根据父UUid查询日志
func GetAuditLogByParentId(parentUUId string) (dao.AuditLog, error) {
	return dao.GetAuditLogByParentId(parentUUId)
}

// 查询单条日志
func GetAuditLogById(logUUId string) (dao.AuditLog, error) {
	return dao.GetAuditLogById(logUUId)
}

func GetAuditLogByModule(name string) ([]dao.AuditLog, error) {
	return dao.GetAuditLogByModule(name)
}
