// Package ezzerolog provides the concrete, zerolog-based implementation of the logging service.
package ezzerolog

import (
	"bytes"
	"eaglechat/common/ezlog"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// rootLogger is the foundational, package-level logger.
// All other loggers are derived from it.
var rootLogger zerolog.Logger

// customWriter is an io.Writer that transforms zerolog's JSON output
// into the desired custom text format.
type customWriter struct {
	file io.Writer
}

// Write receives the JSON log output from zerolog, unmarshals it, and
// writes it in a custom format.
func (w *customWriter) Write(p []byte) (n int, err error) {
	var event map[string]any
	if err := json.Unmarshal(p, &event); err != nil {
		return 0, err
	}

	var b bytes.Buffer

	if level, ok := event["level"]; ok {
		b.WriteString(fmt.Sprintf("[%s] ", level))
	}
	if traceID, ok := event["trace_id"]; ok {
		b.WriteString(fmt.Sprintf("[%s] ", traceID))
	}
	if component, ok := event["component"]; ok {
		b.WriteString(fmt.Sprintf("[%s] ", component))
	}
	if msg, ok := event[zerolog.MessageFieldName]; ok {
		b.WriteString(fmt.Sprintf("%s", msg))
	}

	b.WriteByte('\n')

	_, err = w.file.Write(b.Bytes())
	return len(p), err
}

// init sets up the logging environment once when the package is loaded.
func init() {
	logFilePath := os.Getenv("LOGGER_PATH")
	if logFilePath == "" {
		logFilePath = "/data/ezlog"
	}

	os.MkdirAll(filepath.Dir(logFilePath), 0755)

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file at %s: %v", logFilePath, err))
	}

	// Create our custom writer instead of ConsoleWriter.
	writer := &customWriter{file: file}

	// The root logger is configured to write to our custom writer.
	rootLogger = zerolog.New(writer)
}

// zerologFactory implements the ezlog.Factory interface.
type zerologFactory struct{}

// NewFactory returns a new instance of the zerolog-based factory.
func NewFactory() ezlog.Factory {
	return &zerologFactory{}
}

// New creates a new Logger with a root component name and a new trace ID.
func (f *zerologFactory) New(component string) ezlog.Logger {
	logger := rootLogger.With().
		Str("trace_id", uuid.NewString()).
		Str("component", component).
		Logger()

	return &zerologAdapter{zlog: &logger}
}

// WithComponent creates a new Logger that inherits context (like trace_id)
// from an existing logger but uses a new component name.
func (f *zerologFactory) WithComponent(logger ezlog.Logger, component string) ezlog.Logger {
	if adapter, ok := logger.(*zerologAdapter); ok {
		newLogger := adapter.zlog.With().
			Str("component", component).
			Logger()
		return &zerologAdapter{zlog: &newLogger}
	}
	return logger
}

// zerologAdapter implements the ezlog.Logger interface by wrapping a zerolog.Logger.
type zerologAdapter struct {
	zlog *zerolog.Logger
}

func (l *zerologAdapter) Info(msg string) {
	l.zlog.Info().Msg(msg)
}

func (l *zerologAdapter) Infof(format string, v ...any) {
	l.zlog.Info().Msgf(format, v...)
}

func (l *zerologAdapter) Debug(msg string) {
	l.zlog.Debug().Msg(msg)
}

func (l *zerologAdapter) Debugf(format string, v ...any) {
	l.zlog.Debug().Msgf(format, v...)
}

func (l *zerologAdapter) Warn(msg string) {
	l.zlog.Warn().Msg(msg)
}

func (l *zerologAdapter) Warnf(format string, v ...any) {
	l.zlog.Warn().Msgf(format, v...)
}

func (l *zerologAdapter) Error(msg string) {
	l.zlog.Error().Msg(msg)
}

func (l *zerologAdapter) Errorf(format string, v ...any) {
	l.zlog.Error().Msgf(format, v...)
}

func (l *zerologAdapter) Fatal(msg string) {
	l.zlog.Fatal().Msg(msg)
}

func (l *zerologAdapter) Fatalf(format string, v ...any) {
	l.zlog.Fatal().Msgf(format, v...)
}
