package log

var globalLogger Logger
var FrameworkLogger = NewConsoleLog(DEBUG, "[STA : go-groups-utils]")

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

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}
