package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	console io.Writer = os.Stdout
)

type consoleLog struct {
	nonClear bool
	prefix   string
	level    Level
}

func NewConsoleLog(level Level, prefix string) Logger {
	if prefix == "" {
		prefix = PREFIX
	}
	return &consoleLog{
		nonClear: runtime.GOOS == "windows",
		level:    level,
		prefix:   prefix,
	}
}

func (cl *consoleLog) print(level Level, format string, args ...interface{}) {

	var logFormat = "%s %s [%s] ==> %s\n"
	if !cl.nonClear {
		switch level {
		case DEBUG:
			logFormat = "%s \033[36m%s\033[0m [\033[34m%s\033[0m] %s\n"
		case INFO:
			logFormat = "%s \033[36m%s\033[0m [\033[32m%s\033[0m] %s\n"
		case WARNING:
			logFormat = "%s \033[36m%s\033[0m [\033[33m%s\033[0m] %s\n"
		case ERROR:
			logFormat = "%s \033[36m%s\033[0m [\033[31m%s\033[0m] %s\n"
		case FATAL:
			logFormat = "%s \033[36m%s\033[0m [\033[35m%s\033[0m] %s\n"
		}
	}

	fmt.Fprintf(console, logFormat, cl.prefix, time.Now().Format(TIME_FORMAT),
		LEVEL_FLAGS[level], fmt.Sprintf(format, args...))

}

func (cl *consoleLog) println(level Level, args ...interface{}) {
	if level < cl.level {
		return
	}
	var logFormat = "%s %s [%s] ==> "
	if !cl.nonClear {
		switch level {
		case DEBUG:
			logFormat = "%s \033[36m%s\033[0m [\033[34m%s\033[0m] %s\n"
		case INFO:
			logFormat = "%s \033[36m%s\033[0m [\033[32m%s\033[0m] %s\n"
		case WARNING:
			logFormat = "%s \033[36m%s\033[0m [\033[33m%s\033[0m] %s\n"
		case ERROR:
			logFormat = "%s \033[36m%s\033[0m [\033[31m%s\033[0m] %s\n"
		case FATAL:
			logFormat = "%s \033[36m%s\033[0m [\033[35m%s\033[0m] %s\n"
		}
	}

	fmt.Fprintf(console, logFormat, cl.prefix, time.Now().Format(TIME_FORMAT),
		LEVEL_FLAGS[level], args)
	fmt.Fprintln(console)
}

func (cl *consoleLog) SetLevel(level Level) {
	if level < DEBUG || level > FATAL {
		return
	}
	cl.level = level
}

// GetLevel 获取输出端日志级别
func (cl *consoleLog) GetLevel() string {
	return LEVEL_FLAGS[cl.level]
}

func (cl *consoleLog) Debugf(format string, args ...interface{}) {
	cl.print(DEBUG, format, args...)
}

func (cl *consoleLog) Warnf(format string, args ...interface{}) {
	cl.print(WARNING, format, args...)
}

func (cl *consoleLog) Infof(format string, args ...interface{}) {
	cl.print(INFO, format, args...)
}

func (cl *consoleLog) Errorf(format string, args ...interface{}) {
	cl.print(ERROR, format, args...)
}

func (cl *consoleLog) Fatalf(format string, args ...interface{}) {
	cl.print(FATAL, format, args...)
}

func (cl *consoleLog) Debug(args ...interface{}) {
	cl.println(DEBUG, args...)
}

func (cl *consoleLog) Warn(args ...interface{}) {
	cl.println(WARNING, args...)
}

func (cl *consoleLog) Info(args ...interface{}) {
	cl.println(INFO, args...)
}

func (cl *consoleLog) Error(args ...interface{}) {
	cl.println(ERROR, args...)
}

func (cl *consoleLog) Fatal(args ...interface{}) {
	cl.println(FATAL, args...)
}
