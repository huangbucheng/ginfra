package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ginfra/config"
	"ginfra/datasource"
	"ginfra/log"
	mw "ginfra/middleware"
	"ginfra/models"
	"ginfra/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// init config
	cfg, err := config.Parse("")
	if err != nil {
		panic(err)
	}

	// init DB
	fmt.Println(cfg.GetString("db.dialect"), cfg.GetString("db.url"))
	db, err := datasource.InitGormDB(cfg.GetString("db.dialect"), cfg.GetString("db.url"),
		cfg.GetInt("db.maxopenconns"), cfg.GetInt("db.maxidleconns"), cfg.GetBool("db.logmode"))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Post{}, &models.Tag{}, &models.PostTag{})
	db.Model(&models.PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")
	defer db.Close()

	// Set gin mode.
	gin.SetMode(cfg.GetString("runmode"))

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// New Zap logger
	logger := log.NewLogger("ginfra", cfg.GetString("logfile"), "debug")
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
