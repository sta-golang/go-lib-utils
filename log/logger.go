package log

import "sync"

var (
	LEVEL_FLAGS = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	LEVEL_INDEX = map[string]Level{
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
)

// Logger log接口
type Logger interface {
	Debug(args ...interface{})

	Debugf(format string, args ...interface{})

	Info(args ...interface{})

	Infof(format string, args ...interface{})

	Warn(args ...interface{})

	Warnf(format string, args ...interface{})

	Error(args ...interface{})

	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})

	Fatalf(format string, args ...interface{})

	// SetLevel 设置输出端日志级别
	SetLevel(level Level)
	// GetLevel 获取输出端日志级别
	GetLevel() string

	setSkip(skip int)
}
