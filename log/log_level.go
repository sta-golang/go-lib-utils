package log

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var GlobalLevel = INFO

func GetLogLevel() int {
	return GlobalLevel
}
