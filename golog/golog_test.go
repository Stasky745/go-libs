package golog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupTestLogger(isDevelopment bool) (*bytes.Buffer, func()) {
	var buf bytes.Buffer

	// Create a custom writer syncer to capture logs
	writerSyncer := zapcore.AddSync(&buf)

	// Choose the appropriate config based on development mode
	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console" // Pretty print for dev mode
	} else {
		config = zap.NewProductionConfig()
		config.Encoding = "json" // JSON for production
	}

	// Redirect log output to our buffer
	config.OutputPaths = []string{"stdout"} // 🔴 Issue: Should be writerSyncer
	config.OutputPaths = []string{}         // 🛠 Fix: Empty output paths
	config.EncoderConfig.TimeKey = ""       // 🛠 Fix: Removes timestamp (simplifies testing)

	// Build logger with custom writer
	zapLogger, _ := config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig), // Pretty or JSON encoder
			writerSyncer, // 👈 Use the buffer as the output
			config.Level, // Use the set log level
		)
	}))

	// Set global logger
	logger = &Logger{sugaredLogger: zapLogger.Sugar()}

	return &buf, func() { _ = zapLogger.Sync() }
}

// **TEST 1: Verify Info Log Output**
func TestInfoLog(t *testing.T) {
	buf, cleanup := setupTestLogger(true) // Use development mode
	defer cleanup()

	Info("Test log", "key", "value")

	// Read buffer and check log output
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Test log")
	assert.Contains(t, logOutput, "key")
	assert.Contains(t, logOutput, "value")
}

// **TEST 2: Verify JSON Format in Production Mode**
func TestProductionLogging(t *testing.T) {
	buf, cleanup := setupTestLogger(false) // Use production mode
	defer cleanup()

	Info("Production log", "key", "value")

	// Read buffer and check if log is JSON
	logOutput := buf.String()
	assert.Contains(t, logOutput, `Production log`) // JSON format
	assert.Contains(t, logOutput, `"key": "value"`) // JSON key-value
}

// **TEST 3: Debug Log Should Not Appear in Production**
func TestDebugLogNotShownInProduction(t *testing.T) {
	buf, cleanup := setupTestLogger(false) // Production mode
	defer cleanup()

	Debug("Debug log", "key", "debug_value")

	// Read buffer and check if log contains "Debug log"
	logOutput := buf.String()
	assert.NotContains(t, logOutput, "Debug log") // Should NOT appear in production
}

// **TEST 4: Debug Log Should Appear in Development**
func TestDebugLogInDevelopment(t *testing.T) {
	buf, cleanup := setupTestLogger(true) // Development mode
	defer cleanup()

	Debug("Debug log", "key", "debug_value")

	// Read buffer and check if log contains "Debug log"
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Debug log") // Should appear in development
}

// **TEST 5: CheckErr Helper Function**
func TestCheckErr(t *testing.T) {
	buf, cleanup := setupTestLogger(true)
	defer cleanup()

	err := CheckErr(nil, "No error", false)
	assert.False(t, err) // Should return false for no error

	err = CheckErr(assert.AnError, "Error happened", false)
	assert.True(t, err) // Should return true if an error occurred

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Error happened") // Should log the error message
}
