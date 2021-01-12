package workerpool

import (
	"runtime"
	"sync"
)

var frameOnce sync.Once
var framePool Executor

func initPool() {
	frameOnce.Do(func() {
		framePool = New(runtime.NumCPU() + 1)
	})
}

func Submit(task func()) error {
	if framePool == nil {
		initPool()
	}
	return framePool.Submit(task)
}

func SubmitW(task func()) error {
	if framePool == nil {
		initPool()
	}
	return framePool.SubmitWait(task)
}

func Stop() error {
	if framePool == nil {
		initPool()
	}
	return framePool.Stop()
}

func StopWait() error {
	if framePool == nil {
		initPool()
	}
	return framePool.StopWait()
}

func StopGetTasks() []func() {
	if framePool == nil {
		initPool()
	}
	return framePool.StopGetTasks()
}
