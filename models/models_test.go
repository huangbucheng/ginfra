package models

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"
	"testing"

	"ginfra/datasource"
	
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

var db *gorm.DB
var mock sqlmock.Sqlmock

func init() {
	var (
		d *sql.DB
	)

	d, mock, _ = sqlmock.New()
	db, _ = gorm.Open("mysql", d)
	// db.LogMode(false)
	datasource.SetGormDB(db)
}

func Test_GetPostById(t *testing.T) {
	var (
		id    = 1
		title = "post title"
		body  = "blabla..."
		view  = 10
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL AND ((id = ?)) ORDER BY `posts`.`id` ASC LIMIT 1")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"title", "body", "view"}).
			AddRow(title, body, view))

	res, err := GetPostById(context.Background(), strconv.Itoa(id))
	convey.Convey("models.GetPostById", t, func() {
		convey.So(err, convey.ShouldEqual, nil)
	})
	convey.Convey("models.GetPostById", t, func() {
		convey.So(res, convey.ShouldResemble, &Post{Title: title, Body: body, View: view})
	})
}
