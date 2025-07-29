/*
 * @Author: JimZhang
 * @Date: 2025-07-24 12:17:26
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-25 00:10:22
 * @FilePath: /server/pkg/log/logger.go
 * @Description: 日志
 *
 */
package log

import (
	"os"
	"path/filepath"
	"server/common/config"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var logToFile *logrus.Logger
var loggerFile string

func setLogFile(file string) {
	loggerFile = file
}

// 初始化
func init() {
	// 确保日志目录存在
	_ = os.MkdirAll(config.Config.Log.Path, os.ModePerm)
	setLogFile(filepath.Join(config.Config.Log.Path, config.Config.Log.Name))
}

func Log() *logrus.Logger {
	// 文件输出
	if config.Config.Log.Model == "file" {
		return logFile()
	} else {
		// 控制台输出
		if log == nil {
			log = logrus.New()
			log.Out = os.Stdout
			log.Formatter = &logrus.TextFormatter{TimestampFormat: "2008-01-0115:04:05"}
			log.SetLevel(logrus.DebugLevel)
		}
	}
	return log
}

// 创建日志写入器
func createLogWriter() (*rotatelogs.RotateLogs, error) {
	var logConfig = config.Config.Log
	return rotatelogs.New(
		// 分割后的文件名称
		loggerFile+"_%Y-%m-%d.log",
		// 设置最大保存时间
		rotatelogs.WithMaxAge(time.Duration(logConfig.MaxAge)*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
		// 设置日志文件个数
		// rotatelogs.WithRotationCount(uint(logConfig.MaxBackups)),
		// 设置日志文件大小
		// rotatelogs.WithRotationSize(int64(logConfig.MaxSize)*1024*1024),
	)
}

// 创建 lfshook
func createLfHook(writerMap lfshook.WriterMap) *lfshook.LfsHook {
	return lfshook.NewHook(writerMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
func logFile() *logrus.Logger {
	if logToFile == nil {
		logToFile = logrus.New()
		logToFile.SetLevel(logrus.DebugLevel)

		// 创建文件写入器
		logWriter, err := createLogWriter()
		if err != nil {
			// 如果创建失败，降级到仅控制台输出
			logToFile.Out = os.Stdout
			logToFile.Formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"}
			return logToFile
		}

		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}

		// 新增 Hook
		logToFile.AddHook(createLfHook(writeMap))
	}
	return logToFile
}
