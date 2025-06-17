// Package unit provides unit tests for the logger functionality.
package unit

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected internal.LogLevel
	}{
		{
			name:     "debug level",
			input:    "debug",
			expected: internal.DebugLevel,
		},
		{
			name:     "debug level uppercase",
			input:    "DEBUG",
			expected: internal.DebugLevel,
		},
		{
			name:     "info level",
			input:    "info",
			expected: internal.InfoLevel,
		},
		{
			name:     "info level mixed case",
			input:    "InFo",
			expected: internal.InfoLevel,
		},
		{
			name:     "warn level",
			input:    "warn",
			expected: internal.WarnLevel,
		},
		{
			name:     "warning level",
			input:    "warning",
			expected: internal.WarnLevel,
		},
		{
			name:     "error level",
			input:    "error",
			expected: internal.ErrorLevel,
		},
		{
			name:     "silent level",
			input:    "silent",
			expected: internal.SilentLevel,
		},
		{
			name:     "none level",
			input:    "none",
			expected: internal.SilentLevel,
		},
		{
			name:     "invalid level defaults to info",
			input:    "invalid",
			expected: internal.InfoLevel,
		},
		{
			name:     "empty string defaults to info",
			input:    "",
			expected: internal.InfoLevel,
		},
		{
			name:     "whitespace defaults to info",
			input:    "   ",
			expected: internal.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ParseLogLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    internal.LogLevel
		expected internal.LogLevel
	}{
		{
			name:     "debug level logger",
			level:    internal.DebugLevel,
			expected: internal.DebugLevel,
		},
		{
			name:     "info level logger",
			level:    internal.InfoLevel,
			expected: internal.InfoLevel,
		},
		{
			name:     "warn level logger",
			level:    internal.WarnLevel,
			expected: internal.WarnLevel,
		},
		{
			name:     "error level logger",
			level:    internal.ErrorLevel,
			expected: internal.ErrorLevel,
		},
		{
			name:     "silent level logger",
			level:    internal.SilentLevel,
			expected: internal.SilentLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := internal.NewLogger(tt.level)
			require.NotNil(t, logger)
			// We can't directly access the level field, but we can test behavior
			// This will be tested in the logging method tests
		})
	}
}

