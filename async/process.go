package async

import (
	"sync/atomic"
)

type ProcessSource int32

const (
	tryProcess  = 1
	initProcess = 0
)

func NewProcessSource() ProcessSource {
	return initProcess
}

func TryProcess(flag *ProcessSource) bool {
	return atomic.CompareAndSwapInt32((*int32)(flag), initProcess, tryProcess)
}

func EndProcess(flag *ProcessSource) {
	atomic.CompareAndSwapInt32((*int32)(flag), tryProcess, initProcess)
}
