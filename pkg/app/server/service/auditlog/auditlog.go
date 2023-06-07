package auditlog

import (
	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
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

func New(module, action, msg string, u dao.User) *AuditLog {
	return &AuditLog{
		LogUUID:    uuid.New().String(),
		Module:     module,
		Status:     StatusRunning,
		OperatorID: u.ID,
		Action:     action,
		Message:    msg,
	}
}

func Add(log *dao.AuditLog) error {
	return log.Record()
}

// 修改日志的操作状态
func UpdateStatus(log *dao.AuditLog, status string) error {
	return log.UpdateStatus(status)
}

// 查询所有日志
func Get() ([]dao.AuditLog, error) {
	return dao.GetAuditLog()
}

// 根据父UUid查询日志
func GetByParentId(parentUUId string) (dao.AuditLog, error) {
	return dao.GetAuditLogByParentId(parentUUId)
}

// 查询单条日志
func GetById(logUUId string) (dao.AuditLog, error) {
	return dao.GetAuditLogById(logUUId)
}

func GetByModule(name string) ([]dao.AuditLog, error) {
	return dao.GetAuditLogByModule(name)
}
