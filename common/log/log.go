// Copyright 2017 The margin Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package log provides logging interfaces.
package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Log Level
const (
	FATAL = 0
	ERROR = 1
	WARN  = 2
	INFO  = 3
	DEBUG = 4
)

// Logger represents an active logging object that generates lines of output to an io.Writer.
//Each logging operation makes a single call to theWriter's Write method.
//A Logger can be used simultaneously from multiple goroutines
//it guarantees to serialize access to the Writer.
type Logger struct {
	log   *log.Logger
	file  *os.File
	level int
}

var (
	defaultLogLevel = DEBUG
	defaultLogger   = &Logger{log: log.New(os.Stdout, "", log.LstdFlags), file: nil, level: defaultLogLevel}
	errLevels       = []int{FATAL, ERROR, WARN, INFO, DEBUG}
	strLevels       = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
	selfHold        *Logger
)

// Setup creates a new Logger and hold it in this package.
// The out variable sets the destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The file argument defines the write log file path.
// if any error the os.Stdout will return
func Setup(file string, level int) (err error) {
	if selfHold != nil {
		return nil
	}
	selfHold, err = New(file, level)
	return
}

// New returns an Logger
func New(file string, levelIn int) (*Logger, error) {
	level := defaultLogLevel
	for _, v := range errLevels {
		if v == levelIn {
			level = v
		}
	}
	if file != "" {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return defaultLogger, err
		}
		logger := log.New(f, "", log.LstdFlags)
		return &Logger{log: logger, file: f, level: level}, nil
	}
	return &Logger{log: log.New(os.Stdout, "", log.LstdFlags), file: nil, level: level}, nil
}

// Close closes the open log file.
func Close() error {
	if selfHold.file != nil {
		return selfHold.file.Close()
	}
	return nil
}

// Error use the Error log level write data
func Error(format string, args ...interface{}) {
	if selfHold.level >= ERROR {
		selfHold.logCore(ERROR, format, args...)
	}
}

// Errorln like Error but ignore format
func Errorln(args ...interface{}) {
	if selfHold.level >= ERROR {
		selfHold.logCore(ERROR, "", args...)
	}
}

// Debug use the Error log level write data
func Debug(format string, args ...interface{}) {
	if selfHold.level >= DEBUG {
		selfHold.logCore(DEBUG, format, args...)
	}
}

// Debugln like Debug but ignore format
func Debugln(args ...interface{}) {
	if selfHold.level >= DEBUG {
		selfHold.logCore(DEBUG, "", args...)
	}
}

// Info use the Info log level write data
func Info(format string, args ...interface{}) {
	if selfHold.level >= INFO {
		selfHold.logCore(INFO, format, args...)
	}
}

// Infoln like info bug ignore format
func Infoln(args ...interface{}) {
	if selfHold.level >= INFO {
		selfHold.logCore(INFO, "", args...)
	}
}

// Warn use the Warn level write data
func Warn(format string, args ...interface{}) {
	if selfHold.level >= WARN {
		selfHold.logCore(WARN, format, args...)
	}
}

// Warnln like Warn but ignore format
func Warnln(args ...interface{}) {
	if selfHold.level >= WARN {
		selfHold.logCore(WARN, "", args...)
	}
}

// Fatal use the Fatal level write data
func Fatal(format string, args ...interface{}) {
	if selfHold.level >= FATAL {
		selfHold.logCore(FATAL, format, args...)
	}
}

// Fatalln like fatal bug ignore format
func Fatalln(args ...interface{}) {
	if selfHold.level >= FATAL {
		selfHold.logCore(FATAL, "", args...)
	}
}

// logCore handle the core log proc
func (l *Logger) logCore(level int, format string, args ...interface{}) {
	var (
		file string
		line int
		ok   bool
	)
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	if format == "" {
		l.log.Print(fmt.Sprintf("[%s] %s:%d %s", strLevels[level], file, line, fmt.Sprintln(args...)))
	} else {
		l.log.Print(fmt.Sprintf("[%s] %s:%d %s", strLevels[level], file, line, fmt.Sprintf(format, args...)))
	}
}
