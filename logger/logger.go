// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// @budui copy from "miniflux.app/logger", may be modified later.

package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"wrong.wang/x/go-isso/version"
)

var requestedLevel = InfoLevel
var displayDateTime = false
var displayRuntime = false

// LogLevel type.
type LogLevel uint32

var outputWriter io.Writer = os.Stderr

const (
	// FatalLevel should be used in fatal situations, the app will exit.
	FatalLevel LogLevel = iota

	// ErrorLevel should be used when someone should really look at the error.
	ErrorLevel

	// InfoLevel should be used during normal operations.
	InfoLevel

	// DebugLevel should be used only during development.
	DebugLevel
)

func (level LogLevel) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// EnableDateTime enables date time in log messages.
func EnableDateTime() {
	displayDateTime = true
}

// SetRuntime enables func name, line in log messages.
func SetRuntime(ok bool) {
	displayRuntime = ok
}

// EnableDebug increases logging, more verbose (debug)
func EnableDebug() {
	requestedLevel = DebugLevel
	SetRuntime(true)
	formatMessage(InfoLevel, "Debug mode enabled")
}

// Debug sends a debug log message.
func Debug(format string, v ...interface{}) {
	if requestedLevel >= DebugLevel {
		formatMessage(DebugLevel, format, v...)
	}
}

// Info sends an info log message.
func Info(format string, v ...interface{}) {
	if requestedLevel >= InfoLevel {
		formatMessage(InfoLevel, format, v...)
	}
}

// Error sends an error log message.
func Error(format string, v ...interface{}) {
	if requestedLevel >= ErrorLevel {
		formatMessage(ErrorLevel, format, v...)
	}
}

// Fatal sends a fatal log message and stop the execution of the program.
func Fatal(format string, v ...interface{}) {
	if requestedLevel >= FatalLevel {
		formatMessage(FatalLevel, format, v...)
		os.Exit(1)
	}
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	outputWriter = w
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	return strings.TrimPrefix(name, version.Mod)
}

func formatMessage(level LogLevel, format string, v ...interface{}) {
	var prefix string
	var caller string

	if displayDateTime {
		prefix = fmt.Sprintf("[%s]", time.Now().Format("2006-01-02T15:04:05"))
	}

	prefix += fmt.Sprintf(" [%s]", level)

	if displayRuntime {
		pc, _, _, ok := runtime.Caller(2)
		if !ok {
			caller = "unkown"
		} else {
			fn := runtime.FuncForPC(pc)
			caller = fn.Name()
		}
		prefix += fmt.Sprintf(" %s - ", funcname(caller))
	}

	fmt.Fprintf(outputWriter, prefix+format+"\n", v...)
}
