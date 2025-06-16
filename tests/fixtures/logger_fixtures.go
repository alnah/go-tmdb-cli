// Package fixtures provides test data and fixtures for the TMDB CLI tests.
package fixtures

import "github.com/alnah/tmdb-cli/internal"

// LogLevelTestData contains test data for log level parsing tests.
var LogLevelTestData = map[string]internal.LogLevel{
	"debug":   internal.DebugLevel,
	"DEBUG":   internal.DebugLevel,
	"Debug":   internal.DebugLevel,
	"info":    internal.InfoLevel,
	"INFO":    internal.InfoLevel,
	"Info":    internal.InfoLevel,
	"warn":    internal.WarnLevel,
	"WARN":    internal.WarnLevel,
	"Warn":    internal.WarnLevel,
	"warning": internal.WarnLevel,
	"WARNING": internal.WarnLevel,
	"Warning": internal.WarnLevel,
	"error":   internal.ErrorLevel,
	"ERROR":   internal.ErrorLevel,
	"Error":   internal.ErrorLevel,
	"silent":  internal.SilentLevel,
	"SILENT":  internal.SilentLevel,
	"Silent":  internal.SilentLevel,
	"none":    internal.SilentLevel,
	"NONE":    internal.SilentLevel,
	"None":    internal.SilentLevel,
}

// InvalidLogLevels contains invalid log level strings that should default to InfoLevel.
var InvalidLogLevels = []string{
	"",
	"   ",
	"invalid",
	"unknown",
	"trace",
	"fatal",
	"all",
	"off",
	"123",
	"debug-level",
	"info level",
	"warn_level",
	"error.level",
}

// LoggerTestMessages contains sample messages for testing different log scenarios.
var LoggerTestMessages = struct {
	Simple           string
	WithFormatting   string
	WithMultipleArgs string
	WithSpecialChars string
	WithNewlines     string
	WithUnicode      string
	Empty            string
	VeryLong         string
	WithJSONLike     string
	WithURLs         string
	WithStackTrace   string
	WithHTTPStatus   string
	WithTimestamp    string
	WithUserData     string
}{
	Simple:           "This is a simple log message",
	WithFormatting:   "User %s performed action %s at %s",
	WithMultipleArgs: "Processed %d items in %s with %.2f%% success rate",
	WithSpecialChars: "Special characters: @#$%^&*()_+-={}[]|\\:;\"'<>,.?/~`",
	WithNewlines:     "Message with\nmultiple\nlines",
	WithUnicode:      "Unicode message: 🚀 测试 العربية ñáéíóú",
	Empty:            "",
	VeryLong: "This is a very long message that exceeds normal logging" +
		" message lengths and is used to test how the logger handles extremely" +
		" long content that might be encountered in real-world scenarios such" +
		"as when logging large API responses or detailed error information",
	WithJSONLike:   `Log message with JSON-like content: {"key": "value", "number": 123, "boolean": true}`,
	WithURLs:       "Failed to connect to https://api.themoviedb.org/3/movie/popular?api_key=hidden",
	WithStackTrace: "Error occurred at line 42 in file.go: function() -> subfunction() -> error()",
	WithHTTPStatus: "HTTP request failed with status %d: %s",
	WithTimestamp:  "Event occurred at %s with correlation ID %s",
	WithUserData:   "User %s (ID: %d) from IP %s performed %s operation",
}

// LoggerTestArgs contains sample arguments for testing message formatting.
var LoggerTestArgs = struct {
	StringArgs   []any
	NumberArgs   []any
	MixedArgs    []any
	EmptyArgs    []any
	NilArgs      []any
	BooleanArgs  []any
	FloatArgs    []any
	ComplexArgs  []any
	MismatchArgs []any
}{
	StringArgs:   []any{"john_doe", "login", "2024-01-15T10:30:00Z"},
	NumberArgs:   []any{42, 100, 98.75},
	MixedArgs:    []any{"user123", 42, true, 3.14, "success"},
	EmptyArgs:    []any{},
	NilArgs:      nil,
	BooleanArgs:  []any{true, false, true},
	FloatArgs:    []any{3.14159, 2.71828, 1.41421},
	ComplexArgs:  []any{map[string]int{"count": 5}, []string{"a", "b", "c"}},
	MismatchArgs: []any{"too", "many", "args", "for", "format"},
}