func TestLoggerDebug(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       internal.LogLevel
		message        string
		args           []any
		expectOutput   bool
		expectedPrefix string
	}{
		{
			name:           "debug level should output debug messages",
			logLevel:       internal.DebugLevel,
			message:        "debug message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[DEBUG]",
		},
		{
			name:           "debug level with args",
			logLevel:       internal.DebugLevel,
			message:        "debug message with %s and %d",
			args:           []any{"string", 42},
			expectOutput:   true,
			expectedPrefix: "[DEBUG]",
		},
		{
			name:         "info level should not output debug messages",
			logLevel:     internal.InfoLevel,
			message:      "debug message",
			args:         nil,
			expectOutput: false,
		},
		{
			name:         "warn level should not output debug messages",
			logLevel:     internal.WarnLevel,
			message:      "debug message",
			args:         nil,
			expectOutput: false,
		},
		{
			name:         "error level should not output debug messages",
			logLevel:     internal.ErrorLevel,
			message:      "debug message",
			args:         nil,
			expectOutput: false,
		},
		{
			name:         "silent level should not output debug messages",
			logLevel:     internal.SilentLevel,
			message:      "debug message",
			args:         nil,
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			originalStderr := os.Stderr
			r, w, err := os.Pipe()
			require.NoError(t, err)
			os.Stderr = w

			logger := internal.NewLogger(tt.logLevel)
			logger.Debug(tt.message, tt.args...)

			// Close write end and restore stderr
			_ = w.Close()
			os.Stderr = originalStderr

			// Read captured output
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			require.NoError(t, err)
			output := buf.String()

			if tt.expectOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "debug message")
				if len(tt.args) > 0 {
					assert.Contains(t, output, "string")
					assert.Contains(t, output, "42")
				}
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLoggerInfo(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       internal.LogLevel
		message        string
		args           []any
		expectOutput   bool
		expectedPrefix string
	}{
		{
			name:           "debug level should output info messages",
			logLevel:       internal.DebugLevel,
			message:        "info message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[INFO]",
		},
		{
			name:           "info level should output info messages",
			logLevel:       internal.InfoLevel,
			message:        "info message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[INFO]",
		},
		{
			name:           "info level with formatting",
			logLevel:       internal.InfoLevel,
			message:        "info message: %s=%d",
			args:           []any{"count", 10},
			expectOutput:   true,
			expectedPrefix: "[INFO]",
		},
		{
			name:         "warn level should not output info messages",
			logLevel:     internal.WarnLevel,
			message:      "info message",
			args:         nil,
			expectOutput: false,
		},
		{
			name:         "error level should not output info messages",
			logLevel:     internal.ErrorLevel,
			message:      "info message",
			args:         nil,
			expectOutput: false,
		},
		{
			name:         "silent level should not output info messages",
			logLevel:     internal.SilentLevel,
			message:      "info message",
			args:         nil,
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := helpers.CaptureStderr(t, func() {
				logger := internal.NewLogger(tt.logLevel)
				logger.Info(tt.message, tt.args...)
			})

			if tt.expectOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "info message")
				if len(tt.args) > 0 {
					assert.Contains(t, output, "count=10")
				}
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLoggerWarn(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       internal.LogLevel
		message        string
		args           []any
		expectOutput   bool
		expectedPrefix string
	}{
		{
			name:           "debug level should output warn messages",
			logLevel:       internal.DebugLevel,
			message:        "warn message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[WARN]",
		},
		{
			name:           "info level should output warn messages",
			logLevel:       internal.InfoLevel,
			message:        "warn message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[WARN]",
		},
		{
			name:           "warn level should output warn messages",
			logLevel:       internal.WarnLevel,
			message:        "warn message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[WARN]",
		},
		{
			name:           "warn with complex formatting",
			logLevel:       internal.WarnLevel,
			message:        "warning: failed to process %s after %d attempts",
			args:           []any{"file.txt", 3},
			expectOutput:   true,
			expectedPrefix: "[WARN]",
		},
		{
			name:         "error level should not output warn messages",
			logLevel:     internal.ErrorLevel,
			message:      "warn message",
			args:         nil,
			expectOutput: false,
		},
		{
			name:         "silent level should not output warn messages",
			logLevel:     internal.SilentLevel,
			message:      "warn message",
			args:         nil,
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := helpers.CaptureStderr(t, func() {
				logger := internal.NewLogger(tt.logLevel)
				logger.Warn(tt.message, tt.args...)
			})

			if tt.expectOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "warn")
				if len(tt.args) > 0 {
					assert.Contains(t, output, "file.txt")
					assert.Contains(t, output, "3")
				}
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLoggerError(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       internal.LogLevel
		message        string
		args           []any
		expectOutput   bool
		expectedPrefix string
	}{
		{
			name:           "debug level should output error messages",
			logLevel:       internal.DebugLevel,
			message:        "error message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "info level should output error messages",
			logLevel:       internal.InfoLevel,
			message:        "error message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "warn level should output error messages",
			logLevel:       internal.WarnLevel,
			message:        "error message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "error level should output error messages",
			logLevel:       internal.ErrorLevel,
			message:        "error message",
			args:           nil,
			expectOutput:   true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:           "error with detailed formatting",
			logLevel:       internal.ErrorLevel,
			message:        "API error: %s (status %d)",
			args:           []any{"authentication failed", 401},
			expectOutput:   true,
			expectedPrefix: "[ERROR]",
		},
		{
			name:         "silent level should not output error messages",
			logLevel:     internal.SilentLevel,
			message:      "error message",
			args:         nil,
			expectOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := helpers.CaptureStderr(t, func() {
				logger := internal.NewLogger(tt.logLevel)
				logger.Error(tt.message, tt.args...)
			})

			if tt.expectOutput {
				assert.Contains(t, output, tt.expectedPrefix)
				assert.Contains(t, output, "error")
				if len(tt.args) > 0 {
					assert.Contains(t, output, "authentication failed")
					assert.Contains(t, output, "401")
				}
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLoggerMessageFormatting(t *testing.T) {
	tests := []struct {
		name             string
		logLevel         internal.LogLevel
		logMethod        func(*internal.Logger, string, ...any)
		message          string
		args             []any
		expectedInOutput []string
	}{
		{
			name:     "simple message without args",
			logLevel: internal.InfoLevel,
			logMethod: func(l *internal.Logger, msg string, args ...any) {
				l.Info(msg, args...)
			},
			message:          "simple message",
			args:             nil,
			expectedInOutput: []string{"[INFO]", "simple message"},
		},
		{
			name:     "message with string formatting",
			logLevel: internal.InfoLevel,
			logMethod: func(l *internal.Logger, msg string, args ...any) {
				l.Info(msg, args...)
			},
			message:          "user %s logged in",
			args:             []any{"john"},
			expectedInOutput: []string{"[INFO]", "user john logged in"},
		},
		{
			name:     "message with multiple format specifiers",
			logLevel: internal.DebugLevel,
			logMethod: func(l *internal.Logger, msg string, args ...any) {
				l.Debug(msg, args...)
			},
			message:          "processed %d items in %s with %f%% success rate",
			args:             []any{42, "2.5s", 98.7},
			expectedInOutput: []string{"[DEBUG]", "processed 42 items", "2.5s", "98.7"},
		},
		{
			name:     "message with special characters",
			logLevel: internal.WarnLevel,
			logMethod: func(l *internal.Logger, msg string, args ...any) {
				l.Warn(msg, args...)
			},
			message:          "special chars: %s contains symbols @#$%%",
			args:             []any{"file.txt"},
			expectedInOutput: []string{"[WARN]", "file.txt", "@#$%"},
		},
		{
			name:     "empty message with args",
			logLevel: internal.ErrorLevel,
			logMethod: func(l *internal.Logger, msg string, args ...any) {
				l.Error(msg, args...)
			},
			message:          "",
			args:             []any{"ignored"},
			expectedInOutput: []string{"[ERROR]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := helpers.CaptureStderr(t, func() {
				logger := internal.NewLogger(tt.logLevel)
				tt.logMethod(logger, tt.message, tt.args...)
			})

			for _, expected := range tt.expectedInOutput {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestLoggerConcurrency(t *testing.T) {
	// Test that logger is safe for concurrent use
	const numGoroutines = 10
	const messagesPerGoroutine = 10

	logger := internal.NewLogger(internal.InfoLevel)

	output := helpers.CaptureStderr(t, func() {
		helpers.RunConcurrently(t, numGoroutines, func(goroutineID int) {
			for i := range messagesPerGoroutine {
				logger.Info("goroutine %d message %d", goroutineID, i)
				logger.Debug("debug from goroutine %d", goroutineID)
				logger.Warn("warning from goroutine %d", goroutineID)
				logger.Error("error from goroutine %d", goroutineID)
			}
		})
	})

	// Count expected messages (only Info, Warn, Error should appear with InfoLevel)
	expectedInfoMessages := numGoroutines * messagesPerGoroutine
	expectedWarnMessages := numGoroutines * messagesPerGoroutine
	expectedErrorMessages := numGoroutines * messagesPerGoroutine

	infoCount := strings.Count(output, "[INFO]")
	warnCount := strings.Count(output, "[WARN]")
	errorCount := strings.Count(output, "[ERROR]")
	debugCount := strings.Count(output, "[DEBUG]")

	assert.Equal(
		t,
		expectedInfoMessages,
		infoCount,
		"Expected %d info messages, got %d",
		expectedInfoMessages,
		infoCount,
	)
	assert.Equal(
		t,
		expectedWarnMessages,
		warnCount,
		"Expected %d warn messages, got %d",
		expectedWarnMessages,
		warnCount,
	)
	assert.Equal(
		t,
		expectedErrorMessages,
		errorCount,
		"Expected %d error messages, got %d",
		expectedErrorMessages,
		errorCount,
	)
	assert.Equal(t, 0, debugCount, "Expected 0 debug messages with InfoLevel, got %d", debugCount)
}

func TestLoggerEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  internal.LogLevel
		testFunc  func(*internal.Logger)
		expectLog bool
	}{
		{
			name:     "nil args slice",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				l.Info("message with nil args", nil...)
			},
			expectLog: true,
		},
		{
			name:     "empty args slice",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				l.Info("message with empty args", []any{}...)
			},
			expectLog: true,
		},
		{
			name:     "message with newlines",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				l.Info("message\nwith\nnewlines")
			},
			expectLog: true,
		},
		{
			name:     "very long message",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				longMsg := strings.Repeat("a", 1000)
				l.Info("long message: %s", longMsg)
			},
			expectLog: true,
		},
		{
			name:     "message with unicode",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				l.Info("unicode message: 🚀 测试 العربية")
			},
			expectLog: true,
		},
		{
			name:     "format mismatch - too few args",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				l.Info("format %s %d %s", "only", 1)
			},
			expectLog: true,
		},
		{
			name:     "format mismatch - too many args",
			logLevel: internal.InfoLevel,
			testFunc: func(l *internal.Logger) {
				l.Info("format %s", "arg1", "arg2", "arg3")
			},
			expectLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := helpers.CaptureStderr(t, func() {
				logger := internal.NewLogger(tt.logLevel)
				tt.testFunc(logger)
			})

			if tt.expectLog {
				assert.NotEmpty(t, output, "Expected log output but got none")
				assert.Contains(t, output, "[INFO]")
			} else {
				assert.Empty(t, output, "Expected no log output but got: %s", output)
			}
		})
	}
}

