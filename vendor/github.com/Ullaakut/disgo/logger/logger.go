// Package logger is a simplified logger which only
// handles two basic log levels.
package logger

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// Logger represents a disgo Logger.
// It writes the output on a given io.Writer
// and can toggle debug logs and have an error writer.
type Logger struct {
	// Writer on which the Info and Debug logs are written.
	standardOutput io.Writer
	// Writer on which the Error logs are written.
	errorOutput io.Writer

	// Whether or not Debug logs are enabled. If this is
	// disabled, Debug logs are not shown to the user.
	debug bool
}

// New creates a new Logger and binds the given writer to its outputs.
func New(w io.Writer, options ...func(*Logger)) (*Logger, error) {
	logger := Logger{
		standardOutput: w,
		errorOutput:    w,
	}

	for _, option := range options {
		option(&logger)
	}

	return &logger, nil
}

// WithDebug sets the logger debug mode to true.
func WithDebug(enabled bool) func(*Logger) {
	return func(logger *Logger) {
		logger.debug = enabled
	}
}

// WithErrorOutput binds a second writer on the logger, for logging errors.
func WithErrorOutput(w io.Writer) func(*Logger) {
	return func(l *Logger) {
		l.errorOutput = w
	}
}

// WithColors sets the use of colors in the logger. By default, whether or not
// colors are enabled depends on the user's TTY, but this option can be used
// to force colors to be enabled or disabled.
func WithColors(enabled bool) func(*Logger) {
	return func(_ *Logger) {
		color.NoColor = !enabled
	}
}

// Info writes an info log on the logger's standard writer.
func (l Logger) Info(a ...interface{}) {
	fmt.Fprint(l.standardOutput, a...)
}

// Infoln writes an info log on the logger's standard writer
// and appends a newline to its input.
func (l Logger) Infoln(a ...interface{}) {
	fmt.Fprintln(l.standardOutput, a...)
}

// Infof formats according to a format specifier and writes
// to the logger's standard writer.
func (l Logger) Infof(format string, a ...interface{}) {
	fmt.Fprintf(l.standardOutput, format, a...)
}

// Debug writes a debug log on the logger's standard writer if
// the debug logs are enabled.
func (l Logger) Debug(a ...interface{}) {
	if !l.debug {
		return
	}

	fmt.Fprint(l.standardOutput, a...)
}

// Debugln writes a debug log on the logger's standard writer if
//  the debug logs are enabled and appends a newline to its input.
func (l Logger) Debugln(a ...interface{}) {
	if !l.debug {
		return
	}

	fmt.Fprintln(l.standardOutput, a...)
}

// Debugf formats according to a format specifier and writes
// to the logger's standard writer if the debug logs are enabled.
func (l Logger) Debugf(format string, a ...interface{}) {
	if !l.debug {
		return
	}

	fmt.Fprintf(l.standardOutput, format, a...)
}

// Error writes an error log on the logger's error writer.
func (l Logger) Error(a ...interface{}) {
	fmt.Fprint(l.errorOutput, a...)
}

// Errorln writes an error log on the logger's error writer.
// It appends a newline to its input.
func (l Logger) Errorln(a ...interface{}) {
	fmt.Fprintln(l.errorOutput, a...)
}

// Errorf formats according to a format specifier and writes
// to the logger's error writer.
func (l Logger) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(l.errorOutput, format, a...)
}
