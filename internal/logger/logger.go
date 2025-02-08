package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var (
	// Logger is the global logger instance
	Logger *log.Logger
)

// InitLogger initializes the global logger
func InitLogger() {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.DebugLevel,
		Prefix:          "ffmpeg-api",
	})
}

// Debug logs a debug message
func Debug(msg interface{}, keyvals ...interface{}) {
	Logger.Debug(msg, keyvals...)
}

// Info logs an info message
func Info(msg interface{}, keyvals ...interface{}) {
	Logger.Info(msg, keyvals...)
}

// Warn logs a warning message
func Warn(msg interface{}, keyvals ...interface{}) {
	Logger.Warn(msg, keyvals...)
}

// Error logs an error message
func Error(msg interface{}, keyvals ...interface{}) {
	Logger.Error(msg, keyvals...)
}

// Fatal logs a fatal message and exits
func Fatal(msg interface{}, keyvals ...interface{}) {
	Logger.Fatal(msg, keyvals...)
}

// With returns a new logger with the given key-value pairs
func With(keyvals ...interface{}) *log.Logger {
	return Logger.With(keyvals...)
}
