package datasource

import (
	"context"
	"errors"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormDBv2 *gorm.DB

//Gormv2 获取gorm v2 默认实例
func Gormv2(ctx context.Context) (*gorm.DB, error) {
	if gormDBv2 == nil {
		return nil, errors.New("DB uninitialized")
	}
	return gormDBv2.WithContext(ctx), nil
}

//InitGormDBv2 初始化gorm v2 实例
func InitGormDBv2(dsn string, maxopen, maxidle int, lv logger.LogLevel) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(lv),
		NamingStrategy: schema.NamingStrategy{
			// table name prefix, table for `User` would be `t_users`
			TablePrefix: "t_",
			// use singular table name, table for `User` would be `user` with this option enabled
			SingularTable: true,
			// use name replacer to change struct/field name before convert it to db name
			NameReplacer: strings.NewReplacer("CID", "Cid"),
		},
	})
	if err == nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxOpenConns(maxopen)
			sqlDB.SetMaxIdleConns(maxidle)
		}
		return db, nil
	}
	return nil, err
}

//InitGormDBv2 初始化gorm v2 实例，并设置为默认实例
func InitDefaultGormDBv2(dsn string, maxopen, maxidle int, lv logger.LogLevel) (*gorm.DB, error) {
	db, err := InitGormDBv2(dsn, maxopen, maxidle, lv)
	if err != nil {
		return nil, err
	}
	gormDBv2 = db
	return db, nil
}
