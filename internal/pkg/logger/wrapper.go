package logger

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// LevelRouterLogger routes logs to different files based on level
type LevelRouterLogger struct {
	fileLoggers *FileLoggers
}

// NewLevelRouterLogger creates a logger that routes logs to appropriate files
func NewLevelRouterLogger(fileLoggers *FileLoggers) log.Logger {
	return &LevelRouterLogger{
		fileLoggers: fileLoggers,
	}
}

// Log implements log.Logger interface
func (l *LevelRouterLogger) Log(level log.Level, keyvals ...interface{}) error {
	switch level {
	case log.LevelInfo:
		return l.fileLoggers.InfoLogger.Log(level, keyvals...)
	case log.LevelWarn:
		return l.fileLoggers.WarningLogger.Log(level, keyvals...)
	case log.LevelError:
		return l.fileLoggers.ErrorLogger.Log(level, keyvals...)
	default:
		return l.fileLoggers.StdoutLogger.Log(level, keyvals...)
	}
}

// WithContext returns a logger with context
func (l *LevelRouterLogger) WithContext(ctx context.Context) log.Logger {
	// Create new loggers with context
	infoLogger := log.NewHelper(l.fileLoggers.InfoLogger).WithContext(ctx)
	warnLogger := log.NewHelper(l.fileLoggers.WarningLogger).WithContext(ctx)
	errorLogger := log.NewHelper(l.fileLoggers.ErrorLogger).WithContext(ctx)
	stdoutLogger := log.NewHelper(l.fileLoggers.StdoutLogger).WithContext(ctx)

	return &LevelRouterLoggerWithContext{
		infoHelper:    infoHelperType{helper: infoLogger},
		warnHelper:    warnHelperType{helper: warnLogger},
		errorHelper:   errorHelperType{helper: errorLogger},
		stdoutHelper:  stdoutHelperType{helper: stdoutLogger},
	}
}

// Helper types for context-aware logging
type infoHelperType struct {
	helper *log.Helper
}

type warnHelperType struct {
	helper *log.Helper
}

type errorHelperType struct {
	helper *log.Helper
}

type stdoutHelperType struct {
	helper *log.Helper
}

// LevelRouterLoggerWithContext is a context-aware logger that routes by level
type LevelRouterLoggerWithContext struct {
	infoHelper   infoHelperType
	warnHelper   warnHelperType
	errorHelper  errorHelperType
	stdoutHelper stdoutHelperType
}

// Log implements log.Logger interface with context
func (l *LevelRouterLoggerWithContext) Log(level log.Level, keyvals ...interface{}) error {
	switch level {
	case log.LevelInfo:
		_ = l.infoHelper.helper.Log(level, keyvals...)
		return nil
	case log.LevelWarn:
		_ = l.warnHelper.helper.Log(level, keyvals...)
		return nil
	case log.LevelError:
		_ = l.errorHelper.helper.Log(level, keyvals...)
		return nil
	default:
		_ = l.stdoutHelper.helper.Log(level, keyvals...)
		return nil
	}
}

// WithContext returns itself (already has context)
func (l *LevelRouterLoggerWithContext) WithContext(ctx context.Context) log.Logger {
	return l
}

