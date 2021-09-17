package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ginfra/config"
	"ginfra/datasource"
	"ginfra/log"
	mw "ginfra/middleware"
	"ginfra/models"
	"ginfra/router"

	"gorm.io/gorm/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()

	// init config
	cfg, err := config.Parse("")
	if err != nil {
		panic(err)
	}

	if len(cfg.GetString("db.url")) > 0 {
		// init DB
		fmt.Println(cfg.GetString("db.dialect"), cfg.GetString("db.url"))
		var lv logger.LogLevel = logger.Silent
		if cfg.GetBool("db.logmode") {
			lv = logger.Info
		}
		db, err := datasource.InitDefaultGormDBv2(cfg.GetString("db.url"),
			cfg.GetInt("db.maxopenconns"), cfg.GetInt("db.maxidleconns"), lv)
		if err != nil {
			panic(err)
		}
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&models.Post{},
			&models.Tag{},
			&models.PostTag{},
		)
	}

	// Set gin mode.
	gin.SetMode(cfg.GetString("runmode"))

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// New Zap logger
	logger := log.NewZapLogger("ginfra", cfg.GetString("logfile"), "debug")
	log.ZLog = logger
	defer logger.Sync()

	// Create the Gin engine.
	g := router.New(
		// gin.Context to context
		mw.GinContextToContextMiddleware(),
		// Middlwares. RequestID
		mw.RequestId(),
		// Middlwares. Customize logger, should behind RequestId
		mw.ContextLogger(logger),
		// Middlwares. Request time out
		mw.Timeout(cfg.GetDuration("timeout")),
		// cors
		cors.New(cors.Config{
			AllowOrigins:     cfg.GetStringSlice("cors.origins"),
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return strings.HasSuffix(origin, "qq.com")
			},
			MaxAge: 12 * time.Hour,
		}),
	)

	srv := &http.Server{
		Addr:           cfg.GetString("addr"),
		Handler:        g,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Info("Server Started...")

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of N seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
	}
	logger.Info("Server exiting")
}
