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
	errLevels       = []int{Fatal, Error, Warn, Info, Debug}
	strLevels       = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
	selfHold        *Logger
)

/*
   Setup creates a new Logger and hold it in this package. The out variable sets the destination to
   which log data will be written. The prefix appears at the beginning of
   each generated log line. The file argument defines the write log file path.
   if any error the os.Stdout will return
*/

func Setup(file string, level int) (err error) {
	selfHold, err = New(file, level)
	return
}

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
			return DefaultLogger, err
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
	if selfHold.level >= Error {
		selfHold.logCore(Error, format, args...)
	}
}
func Errorln(args ...interface{}) {
	if selfHold.level >= Error {
		selfHold.logCore(Error, "", args...)
	}
}

// Debug use the Error log level write data
func Debug(format string, args ...interface{}) {
	if selfHold.level >= Debug {
		selfHold.logCore(Debug, format, args...)
	}
}
func Debugln(args ...interface{}) {
	if selfHold.level >= Debug {
		selfHold.logCore(Debug, "", args...)
	}
}

// Info use the Info log level write data
func Info(format string, args ...interface{}) {
	if selfHold.level >= Info {
		selfHold.logCore(Info, format, args...)
	}
}
func Infoln(args ...interface{}) {
	if selfHold.level >= Info {
		selfHold.logCore(Info, "", args...)
	}
}

// Warn use the Warn level write data
func Warn(format string, args ...interface{}) {
	if selfHold.level >= Warn {
		selfHold.logCore(Warn, format, args...)
	}
}
func Warnln(args ...interface{}) {
	if selfHold.level >= Warn {
		selfHold.logCore(Warn, "", args...)
	}
}

// Fatal use the Fatal level write data
func Fatal(format string, args ...interface{}) {
	if selfHold.level >= Fatal {
		selfHold.logCore(Fatal, format, args...)
	}
}
func Fatalln(args ...interface{}) {
	if selfHold.level >= Fatal {
		selfHold.logCore(Fatal, "", args...)
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
