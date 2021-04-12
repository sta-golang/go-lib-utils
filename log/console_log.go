package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/str"
	tm "github.com/sta-golang/go-lib-utils/time"
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
	if prefix == "" {
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

func (cl *consoleLog) printContext(ctx context.Context, level Level, format string, args ...interface{}) {
	if level < cl.level {
		return
	}
	var logFormat = "%s %s [%s] %s ==> %s ctx: %s\n"
	if !cl.nonClear {
		switch level {
		case DEBUG:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[34m%s\033[0m] %s ==> %s ctx: %s\n"
		case INFO:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[32m%s\033[0m] %s ==> %s ctx: %s\n"
		case WARNING:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[33m%s\033[0m] %s ==> %s ctx: %s\n"
		case ERROR:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[31m%s\033[0m] %s ==> %s ctx: %s\n"
		case FATAL:
			logFormat = "\033[1;38m%s\033[0m \033[36m%s\033[0m [\033[35m%s\033[0m] %s ==> %s ctx: %s\n"
		}
	}
	_, transFile, transLine, _ := runtime.Caller(cl.skip)
	ctxInfo := ""
	if ctx != nil {
		keyMap := ctx.Value(LogContextKey)
		if keyMap != nil {
			bys, _ := codec.API.JsonAPI.Marshal(keyMap)
			ctxInfo = str.BytesToString(bys)
		}
	}
	fmt.Fprintf(console, logFormat, cl.prefix, tm.GetNowDateTimeStr(),
		levelFlages[level], fmt.Sprintf("%s:%d", transFile, transLine), fmt.Sprintf(format, args...), ctxInfo)

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

func (cl *consoleLog) printlnContext(ctx context.Context, level Level, args ...interface{}) {
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
	ctxInfo := ""
	if ctx != nil {
		keyMap := ctx.Value(LogContextKey)
		if keyMap != nil {
			bys, _ := codec.API.JsonAPI.Marshal(keyMap)
			ctxInfo = str.BytesToString(bys)
		}
	}
	fmt.Fprintf(console, fmt.Sprintf("%s%s%s\n", fmt.Sprintf(logFormat, cl.prefix, tm.GetNowDateTimeStr(),
		levelFlages[level], fmt.Sprintf("%s:%d", transFile, transLine)), fmt.Sprint(args...), fmt.Sprint(" ctx: ", ctxInfo)))

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

func (cl *consoleLog) DebugContextf(ctx context.Context, format string, args ...interface{}) {
	cl.printContext(ctx, DEBUG, format, args...)
}

func (cl *consoleLog) Warnf(format string, args ...interface{}) {
	cl.print(WARNING, format, args...)
}

func (cl *consoleLog) WarnContextf(ctx context.Context, format string, args ...interface{}) {
	cl.printContext(ctx, WARNING, format, args...)
}

func (cl *consoleLog) Infof(format string, args ...interface{}) {
	cl.print(INFO, format, args...)
}

func (cl *consoleLog) InfoContextf(ctx context.Context, format string, args ...interface{}) {
	cl.printContext(ctx, INFO, format, args...)
}

func (cl *consoleLog) Errorf(format string, args ...interface{}) {
	cl.print(ERROR, format, args...)
}

func (cl *consoleLog) ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	cl.printContext(ctx, ERROR, format, args...)
}

func (cl *consoleLog) Fatalf(format string, args ...interface{}) {
	cl.print(FATAL, format, args...)
}

func (cl *consoleLog) FatalContextf(ctx context.Context, format string, args ...interface{}) {
	cl.printContext(ctx, FATAL, format, args...)
}

func (cl *consoleLog) Debug(args ...interface{}) {
	cl.println(DEBUG, args...)
}

func (cl *consoleLog) DebugContext(ctx context.Context, args ...interface{}) {
	cl.printlnContext(ctx, DEBUG, args...)
}

func (cl *consoleLog) Warn(args ...interface{}) {
	cl.println(WARNING, args...)
}

func (cl *consoleLog) WarnContext(ctx context.Context, args ...interface{}) {
	cl.printlnContext(ctx, WARNING, args...)
}

func (cl *consoleLog) Info(args ...interface{}) {
	cl.println(INFO, args...)
}

func (cl *consoleLog) InfoContext(ctx context.Context, args ...interface{}) {
	cl.printlnContext(ctx, INFO, args...)
}

func (cl *consoleLog) Error(args ...interface{}) {
	cl.println(ERROR, args...)
}

func (cl *consoleLog) ErrorContext(ctx context.Context, args ...interface{}) {
	cl.printlnContext(ctx, ERROR, args...)
}

func (cl *consoleLog) Fatal(args ...interface{}) {
	cl.println(FATAL, args...)
}

func (cl *consoleLog) FatalContext(ctx context.Context, args ...interface{}) {
	cl.printlnContext(ctx, FATAL, args...)
}
