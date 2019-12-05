package middleware

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggers         map[string]*zap.Logger
	loggerMu        sync.Mutex
	firstloggername string // first logger used as default logger
)

func init() {
	loggers = make(map[string]*zap.Logger)
}

// NewLogger returns a new zap logger
func NewLogger(loggername, filename, level string) *zap.Logger {
	if _, ok := loggers[loggername]; ok {
		return loggers[loggername]
	}

	hook := lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		// EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.UnmarshalText([]byte(level))

	multiSyncer := []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	if atomicLevel.Enabled(zapcore.DebugLevel) {
		multiSyncer = append(multiSyncer, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),       // 编码器配置
		zapcore.NewMultiWriteSyncer(multiSyncer...), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("LoggerName", loggername))
	// 构造日志
	logger := zap.New(core, caller, development, filed)

	loggerMu.Lock()
	defer loggerMu.Unlock()
	loggers[loggername] = logger
	if firstloggername == "" {
		firstloggername = loggername
	}

	return logger
}

// GinCustomLogFormat defines custom log format
func GinCustomLogFormat(param gin.LogFormatterParams) string {

	// your custom format
	requestId := ""
	if v, found := param.Keys["X-Request-Id"]; found {
		requestId = v.(string)
	}
	return fmt.Sprintf("%s | %d | %s | %s | %s | %s | %s | %s | %s\n",
		param.TimeStamp.Format(time.RFC3339),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		requestId,
		param.Request.Method,
		param.Path,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// ContextLogger middleware creates a derived logger to include logging of the
// Request ID, and inserts it into the context object
func ContextLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var derivedLogger = logger.With(zap.String("client_ip", c.ClientIP()))
		if ctxReqId, ok := c.Value(CtxRequestID).(string); ok {
			derivedLogger = derivedLogger.With(zap.String("request_id", ctxReqId))
		}
		// c.Set(CtxLoggerFields, make(map[string]string))

		ctx := context.WithValue(c.Request.Context(), CtxLoggerKey, derivedLogger)
		ctx = context.WithValue(ctx, CtxLoggerFields, make(map[string]string))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func Ginzap() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				Logger(c).Error(e)
			}
		} else {
			Logger(c).Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("etime", end.Format(time.RFC3339)),
				zap.Duration("latency", latency),
			)
		}
	}
}

// Log returns the Request ID scoped logger from the request Context
// and panics if it cannot be found. This function is only ever used
// by your controllers if your app uses the RequestID middlewares,
// otherwise you should use the controller's receiver logger directly.
func Logger(c *gin.Context) *zap.Logger {
	if c == nil {
		return GetDefaultLogger()
	}

	v := c.Request.Context().Value(CtxLoggerKey)
	log, ok := v.(*zap.Logger)
	if !ok {
		// panic("cannot get derived request id logger from context object")
		return GetDefaultLogger()
	}

	// with custom log fields
	if ctxLoggerFields, ok := c.Request.Context().Value(CtxLoggerFields).(map[string]string); ok {
		for loggerField := range ctxLoggerFields {
			log = log.With(zap.String(loggerField, ctxLoggerFields[loggerField]))
		}
	}

	return log
}

func Logger2(c context.Context) *zap.Logger {
	v := c.Value(CtxLoggerKey)
	log, ok := v.(*zap.Logger)
	if !ok {
		// panic("cannot get derived request id logger from context object")
		return GetDefaultLogger()
	}

	// with custom log fields
	if ctxLoggerFields, ok := c.Value(CtxLoggerFields).(map[string]string); ok {
		for loggerField := range ctxLoggerFields {
			log = log.With(zap.String(loggerField, ctxLoggerFields[loggerField]))
		}
	}

	return log
}

func SetLoggerField(c *gin.Context, fields map[string]string) {
	// set custom log fields
	if ctxLoggerFields, ok := c.Request.Context().Value(CtxLoggerFields).(map[string]string); ok {
		for field, value := range fields {
			ctxLoggerFields[field] = value
		}
	}
}

// GetDefaultLogger return the first logger
func GetDefaultLogger() *zap.Logger {
	if firstloggername == "" {
		return NewLogger("default", "/tmp/ginfra.log", "debug")
	}
	return GetLogger(firstloggername)
}

// GetDefaultLogger return the logger with the specified logger name
func GetLogger(loggername string) *zap.Logger {
	if logger, ok := loggers[loggername]; ok {
		return logger
	}

	return nil
}