func TestLoggerLevelHierarchy(t *testing.T) {
	// Test that log levels properly filter messages
	levels := []internal.LogLevel{
		internal.DebugLevel,
		internal.InfoLevel,
		internal.WarnLevel,
		internal.ErrorLevel,
		internal.SilentLevel,
	}

	levelNames := []string{"DebugLevel", "InfoLevel", "WarnLevel", "ErrorLevel", "SilentLevel"}

	for i, level := range levels {
		t.Run("level_"+levelNames[i], func(t *testing.T) {
			output := helpers.CaptureStderr(t, func() {
				logger := internal.NewLogger(level)
				logger.Debug("debug message")
				logger.Info("info message")
				logger.Warn("warn message")
				logger.Error("error message")
			})

			debugCount := strings.Count(output, "[DEBUG]")
			infoCount := strings.Count(output, "[INFO]")
			warnCount := strings.Count(output, "[WARN]")
			errorCount := strings.Count(output, "[ERROR]")

			switch level {
			case internal.DebugLevel:
				assert.Equal(t, 1, debugCount)
				assert.Equal(t, 1, infoCount)
				assert.Equal(t, 1, warnCount)
				assert.Equal(t, 1, errorCount)
			case internal.InfoLevel:
				assert.Equal(t, 0, debugCount)
				assert.Equal(t, 1, infoCount)
				assert.Equal(t, 1, warnCount)
				assert.Equal(t, 1, errorCount)
			case internal.WarnLevel:
				assert.Equal(t, 0, debugCount)
				assert.Equal(t, 0, infoCount)
				assert.Equal(t, 1, warnCount)
				assert.Equal(t, 1, errorCount)
			case internal.ErrorLevel:
				assert.Equal(t, 0, debugCount)
				assert.Equal(t, 0, infoCount)
				assert.Equal(t, 0, warnCount)
				assert.Equal(t, 1, errorCount)
			case internal.SilentLevel:
				assert.Equal(t, 0, debugCount)
				assert.Equal(t, 0, infoCount)
				assert.Equal(t, 0, warnCount)
				assert.Equal(t, 0, errorCount)
			}
		})
	}
}
