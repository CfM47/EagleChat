package ezlog

import (
	"context"
	"fmt"
	"log"
)

// Logger defines a standardized logging interface for the application.
// This allows the logging implementation to be swapped out without affecting
// the application code that uses it.
type Logger interface {
	// Info logs an informational message.
	Info(msg string)
	// Infof logs a formatted informational message.
	Infof(format string, v ...any)
	// Debug logs a debug message.
	Debug(msg string)
	// Debugf logs a formatted debug message.
	Debugf(format string, v ...any)
	// Warn logs a warning message.
	Warn(msg string)
	// Warnf logs a formatted warning message.
	Warnf(format string, v ...any)
	// Error logs an error message.
	Error(msg string)
	// Errorf logs a formatted error message.
	Errorf(format string, v ...any)
	// Fatal logs a fatal error message and exits the application.
	Fatal(msg string)
	// Fatalf logs a formatted fatal error message and exits the application.
	Fatalf(format string, v ...any)
}

// Factory defines the interface for creating Logger instances.
// This abstracts the specific logger creation logic from the domain.
type Factory interface {
	// New creates a new Logger with a root component name and a new trace ID.
	New(component string) Logger
	// WithComponent creates a new Logger that inherits the trace ID from an
	// existing logger but has a new component name.
	WithComponent(logger Logger, component string) Logger
}

// loggerFactory holds the concrete implementation of the Factory interface.
// It is set at application startup.
var loggerFactory Factory

// SetLoggerFactory injects the concrete logging factory implementation.
// This must be called once at application startup.
func SetLoggerFactory(factory Factory) {
	if factory == nil {
		panic("logger factory cannot be nil")
	}
	if loggerFactory != nil {
		panic("logger factory already set")
	}
	loggerFactory = factory
}

// loggerKey is an unexported type to be used as the key for storing the
// logger in a context.Context.
type loggerKey string

const loggerCtxKey loggerKey = "logger"

// NewLoggerContext creates a new context with a logger instance for a given
// component. This is the entry point for creating a new traceable context.
func NewLoggerContext(component string) context.Context {
	if loggerFactory == nil {
		panic("logger factory not set; call SetLoggerFactory at startup")
	}
	logger := loggerFactory.New(component)
	return context.WithValue(context.Background(), loggerCtxKey, logger)
}

// WithComponentPrefix creates a new context that inherits the trace ID from the
// parent context but uses a new component name for subsequent logs.
func WithComponentPrefix(ctx context.Context, name string) context.Context {
	if loggerFactory == nil {
		panic("logger factory not set; call SetLoggerFactory at startup")
	}
	currentLogger := Log(ctx)
	newLogger := loggerFactory.WithComponent(currentLogger, name)
	return context.WithValue(ctx, loggerCtxKey, newLogger)
}

// Log retrieves the Logger from the context.
// If no logger is found, it returns a "no-op" logger that does nothing.
// This ensures that logging calls are always safe.
func Log(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerCtxKey).(Logger); ok {
		return logger
	}

	log.Print("WARNING: nop logger called")
	return nop
}

// nopLogger is a no-op implementation of the Logger interface.
type nopLogger struct{}

func (n *nopLogger) Info(msg string)                {}
func (n *nopLogger) Infof(format string, v ...any)  {}
func (n *nopLogger) Debug(msg string)               {}
func (n *nopLogger) Debugf(format string, v ...any) {}
func (n *nopLogger) Warn(msg string)                {}
func (n *nopLogger) Warnf(format string, v ...any)  {}
func (n *nopLogger) Error(msg string)               {}
func (n *nopLogger) Errorf(format string, v ...any) {}
func (n *nopLogger) Fatal(msg string) {
	panic(msg)
}

func (n *nopLogger) Fatalf(format string, v ...any) {
	err := fmt.Errorf(format, v...)
	panic(err)
}

// nop is a single, shared instance of the no-op logger.
var nop = &nopLogger{}
