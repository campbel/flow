package logger

import (
	"io"
	"os"

	"github.com/campbel/flow/meta/config"
	"github.com/charmbracelet/log"
)

var DefaultLogger = NewLogger(os.Stderr, config.DebugLogLevel)

func NewLogger(w io.Writer, debug bool) *log.Logger {
	logger := log.NewWithOptions(w, log.Options{
		TimeFormat:      "15:04:05",
		ReportTimestamp: true,
	})
	if debug {
		logger.SetLevel(log.DebugLevel)
		logger.SetReportCaller(true)
		logger.SetCallerOffset(1)
	}

	return logger
}

func Info(msg any, args ...any) {
	DefaultLogger.Info(msg, args...)
}

func Debug(msg any, args ...any) {
	DefaultLogger.Debug(msg, args...)
}

func Error(msg any, args ...any) {
	DefaultLogger.Error(msg, args...)
}

func Warn(msg any, args ...any) {
	DefaultLogger.Warn(msg, args...)
}
