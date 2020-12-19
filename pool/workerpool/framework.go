package workerpool

import (
	"runtime"
	"sync"
)

var frameOnce sync.Once
var framePool *workerPool

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

func Stop() {
	if framePool == nil {
		initPool()
	}
	framePool.Stop()
}

func StopWait() {
	if framePool == nil {
		initPool()
	}
	framePool.StopWait()
}

func StopGetTasks() {
	if framePool == nil {
		initPool()
	}
	framePool.StopGetTasks()
}
