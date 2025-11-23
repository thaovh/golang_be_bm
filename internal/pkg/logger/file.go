package logger

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// LogDir is the directory where log files are stored
	LogDir = "logs"
	
	// MaxSize is the maximum size of a log file in megabytes before rotation
	MaxSize = 10 // 10MB
	
	// MaxBackups is the maximum number of old log files to retain
	MaxBackups = 10
	
	// MaxAge is the maximum number of days to retain old log files
	MaxAge = 30 // 30 days
	
	// Compress determines if rotated log files should be compressed
	Compress = true
)

// FileLoggerConfig configures file logging
type FileLoggerConfig struct {
	LogDir    string
	MaxSize   int // MB
	MaxBackups int
	MaxAge    int // days
	Compress  bool
}

// DefaultFileLoggerConfig returns default configuration
func DefaultFileLoggerConfig() *FileLoggerConfig {
	return &FileLoggerConfig{
		LogDir:     LogDir,
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   Compress,
	}
}

// FileLoggers contains separate loggers for different log levels
type FileLoggers struct {
	InfoLogger    log.Logger
	WarningLogger log.Logger
	ErrorLogger   log.Logger
	StdoutLogger  log.Logger // For console output
}

// NewFileLoggers creates file loggers with rotation support
func NewFileLoggers(config *FileLoggerConfig, serviceID, serviceName, serviceVersion string) (*FileLoggers, error) {
	if config == nil {
		config = DefaultFileLoggerConfig()
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, err
	}

	// Create lumberjack writers for each log level
	infoWriter := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, "info.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	warningWriter := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, "warning.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	errorWriter := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, "error.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// Create loggers with structured fields
	// Info logger: write to file only (stdout will be handled by helper)
	infoLogger := log.With(
		log.NewStdLogger(infoWriter),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", serviceID,
		"service.name", serviceName,
		"service.version", serviceVersion,
		"level", "info",
	)

	// Warning logger: write to file only
	warningLogger := log.With(
		log.NewStdLogger(warningWriter),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", serviceID,
		"service.name", serviceName,
		"service.version", serviceVersion,
		"level", "warning",
	)

	// Error logger: write to file only
	errorLogger := log.With(
		log.NewStdLogger(errorWriter),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", serviceID,
		"service.name", serviceName,
		"service.version", serviceVersion,
		"level", "error",
	)

	// Stdout logger for general output
	stdoutLogger := log.With(
		log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", serviceID,
		"service.name", serviceName,
		"service.version", serviceVersion,
	)

	return &FileLoggers{
		InfoLogger:    infoLogger,
		WarningLogger: warningLogger,
		ErrorLogger:   errorLogger,
		StdoutLogger:  stdoutLogger,
	}, nil
}

// GetLogger returns appropriate logger based on level
func (fl *FileLoggers) GetLogger(level string) log.Logger {
	switch level {
	case "info":
		return fl.InfoLogger
	case "warning", "warn":
		return fl.WarningLogger
	case "error":
		return fl.ErrorLogger
	default:
		return fl.StdoutLogger
	}
}

