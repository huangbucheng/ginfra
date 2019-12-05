package datasource

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GormWithContext(ctx context.Context) (*gorm.DB, error) {
	db, err := Gorm()
	if db == nil {
		return nil, err
	}

	ctxdb, err := gorm.Open("mysql", &sqlCtxDB{
		underlying: db.DB(),
		ctx:        ctx,
	})
	if err == nil {
		ctxdb.SingularTable(true)
	}
	return ctxdb, err
}

func Gorm() (*gorm.DB, error) {
	if gormDB == nil {
		return nil, errors.New("DB uninitialized")
	}
	return gormDB, nil
}

var gormDB *gorm.DB

func InitGormDB(dialect, source string, maxopen, maxidle int, logmode bool) (*gorm.DB, error) {

	db, err := gorm.Open(dialect, source)
	if err == nil {
		gormDB = db
		db.DB().SetMaxOpenConns(maxopen)
		db.DB().SetMaxIdleConns(maxidle)
		db.SingularTable(true)
		db.LogMode(logmode)
		return db, nil
	}
	return nil, err
}

func SetGormDB(db *gorm.DB) {
	gormDB = db
}

type sqlCtxDB struct {
	underlying *sql.DB
	ctx        context.Context
}

func (db *sqlCtxDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.underlying.ExecContext(db.ctx, query, args...)
}

func (db *sqlCtxDB) Prepare(query string) (*sql.Stmt, error) {
	return db.underlying.PrepareContext(db.ctx, query)
}

func (db *sqlCtxDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.underlying.QueryContext(db.ctx, query, args...)
}

func (db *sqlCtxDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.underlying.QueryRowContext(db.ctx, query, args...)
}

func (db *sqlCtxDB) Begin() (*sql.Tx, error) {
	return db.underlying.BeginTx(db.ctx, nil)
}

func (db *sqlCtxDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.underlying.BeginTx(ctx, opts)
}
