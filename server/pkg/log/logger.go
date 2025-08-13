/*
 * @Author: JimZhang
 * @Date: 2025-07-24 12:17:26
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-08-11 00:00:00
 * @FilePath: /server/pkg/log/logger.go
 * @Description: 使用 Zap 记录日志
 */
package log

import (
	"os"
	"path/filepath"
	"server/common/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

func ensureDir(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}

func newZapLogger() *zap.SugaredLogger {
	logCfg := config.Config.Log

	// 编码器配置
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var ws zapcore.WriteSyncer
	if logCfg.Model == "file" {
		ensureDir(logCfg.Path)
		filePath := filepath.Join(logCfg.Path, logCfg.Name)
		ws = zapcore.AddSync(&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    logCfg.MaxSize, // MB
			MaxBackups: logCfg.MaxBackups,
			MaxAge:     logCfg.MaxAge, // days
			Compress:   true,
		})
	} else {
		ws = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), ws, zapcore.InfoLevel)
	// 开启 caller 和 stacktrace
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return zapLogger.Sugar()
}

// Log 返回全局日志实例（保持函数名不变，便于兼容）
func Log() *zap.SugaredLogger {
	if logger == nil {
		logger = newZapLogger()
	}
	return logger
}
