package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ginfra/datasource"
	"ginfra/log"
	mw "ginfra/middleware"

	. "github.com/agiledragon/gomonkey"
	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

var g *gin.Engine

func init() {
	// Set gin mode.
	gin.SetMode(gin.ReleaseMode)

	// Create the Gin engine.
	g = gin.New()

	// init DB
	_, err := datasource.InitGormDB("sqlite3", "/tmp/gorm.db", 1, 1, false)
	if err != nil {
		panic(err)
	}

	// New Zap logger
	logger := log.NewLogger("ginfra", "/tmp/gin.log", "info")

	// Routes and Middlewares.
	g.Use(mw.ContextLogger(logger))
	g.Use(mw.Timeout(time.Second * 2))
	g.Use(mw.RequestId())

	g.GET("/ping", Ping)
	g.POST("/PostCreate", PostCreate)
	g.GET("/UseHttpClient", UseHttpClient)
}

// Ping - httpexpect
func Test_Ping(t *testing.T) {

	// create a server for testing
	server := httptest.NewServer(g)
	defer server.Close()

	// create a test engine from server
	e := httpexpect.New(t, server.URL)

	// expect get / status is 200
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		Text().Equal("pong")
}

// Http Client
func Test_UseHttpClient(t *testing.T) {
	patches := ApplyFunc(req.Get, func(_ string, _ ...interface{}) (*req.Resp, error) {
		return nil, errors.New("failed")
	})
	defer patches.Reset()

	// create a server for testing
	server := httptest.NewServer(g)
	defer server.Close()

	// create a test engine from server
	e := httpexpect.New(t, server.URL)

	// expect get / status is 200
	e.GET("/UseHttpClient").
		Expect().
		Status(http.StatusInternalServerError)
}
