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
	Logger *LoggerWrap
)

type LoggerWrap struct {
	*zap.Logger
}

func NewLogger(loggername, filename, level string) *LoggerWrap {
	Logger = &LoggerWrap{
		Logger: NewZapLogger(loggername, filename, level),
	}
	return Logger
}

// NewLogger returns a new zap logger
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
	filed := zap.Fields(zap.String("Logger", loggername))
	// 构造日志
	logger := zap.New(core, caller, development, filed)

	return logger
}

// Log returns the Request ID scoped logger from the request Context
// and panics if it cannot be found. This function is only ever used
// by your controllers if your app uses the RequestID middlewares,
// otherwise you should use the controller's receiver logger directly.
func (log *LoggerWrap) WithGinContext(c *gin.Context) *LoggerWrap {
	if c == nil {
		return log
	}

	// with custom log fields
	return log.WithContext(c.Request.Context())
}

func (log *LoggerWrap) WithContext(c context.Context) *LoggerWrap {
	var l *zap.Logger
	// with custom log fields
	if ctxLoggerFields, ok := c.Value(CtxLoggerFields).(map[string]string); ok {
		for loggerField := range ctxLoggerFields {
			l = l.With(zap.String(loggerField, ctxLoggerFields[loggerField]))
		}
	}

	return &LoggerWrap{
		Logger: l,
	}
}

func SetFieldsByGin(c *gin.Context, fields map[string]string) {
	if c == nil {
		return
	}

	// set custom log fields
	SetFields(c.Request.Context(), fields)
}

func SetFields(c context.Context, fields map[string]string) {
	// set custom log fields
	if ctxLoggerFields, ok := c.Value(CtxLoggerFields).(map[string]string); ok {
		for field, value := range fields {
			ctxLoggerFields[field] = value
		}
	}
}
