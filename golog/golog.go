package golog

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	logger *Logger
)

type Logger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewLogger(isDevelopment bool) (*Logger, error) {
	var config zap.Config

	if isDevelopment {
		config = zap.NewDevelopmentConfig() // Defaults to DebugLevel
		config.Encoding = "console"         // Pretty-print for dev mode
	} else {
		config = zap.NewProductionConfig() // Defaults to InfoLevel
		config.Encoding = "json"           // JSON for production
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{sugaredLogger: zapLogger.Sugar()}, nil
}

func InitLogger(isDevelopment bool) {
	once.Do(func() {
		var err error
		logger, err = NewLogger(isDevelopment)
		if err != nil {
			panic("failed to initialize logger")
		}
	})
}

func GetLogger() *Logger {
	return logger
}

// Info logs an info message with key-value pairs.
func Info(msg string, keysAndValues ...interface{}) {
	GetLogger().sugaredLogger.Infow(msg, keysAndValues...)
}

// Debug logs a debug message with key-value pairs.
func Debug(msg string, keysAndValues ...interface{}) {
	GetLogger().sugaredLogger.Debugw(msg, keysAndValues...)
}

// Warn logs a warning message with key-value pairs.
func Warn(msg string, keysAndValues ...interface{}) {
	GetLogger().sugaredLogger.Warnw(msg, keysAndValues...)
}

// Error logs an error message with key-value pairs.
func Error(msg string, keysAndValues ...interface{}) {
	GetLogger().sugaredLogger.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message with key-value pairs and terminates the application.
func Fatal(msg string, keysAndValues ...interface{}) {
	GetLogger().sugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panic logs a panic message with key-value pairs and panics the application.
func Panic(msg string, keysAndValues ...interface{}) {
	GetLogger().sugaredLogger.Panicw(msg, keysAndValues...)
}

// Debugf logs a debug message with formatted text.
func Debugf(template string, args ...interface{}) {
	GetLogger().sugaredLogger.Debugf(template, args...)
}

// Infof logs an info message with formatted text.
func Infof(template string, args ...interface{}) {
	GetLogger().sugaredLogger.Infof(template, args...)
}

// Warnf logs a warning message with formatted text.
func Warnf(template string, args ...interface{}) {
	GetLogger().sugaredLogger.Warnf(template, args...)
}

// Errorf logs an error message with formatted text.
func Errorf(template string, args ...interface{}) {
	GetLogger().sugaredLogger.Errorf(template, args...)
}

// Fatalf logs a fatal message with formatted text and terminates the application.
func Fatalf(template string, args ...interface{}) {
	GetLogger().sugaredLogger.Fatalf(template, args...)
}

// Panicf logs a panic message with formatted text and panics the application.
func Panicf(template string, args ...interface{}) {
	GetLogger().sugaredLogger.Panicf(template, args...)
}

// CheckErr checks if an error is nil. If not, it logs it and optionally exits the program.
func CheckErr(err error, message string, panic bool) bool {
	if err != nil {
		if panic {
			GetLogger().sugaredLogger.DPanicf("%s: %v", message, err)
		}
		GetLogger().sugaredLogger.Errorf("%s: %v", message, err)
		return true
	}
	return false
}
