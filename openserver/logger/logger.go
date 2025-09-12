// logger/logger.go
package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func Init(level string, outputPath string) {
	once.Do(func() {
		var lvl zap.AtomicLevel
		switch level {
		case "debug":
			lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "warn":
			lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
		default:
			lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
		}

		// 编码配置
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 友好时间格式
		encoderConfig.StacktraceKey = ""                      // 关闭自动堆栈跟踪（Error 才需要）

		// JSON 编码器
		encoder := zapcore.NewJSONEncoder(encoderConfig)

		// 日志输出目标
		var cores []zapcore.Core

		// 1. 控制台输出
		consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl)
		cores = append(cores, consoleCore)

		// 2. 文件输出（如果指定了路径）
		if outputPath != "" {
			fileCore := zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
				Filename:   outputPath, // 日志文件路径
				MaxSize:    100,        // 每个文件最大 100MB
				MaxBackups: 7,          // 最多保留 7 个备份
				MaxAge:     30,         // 最长保存 30 天
				Compress:   true,       // 启用压缩
			}), lvl)
			cores = append(cores, fileCore)
		}

		// 合并所有 core
		core := zapcore.NewTee(cores...)

		// 构建最终 logger
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

		// 替换全局 logger，以便第三方库也能使用
		zap.ReplaceGlobals(logger)
	})
}

func getLogger() *zap.Logger {
	if logger == nil {
		Init("info", "")
	}
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	getLogger().Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}

func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}
