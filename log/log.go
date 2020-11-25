package log

var gloabLogger Logger
var FrameworkLogger = NewConsoleLog(INFO, "[STA : utils-log]")

const (
	CurrencyError    = 00
	FileCreateError  = 01
	FileWriterError  = 03
	FileCloseError   = 04
	FileReadDirError = 05
)

func init() {
	gloabLogger = NewConsoleLog(INFO, PREFIX)
}

func SetPrefix(str string) {
	PREFIX = str
}

func SetLevel(level Level) {
	gloabLogger.SetLevel(level)
}

// GetLevel 获取输出端日志级别
func GetLevel() string {
	return gloabLogger.GetLevel()
}

func Debugf(format string, args ...interface{}) {
	gloabLogger.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	gloabLogger.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	gloabLogger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	gloabLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	gloabLogger.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	gloabLogger.Debug(args...)
}

func Warn(args ...interface{}) {
	gloabLogger.Warn(args...)
}

func Info(args ...interface{}) {
	gloabLogger.Info(args...)
}

func Error(args ...interface{}) {
	gloabLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	gloabLogger.Fatal(args...)
}
