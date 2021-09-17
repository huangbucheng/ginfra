package models

import (
	"strconv"

	"gorm.io/gorm"
	//_ "github.com/go-sql-driver/mysql"
)

//OperationLog 操作日志的表
type OperationLog struct {
	gorm.Model
	Event   string `gorm:"index;size:64"`
	ErrCode string `gorm:"size:64"`

	// private field, ignored from gorm
	TableID uint `gorm:"-"`
}

//OperationLogTable 自定义Table Name
func OperationLogTable(s *OperationLog) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("t_operation_log_" + strconv.Itoa(int(s.TableID)))
	}
}

//AutoMigrate 初始化表
func (s *OperationLog) AutoMigrate(db *gorm.DB) error {
	return db.Scopes(OperationLogTable(s)).AutoMigrate(&OperationLog{})
}

//Insert 插入操作日志
func (s *OperationLog) Insert(db *gorm.DB) error {
	return db.Scopes(OperationLogTable(s)).Create(s).Error
	//return db.Create(s).Error
}

//Delete 軟刪除操作日志
func (s *OperationLog) Delete(db *gorm.DB) error {
	return db.Scopes(OperationLogTable(s)).Delete(s).Error
}
