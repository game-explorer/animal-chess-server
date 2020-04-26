package log

import (
	"github.com/op/go-logging"
	"os"
)

var logger *logging.Logger

func Debug(args ...interface{})   { logger.Debug(args...) }
func Info(args ...interface{})    { logger.Info(args...) }
func Warning(args ...interface{}) { logger.Warning(args...) }
func Error(args ...interface{})   { logger.Error(args...) }
func Fatal(args ...interface{})   { logger.Fatal(args...) }
func Panic(args ...interface{})   { logger.Panic(args...) }

func Debugf(format string, args ...interface{})   { logger.Debugf(format, args...) }
func Infof(format string, args ...interface{})    { logger.Infof(format, args...) }
func Warningf(format string, args ...interface{}) { logger.Warningf(format, args...) }
func Errorf(format string, args ...interface{})   { logger.Errorf(format, args...) }
func Fatalf(format string, args ...interface{})   { logger.Fatalf(format, args...) }
func Panicf(format string, args ...interface{})   { logger.Panicf(format, args...) }

func Native() *logging.Logger {
	return logger
}

func init() {
	SetDebug(true)
}

// debug会影响颜色 和 最低等级
func SetDebug(logDebug bool) {
	logger = New(logDebug)
}

func New(logDebug bool) *logging.Logger {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := "%{time:2006-01-02 15:04:05.999} %{longfile} %{shortfunc} >> [%{level:.4s}] %{message}"
	// debug模式会有颜色, 但在主机上, 颜色代码会乱码, 所以生产环境不应该启用
	if logDebug {
		format = "%{color}%{time:2006-01-02 15:04:05.999} %{longfile} %{shortfunc} >> [%{level:.4s}]%{color:reset} %{message}"
	}
	f := logging.MustStringFormatter(format)
	backendFormatter := logging.NewBackendFormatter(backend, f)
	b := logging.MultiLogger(backendFormatter)
	if !logDebug {
		b.SetLevel(logging.INFO, "")
	}

	logger := logging.MustGetLogger("")
	logger.ExtraCalldepth = 1
	logger.SetBackend(b)
	return logger
}