// LoggerExpectedOutputs contains expected output patterns for different log levels.
var LoggerExpectedOutputs = struct {
	DebugPrefix string
	InfoPrefix  string
	WarnPrefix  string
	ErrorPrefix string
}{
	DebugPrefix: "[DEBUG]",
	InfoPrefix:  "[INFO]",
	WarnPrefix:  "[WARN]",
	ErrorPrefix: "[ERROR]",
}

// LoggerConcurrencyTestData contains test data for concurrent logging tests.
var LoggerConcurrencyTestData = struct {
	NumGoroutines        int
	MessagesPerGoroutine int
	TestMessage          string
}{
	NumGoroutines:        20,
	MessagesPerGoroutine: 50,
	TestMessage:          "Concurrent message from goroutine %d, iteration %d",
}

// LoggerEdgeCaseTestData contains test data for edge cases and error conditions.
var LoggerEdgeCaseTestData = struct {
	MalformedFormat   string
	TooFewArgs        []any
	TooManyArgs       []any
	NullByte          string
	ControlCharacters string
	MaxLengthMessage  string
}{
	MalformedFormat:   "Format with wrong specifiers %s %d %s",
	TooFewArgs:        []any{"only_one"},
	TooManyArgs:       []any{"arg1", "arg2", "arg3", "arg4", "arg5"},
	NullByte:          "Message with null byte: \x00",
	ControlCharacters: "Message with control chars: \t\r\n\b\f",
	MaxLengthMessage:  string(make([]byte, 10000)), // 10KB message
}

// GetLogLevelTestCases returns test cases for log level functionality.
func GetLogLevelTestCases() []struct {
	Name           string
	LogLevel       internal.LogLevel
	ShouldLogDebug bool
	ShouldLogInfo  bool
	ShouldLogWarn  bool
	ShouldLogError bool
} {
	return []struct {
		Name           string
		LogLevel       internal.LogLevel
		ShouldLogDebug bool
		ShouldLogInfo  bool
		ShouldLogWarn  bool
		ShouldLogError bool
	}{
		{
			Name:           "DebugLevel logs all messages",
			LogLevel:       internal.DebugLevel,
			ShouldLogDebug: true,
			ShouldLogInfo:  true,
			ShouldLogWarn:  true,
			ShouldLogError: true,
		},
		{
			Name:           "InfoLevel logs info, warn, error",
			LogLevel:       internal.InfoLevel,
			ShouldLogDebug: false,
			ShouldLogInfo:  true,
			ShouldLogWarn:  true,
			ShouldLogError: true,
		},
		{
			Name:           "WarnLevel logs warn, error",
			LogLevel:       internal.WarnLevel,
			ShouldLogDebug: false,
			ShouldLogInfo:  false,
			ShouldLogWarn:  true,
			ShouldLogError: true,
		},
		{
			Name:           "ErrorLevel logs only error",
			LogLevel:       internal.ErrorLevel,
			ShouldLogDebug: false,
			ShouldLogInfo:  false,
			ShouldLogWarn:  false,
			ShouldLogError: true,
		},
		{
			Name:           "SilentLevel logs nothing",
			LogLevel:       internal.SilentLevel,
			ShouldLogDebug: false,
			ShouldLogInfo:  false,
			ShouldLogWarn:  false,
			ShouldLogError: false,
		},
	}
}

// GetFormattingTestCases returns test cases for message formatting.
func GetFormattingTestCases() []struct {
	Name          string
	Message       string
	Args          []any
	ExpectedParts []string
} {
	return []struct {
		Name          string
		Message       string
		Args          []any
		ExpectedParts []string
	}{
		{
			Name:          "Simple string formatting",
			Message:       "Hello %s",
			Args:          []any{"world"},
			ExpectedParts: []string{"Hello world"},
		},
		{
			Name:          "Multiple argument types",
			Message:       "User %s has %d points (%.1f%% completion)",
			Args:          []any{"alice", 150, 75.5},
			ExpectedParts: []string{"User alice", "150 points", "75.5%"},
		},
		{
			Name:          "Boolean formatting",
			Message:       "Success: %t, Failed: %t",
			Args:          []any{true, false},
			ExpectedParts: []string{"Success: true", "Failed: false"},
		},
		{
			Name:          "Hexadecimal formatting",
			Message:       "Memory address: %x, Value: %X",
			Args:          []any{255, 255},
			ExpectedParts: []string{"ff", "FF"},
		},
		{
			Name:          "No formatting placeholders",
			Message:       "Static message without formatting",
			Args:          []any{"ignored", "args"},
			ExpectedParts: []string{"Static message without formatting"},
		},
	}
}
