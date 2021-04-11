package log

import "context"

var globalLogger Logger
var ConsoleLogger = NewConsoleLog(DEBUG, "[STA : go-groups-utils]")

const (
	CurrencyError    = 00
	FileCreateError  = 01
	FileWriterError  = 03
	FileCloseError   = 04
	FileReadDirError = 05

	dfsStep    = 2
	globalSkip = 3
)

func init() {
	globalLogger = NewConsoleLog(INFO, PREFIX)
	globalLogger.setSkip(globalSkip)
}

func SetGlobalLogger(logger Logger) {
	globalLogger = logger
	globalLogger.setSkip(globalSkip)
}

func SetPrefix(str string) {
	PREFIX = str
}

func SetLevel(level Level) {
	globalLogger.SetLevel(level)
}

// GetLevel 获取输出端日志级别
func GetLevel() string {
	return globalLogger.GetLevel()
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func DebugContextf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.DebugContextf(ctx, format, args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func WarnContextf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.WarnContextf(ctx, format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func InfoContextf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.InfoContextf(ctx, format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.ErrorContextf(ctx, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func FatalContextf(ctx context.Context, format string, args ...interface{}) {
	globalLogger.FatalContextf(ctx, format, args...)
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func DebugContext(ctx context.Context, args ...interface{}) {
	globalLogger.DebugContext(ctx, args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func WarnContext(ctx context.Context, args ...interface{}) {
	globalLogger.WarnContext(ctx, args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func InfoContext(ctx context.Context, args ...interface{}) {
	globalLogger.InfoContext(ctx, args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func ErrorContext(ctx context.Context, args ...interface{}) {
	globalLogger.ErrorContext(ctx, args...)
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func FatalContext(ctx context.Context, args ...interface{}) {
	globalLogger.FatalContext(ctx, args...)
}
