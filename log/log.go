package log

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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

// CheckErr logs an error and optionally panics if the `shouldPanic` flag is set.
//
// This function is designed to handle errors consistently by logging them in a structured manner.
// It supports key-value pairs for additional context and has special handling for *http.Response
// values, dumping their contents for better debugging.
//
// Behavior:
// - If `err` is nil, the function returns false and does nothing.
// - If `err` is not nil, it logs the error using the structured logger (`zap.SugaredLogger`).
// - If `shouldPanic` is true, it logs the error as a panic and terminates the program.
// - If `keysAndValues` contains an *http.Response, it dumps the response body for debugging.
//
// Parameters:
// - `err` (error): The error to check and log. If nil, the function exits early.
// - `shouldPanic` (bool): If true, logs the error as a panic and terminates the program.
// - `message` (string): A descriptive error message for the log entry.
// - `keysAndValues` (variadic interface{}): Optional key-value pairs providing context.
//   - If a value is an *http.Response, the response body is dumped for easier debugging.
//   - If an odd number of arguments is provided, a placeholder `"<missing value>"` is added.
//
// Returns:
// - `bool`: Returns true if an error was logged, false otherwise.
func CheckErr(err error, shouldPanic bool, message string, keysAndValues ...interface{}) bool {
	if err == nil {
		return false
	}

	// Ensure keysAndValues has an even number of elements
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "<missing value>")
	}

	// Create a new slice to store processed key-value pairs
	newKeysAndValues := make([]interface{}, 0, len(keysAndValues)+2)

	// Prepend the error itself
	newKeysAndValues = append(newKeysAndValues, "error", err)

	// Process key-value pairs
	for i := 0; i < len(keysAndValues); i += 2 {
		key, value := keysAndValues[i], keysAndValues[i+1]

		if resp, ok := value.(*http.Response); ok {
			dump, dumpErr := httputil.DumpResponse(resp, true)
			if dumpErr != nil {
				GetLogger().sugaredLogger.Warnw("Failed to dump HTTP response", "error", dumpErr)
				newKeysAndValues = append(newKeysAndValues, key, fmt.Sprintf("Error dumping response: %v", dumpErr))
			} else {
				newKeysAndValues = append(newKeysAndValues, key, string(dump))
			}
		} else {
			newKeysAndValues = append(newKeysAndValues, key, value)
		}
	}

	// Log the error
	logger := GetLogger().sugaredLogger
	if shouldPanic {
		logger.Panicw(message, newKeysAndValues...)
	} else {
		logger.Errorw(message, newKeysAndValues...)
	}

	return true
}
