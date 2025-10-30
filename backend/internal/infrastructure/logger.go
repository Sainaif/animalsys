package infrastructure

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog logger
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new logger instance
func NewLogger(config *LogConfig) (*Logger, error) {
	// Set log level
	level, err := zerolog.ParseLevel(strings.ToLower(config.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output writers
	var writers []io.Writer

	// Console output (pretty format in development)
	if config.Format == "pretty" {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		writers = append(writers, consoleWriter)
	} else {
		// JSON format
		writers = append(writers, os.Stdout)
	}

	// File output (if configured)
	if config.File != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(config.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		writers = append(writers, file)
	}

	// Create multi-writer
	multiWriter := zerolog.MultiLevelWriter(writers...)

	// Create logger with context
	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set as global logger
	log.Logger = logger

	return &Logger{
		logger: logger,
	}, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.logger.Debug()
	if len(fields) > 0 {
		event = addFields(event, fields[0])
	}
	event.Msg(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.logger.Info()
	if len(fields) > 0 {
		event = addFields(event, fields[0])
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.logger.Warn()
	if len(fields) > 0 {
		event = addFields(event, fields[0])
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Error().Err(err)
	if len(fields) > 0 {
		event = addFields(event, fields[0])
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Fatal().Err(err)
	if len(fields) > 0 {
		event = addFields(event, fields[0])
	}
	event.Msg(msg)
}

// With creates a child logger with additional fields
func (l *Logger) With(fields map[string]interface{}) *Logger {
	ctx := l.logger.With()
	for key, value := range fields {
		ctx = ctx.Interface(key, value)
	}
	return &Logger{
		logger: ctx.Logger(),
	}
}

// HTTP creates a logger for HTTP requests
func (l *Logger) HTTP() *Logger {
	return &Logger{
		logger: l.logger.With().Str("component", "http").Logger(),
	}
}

// Database creates a logger for database operations
func (l *Logger) Database() *Logger {
	return &Logger{
		logger: l.logger.With().Str("component", "database").Logger(),
	}
}

// Auth creates a logger for authentication operations
func (l *Logger) Auth() *Logger {
	return &Logger{
		logger: l.logger.With().Str("component", "auth").Logger(),
	}
}

// addFields adds fields to a log event
func addFields(event *zerolog.Event, fields map[string]interface{}) *zerolog.Event {
	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case int64:
			event = event.Int64(key, v)
		case float64:
			event = event.Float64(key, v)
		case bool:
			event = event.Bool(key, v)
		case time.Time:
			event = event.Time(key, v)
		case error:
			event = event.AnErr(key, v)
		default:
			event = event.Interface(key, v)
		}
	}
	return event
}

// GetZerolog returns the underlying zerolog logger
func (l *Logger) GetZerolog() zerolog.Logger {
	return l.logger
}
