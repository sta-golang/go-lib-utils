package log

import (
	"sync"

	"context"
)

var (
	levelFlages = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	levelIndexs = map[string]Level{
		"DEBUG": DEBUG,
		"INFO":  INFO,
		"WARN":  WARNING,
		"ERROR": ERROR,
		"FATAL": FATAL,
	}
	logGlobalMutex sync.Mutex
	PREFIX         = "[FOUR-SEASONS: STA]"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL

	LogContextKey = "staLoggerCtx"
)

func GetLevelName(level Level) string {
	if level < DEBUG || level > FATAL {
		return ""
	}
	return levelFlages[level]
}

func LogContextKeyMap(ctx context.Context, keyMap map[string]string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, LogContextKey, keyMap)
}

// Logger log接口
type Logger interface {
	Debug(args ...interface{})

	DebugContext(ctx context.Context, args ...interface{})

	Debugf(format string, args ...interface{})

	DebugContextf(ctx context.Context, format string, args ...interface{})

	Info(args ...interface{})

	InfoContext(ctx context.Context, args ...interface{})

	Infof(format string, args ...interface{})

	InfoContextf(ctx context.Context, format string, args ...interface{})

	Warn(args ...interface{})

	WarnContext(ctx context.Context, args ...interface{})

	Warnf(format string, args ...interface{})

	WarnContextf(ctx context.Context, format string, args ...interface{})

	Error(args ...interface{})

	ErrorContext(ctx context.Context, args ...interface{})

	Errorf(format string, args ...interface{})

	ErrorContextf(ctx context.Context, format string, args ...interface{})

	Fatal(args ...interface{})

	FatalContext(ctx context.Context, args ...interface{})

	Fatalf(format string, args ...interface{})

	FatalContextf(ctx context.Context, format string, args ...interface{})

	// SetLevel 设置输出端日志级别
	SetLevel(level Level)
	// GetLevel 获取输出端日志级别
	GetLevel() string

	setSkip(skip int)
}
