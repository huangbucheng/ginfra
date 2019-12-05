package router

import (
	"ginfra/handler"
	"ginfra/handler/sd"
	mw "ginfra/middleware"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(handlers ...gin.HandlerFunc) *gin.Engine {
	// Create the Gin engine.
	g := gin.New()

	// pprof router
	pprof.Register(g)

	// Middlewares.
	g.Use(handlers...)

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	// g.Use(gin.LoggerWithFormatter(mw.GinCustomLogFormat))

	// The apmgin middleware will recover panics and send them to Elastic APM,
	// so you do not need to install the gin.Recovery middleware.
	// g.Use(apmgin.Middleware(g))
	g.Use(gin.Recovery())

	// metric
	g.Use(mw.Metric())
	g.GET("/metrics", gin.WrapH(promhttp.InstrumentMetricHandler(
		mw.GRegistry, promhttp.HandlerFor(mw.GRegistry, promhttp.HandlerOpts{}),
	)))

	// load routes
	load(g)

	return g
}

// load loads routes.
func load(g *gin.Engine) {
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found.")
	})

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	// User handlers
	g.GET("/ping", handler.Ping)
	g.GET("/timeout", handler.TimedHandler)
	g.GET("/dbtimeout", handler.DBTimedHandler)

	// curl -X POST -H 'content-type: application/x-www-form-urlencoded' -d 'value=abc' "http://127.0.0.1:21000/TagCreate"
	// g.POST("/TagCreate", handler.TagCreate)
	// curl -X POST -H 'content-type: application/x-www-form-urlencoded' -d 'tags=abc&title=First Post&body=balabala...&isPublished=on' "http://127.0.0.1:21000/PostCreate"
	g.POST("/PostCreate", handler.PostCreate)
	// curl "http://127.0.0.1:21000/TagGet?tag=abc"
	// g.GET("/TagGet", handler.TagGet)
	// curl "http://127.0.0.1:21000/UseHttpClient"
	g.GET("/UseHttpClient", handler.UseHttpClient)

	return
}
