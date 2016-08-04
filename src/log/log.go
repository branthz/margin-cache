package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	// Log Level
	Fatal = 0
	Error = 1
	Warn  = 2
	Info  = 3
	Debug = 4
)

/*
   A Logger represents an active logging object that generates lines of output
   to an io.Writer. Each logging operation makes a single call to the
   Writer's Write method. A Logger can be used simultaneously from multiple
   goroutines; it guarantees to serialize access to the Writer.
*/
type Logger struct {
	log   *log.Logger
	file  *os.File
	level int
}

var (
	defaultLogLevel = Debug
	DefaultLogger   = &Logger{log: log.New(os.Stdout, "", log.LstdFlags), file: nil, level: defaultLogLevel}
	errLevels       = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
)

/*
   New creates a new Logger. The out variable sets the destination to
   which log data will be written. The prefix appears at the beginning of
   each generated log line. The file argument defines the write log file path.
   if any error the os.Stdout will return
*/
func New(file string, levelStr string) (*Logger, error) {
	level := defaultLogLevel
	for lv, str := range errLevels {
		if str == levelStr {
			level = lv
		}
	}
	if file != "" {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return DefaultLogger, err
		}
		logger := log.New(f, "", log.LstdFlags)
		return &Logger{log: logger, file: f, level: level}, nil
	}
	return &Logger{log: log.New(os.Stdout, "", log.LstdFlags), file: nil, level: level}, nil
}

// Close closes the open log file.
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Error use the Error log level write data
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level >= Error {
		l.logCore(Error, format, args...)
	}
}
func (l *Logger) Errorln(args ...interface{}) {
	if l.level >= Error {
		l.logCore(Error,"", args...)
	}
}

// Debug use the Error log level write data
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level >= Debug {
		l.logCore(Debug, format, args...)
	}
}
func (l *Logger) Debugln(args ...interface{}) {
	if l.level >= Debug {
		l.logCore(Debug,"", args...)
	}
}

// Info use the Info log level write data
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level >= Info {
		l.logCore(Info, format, args...)
	}
}
func (l *Logger) Infoln(args ...interface{}) {
	if l.level >= Info {
		l.logCore(Info,"", args...)
	}
}

// Warn use the Warn level write data
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level >= Warn {
		l.logCore(Warn, format, args...)
	}
}
func (l *Logger) Warnln(args ...interface{}) {
	if l.level >= Warn {
		l.logCore(Warn,"", args...)
	}
}

// Fatal use the Fatal level write data
func (l *Logger) Fatal(format string, args ...interface{}) {
	if l.level >= Fatal {
		l.logCore(Fatal, format, args...)
	}
}
func (l *Logger) Fatalln(args ...interface{}) {
	if l.level >= Fatal {
		l.logCore(Fatal,"", args...)
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
    if format==""{
	    l.log.Print(fmt.Sprintf("[%s] %s:%d %s", errLevels[level], file, line, fmt.Sprintln(args...)))
    }else{
	    l.log.Print(fmt.Sprintf("[%s] %s:%d %s", errLevels[level], file, line, fmt.Sprintf(format, args...)))
    }
}
