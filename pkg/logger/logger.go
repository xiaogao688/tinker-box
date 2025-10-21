package logger

import (
	"os"
	"path/filepath"

	"yourapp/internal/global"
	appconfig "yourapp/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Init initializes the logger
func Init() error {
	cfg := global.GetConfig()
	if cfg == nil {
		// Use default configuration if global config is not set
		return initDefaultLogger()
	}

	return initLoggerWithConfig(cfg)
}

// initDefaultLogger initializes logger with default configuration
func initDefaultLogger() error {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}

	sugar = logger.Sugar()
	return nil
}

// initLoggerWithConfig initializes logger with application configuration
func initLoggerWithConfig(cfg *appconfig.Config) error {
	// Set log level
	level, err := parseLogLevel(cfg.Logging.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Set encoder based on format
	var encoder zapcore.Encoder
	switch cfg.Logging.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	var writeSyncer zapcore.WriteSyncer
	switch cfg.Logging.Output {
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	case "file":
		if cfg.Logging.FilePath != "" {
			if err := setupFileOutput(cfg.Logging.FilePath); err != nil {
				return err
			}
			file, err := os.OpenFile(cfg.Logging.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			writeSyncer = zapcore.AddSync(file)
		} else {
			writeSyncer = zapcore.AddSync(os.Stdout)
		}
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = logger.Sugar()

	return nil
}

// setupFileOutput sets up file output for logging
func setupFileOutput(filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

// parseLogLevel parses log level string to zapcore.Level
func parseLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		// Initialize with default settings if not already initialized
		initDefaultLogger()
	}
	return logger
}

// GetSugar returns the sugared logger instance
func GetSugar() *zap.SugaredLogger {
	if sugar == nil {
		// Initialize with default settings if not already initialized
		initDefaultLogger()
	}
	return sugar
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Infof logs a formatted info message
func Infof(template string, args ...interface{}) {
	GetSugar().Infof(template, args...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Errorf logs a formatted error message
func Errorf(template string, args ...interface{}) {
	GetSugar().Errorf(template, args...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Debugf logs a formatted debug message
func Debugf(template string, args ...interface{}) {
	GetSugar().Debugf(template, args...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Warnf logs a formatted warning message
func Warnf(template string, args ...interface{}) {
	GetSugar().Warnf(template, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(template string, args ...interface{}) {
	GetSugar().Fatalf(template, args...)
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// WithSugar creates a child sugared logger with the given fields
func WithSugar(args ...interface{}) *zap.SugaredLogger {
	return GetSugar().With(args...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}
