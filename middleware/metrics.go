package middleware

import (
	"strconv"
	"time"

	"ginfra/config"
	"ginfra/log"
	"ginfra/plugin/atta"
	"ginfra/protocol"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var httpRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ginfra_http_request_count",
		Help: "http request count",
	},
	[]string{"method", "path", "status"},
)

var httpRequestDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "ginfra_http_request_duration",
		Help: "http request duration",
	},
	[]string{"method", "path"},
)

var (
	GRegistry *prometheus.Registry

	enableATTA bool
	attaID     string
	attaToken  string
)

func init() {
	GRegistry = prometheus.NewRegistry()
	GRegistry.Register(httpRequestCount)
	GRegistry.Register(httpRequestDuration)
	// GRegistry.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	// GRegistry.Register(prometheus.NewGoCollector())

	cfg, err := config.Parse("")
	if err != nil {
		panic(err)
	}

	// atta
	enableATTA = cfg.GetBool("atta.enable")
	attaID = cfg.GetString("atta.attaid")
	attaToken = cfg.GetString("atta.token")
}

//DumpPromMetrics -
func DumpPromMetrics(filename string) {
	// Write metrics to files
	prometheus.WriteToTextfile(filename, GRegistry)

	// APIs for reference

	// 1. Parse Text Metrics
	// var parser expfmt.TextParser
	// metricFamilies, err := parser.TextToMetricFamilies(in)

	// 2. Gather metircs
	// metrics, err := GRegistry.Gather()
}

//Metric metric middleware
func Metric() gin.HandlerFunc {
	return func(c *gin.Context) {
		tBegin := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		if path == "/metrics" || path == "/" {
			return
		}

		duration := float64(time.Since(tBegin)) / float64(time.Second)
		tEnd := time.Now()
		latency := tEnd.Sub(tBegin)

		// 请求数加1
		httpRequestCount.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Inc()

		// 上报ATTA
		if enableATTA {
			atta.ReportBackendRequestStatus(attaID, attaToken, protocol.GetUserId(c),
				path, protocol.GetResponseCode(c), c.Writer.Status(), int(latency/time.Millisecond))
		}

		//  记录本次请求处理时间
		httpRequestDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
		}).Observe(duration)

		fields := []zap.Field{zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("userid", protocol.GetUserId(c)),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("etime", tEnd.Format(time.RFC3339)),
			zap.Duration("latency", latency),
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				log.WithGinContext(c).Error(e, fields...)
			}
		} else {
			log.WithGinContext(c).Info("summary", fields...)
		}
	}
}
