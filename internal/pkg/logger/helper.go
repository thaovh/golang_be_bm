package logger

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Helper wraps log.Helper with level-aware logging
type Helper struct {
	fileLoggers *FileLoggers
	helper      *log.Helper
}

// NewHelper creates a new logger helper with file logging support
func NewHelper(fileLoggers *FileLoggers) *Helper {
	return &Helper{
		fileLoggers: fileLoggers,
		helper:      log.NewHelper(fileLoggers.StdoutLogger),
	}
}

// WithContext returns a helper with context
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		fileLoggers: h.fileLoggers,
		helper:      log.NewHelper(h.fileLoggers.StdoutLogger).WithContext(ctx),
	}
}

// Log logs at info level
func (h *Helper) Log(level log.Level, keyvals ...interface{}) error {
	switch level {
	case log.LevelInfo:
		return h.fileLoggers.InfoLogger.Log(level, keyvals...)
	case log.LevelWarn:
		return h.fileLoggers.WarningLogger.Log(level, keyvals...)
	case log.LevelError:
		return h.fileLoggers.ErrorLogger.Log(level, keyvals...)
	default:
		return h.fileLoggers.StdoutLogger.Log(level, keyvals...)
	}
}

// Info logs at info level
func (h *Helper) Info(args ...interface{}) {
	h.helper.Info(args...)
	// Also write to info file
	log.NewHelper(h.fileLoggers.InfoLogger).Info(args...)
}

// Infof logs at info level with format
func (h *Helper) Infof(format string, args ...interface{}) {
	h.helper.Infof(format, args...)
	// Also write to info file
	log.NewHelper(h.fileLoggers.InfoLogger).Infof(format, args...)
}

// Warn logs at warning level
func (h *Helper) Warn(args ...interface{}) {
	h.helper.Warn(args...)
	// Also write to warning file
	log.NewHelper(h.fileLoggers.WarningLogger).Warn(args...)
}

// Warnf logs at warning level with format
func (h *Helper) Warnf(format string, args ...interface{}) {
	h.helper.Warnf(format, args...)
	// Also write to warning file
	log.NewHelper(h.fileLoggers.WarningLogger).Warnf(format, args...)
}

// Error logs at error level
func (h *Helper) Error(args ...interface{}) {
	h.helper.Error(args...)
	// Also write to error file
	log.NewHelper(h.fileLoggers.ErrorLogger).Error(args...)
}

// Errorf logs at error level with format
func (h *Helper) Errorf(format string, args ...interface{}) {
	h.helper.Errorf(format, args...)
	// Also write to error file
	log.NewHelper(h.fileLoggers.ErrorLogger).Errorf(format, args...)
}

// Debug logs at debug level
func (h *Helper) Debug(args ...interface{}) {
	h.helper.Debug(args...)
}

// Debugf logs at debug level with format
func (h *Helper) Debugf(format string, args ...interface{}) {
	h.helper.Debugf(format, args...)
}

