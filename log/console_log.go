package log

import (
	"fmt"
	tm "github.com/sta-golang/go-lib-utils/time"
	"io"
	"os"
	"runtime"
)

var (
	console io.Writer = os.Stdout
)

type consoleLog struct {
	nonClear bool
	prefix   string
	level    Level
	skip     int
}

func NewConsoleLog(level Level, prefix string) Logger {
	if prefix == "windows" {
		prefix = PREFIX
	}
	return &consoleLog{
		nonClear: runtime.GOOS == "windows",
		level:    level,
		prefix:   prefix,
		skip:     dfsStep,
	}
}

func (cl *consoleLog) setSkip(skip int) {
	cl.skip = skip
}

func (cl *consoleLog) print(level Level, format string, args ...interface{}) {
	if level < cl.level {
		return
	}
	var logFormat = "%s %s [%s] %s ==> %s\n"
	if !cl.nonClear {
		switch level {
		case DEBUG:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[34m%s\033[0m] %s ==> %s\n"
		case INFO:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[32m%s\033[0m] %s ==> %s\n"
		case WARNING:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[33m%s\033[0m] %s ==> %s\n"
		case ERROR:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[31m%s\033[0m] %s ==> %s\n"
		case FATAL:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[35m%s\033[0m] %s ==> %s\n"
		}
	}
	_, transFile, transLine, _ := runtime.Caller(cl.skip)
	fmt.Fprintf(console, logFormat, cl.prefix, tm.GetNowDateTimeStr(),
		levelFlages[level], fmt.Sprintf("%s:%d", transFile, transLine), fmt.Sprintf(format, args...))

}

func (cl *consoleLog) println(level Level, args ...interface{}) {
	if level < cl.level {
		return
	}
	var logFormat = "%s %s [%s] %s ==> "
	if !cl.nonClear {
		switch level {
		case DEBUG:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[34m%s\033[0m] %s ==> "
		case INFO:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[32m%s\033[0m] %s ==> "
		case WARNING:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[33m%s\033[0m] %s ==> "
		case ERROR:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[31m%s\033[0m] %s ==> "
		case FATAL:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[35m%s\033[0m] %s ==> "
		}
	}
	_, transFile, transLine, _ := runtime.Caller(cl.skip)
	fmt.Fprintf(console, fmt.Sprintf("%s%s\n", fmt.Sprintf(logFormat, cl.prefix, tm.GetNowDateTimeStr(),
		levelFlages[level], fmt.Sprintf("%s:%d", transFile, transLine)), fmt.Sprint(args...)))

}

func (cl *consoleLog) SetLevel(level Level) {
	if level < DEBUG || level > FATAL {
		return
	}
	cl.level = level
}

// GetLevel 获取输出端日志级别
func (cl *consoleLog) GetLevel() string {
	return levelFlages[cl.level]
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
