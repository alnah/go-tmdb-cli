// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"os"
	"strings"
)

// LogLevel represents the logging level for the application.
type LogLevel int

// Available log levels.
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	SilentLevel
)

// ParseLogLevel parses a string log level into LogLevel.
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "silent", "none":
		return SilentLevel
	default:
		return InfoLevel
	}
}

// Logger provides structured logging for the application.
type Logger struct {
	level LogLevel
}

// NewLogger creates a new logger with the specified level.
func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...any) {
	if l.level <= DebugLevel {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+msg+"\n", args...)
	}
}

// Info logs an info message.
func (l *Logger) Info(msg string, args ...any) {
	if l.level <= InfoLevel {
		fmt.Fprintf(os.Stderr, "[INFO] "+msg+"\n", args...)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...any) {
	if l.level <= WarnLevel {
		fmt.Fprintf(os.Stderr, "[WARN] "+msg+"\n", args...)
	}
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...any) {
	if l.level <= ErrorLevel {
		fmt.Fprintf(os.Stderr, "[ERROR] "+msg+"\n", args...)
	}
}
