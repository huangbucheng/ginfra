package log

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ZLog *zap.Logger
)

//ContextLogger 日志封装
type ContextLogger struct {
	*zap.Logger

	// context info associated with the logger
	LogContext map[string]string
}

// NewZapLogger returns a new zap logger
func NewZapLogger(loggername, filename, level string) *zap.Logger {
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
	filed := zap.Fields(zap.String("logger", loggername))
	// 构造日志
	logger := zap.New(core, caller, development, filed)

	return logger
}

//NewContextLogger new 日志实例
func NewContextLogger(l *zap.Logger) *ContextLogger {
	if l == nil {
		panic("zap logger is nil")
	}

	return &ContextLogger{
		Logger:     l,
		LogContext: make(map[string]string),
	}
}

//WithGinContext 从gin context中取出日志实例
func WithGinContext(c *gin.Context) *zap.Logger {
	if c == nil {
		return ZLog
	}

	return WithContext(c.Request.Context())
}

//WithGinContext 从context中取出日志实例
func WithContext(c context.Context) *zap.Logger {
	var logger *zap.Logger
	l := Logger2(c)
	logger = l.Logger

	// with context fields
	for key := range l.LogContext {
		logger = logger.With(zap.String(key, l.LogContext[key]))
	}

	return logger
}

//WithGinContext 从gin context中取出日志实例
func Logger(c *gin.Context) *ContextLogger {
	if c == nil {
		return NewContextLogger(ZLog)
	}

	return Logger2(c.Request.Context())
}

//WithGinContext 从context中取出日志实例
func Logger2(c context.Context) *ContextLogger {
	if l, ok := c.Value(CtxLoggerKey).(*ContextLogger); ok {
		return l
	}

	return NewContextLogger(ZLog)
}

//Set 往日志实例中设置关键字段
func (l *ContextLogger) Set(key, value string) *ContextLogger {
	// set context field
	l.LogContext[key] = value

	return l
}

//UnSet 取消日志实例中关键字段设置
func (l *ContextLogger) UnSet(key string) *ContextLogger {
	// unset context field
	delete(l.LogContext, key)

	return l
}

//With 返回携带关注字段的日志实例
func (l *ContextLogger) With() *zap.Logger {
	var logger *zap.Logger
	logger = l.Logger

	// with context fields
	for key := range l.LogContext {
		logger = logger.With(zap.String(key, l.LogContext[key]))
	}

	return logger
}
