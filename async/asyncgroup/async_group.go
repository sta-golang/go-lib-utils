// Package asyncgroup 异步任务执行组，封装了异步任务，使得异步任务调用变简单
package asyncgroup

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
)

// AsyncGroupStatus 异步组的状态
type AsyncGroupStatus int32

var pool *sync.Pool

const (
	TaskStatusInit = iota
	TaskStatusRunning
	TaskStatusFinish

	AsyncGroupStatusInit = iota
	AsyncGroupStatusClose

	maxAsyncCloseSize = 100
)

var NameRepeatErr = errors.New("you use the wrong name")
var UseClosedErr = errors.New("you used the closed component")

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return newTask()
		},
	}
}

// Group 执行组
type Group struct {
	wg       sync.WaitGroup
	lockFlag bool
	mu       sync.Mutex
	tasks    map[string]*task
	Status   AsyncGroupStatus
	size     uint32
	executor workerpool.Executor
}

// 创建函数
func New(opts ...Option) *Group {
	ag := defaultAsyncGroup()
	for _, fn := range opts {
		fn(ag)
	}
	return ag
}

func defaultAsyncGroup() *Group {
	return &Group{
		wg:       sync.WaitGroup{},
		lockFlag: false,
		Status:   AsyncGroupStatusInit,
	}
}

// Add 添加一个任务到异步执行组
func (ag *Group) Add(name string, fn func() (interface{}, error)) error {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		return UseClosedErr
	}
	curTask := pool.Get().(*task)
	if ag.lockFlag {
		ag.mu.Lock()
	}
	if ag.tasks == nil {
		ag.tasks = make(map[string]*task)
	}
	if _, ok := ag.tasks[name]; ok {
		pool.Put(curTask)
		if ag.lockFlag {
			ag.mu.Unlock()
		}
		return NameRepeatErr
	}
	ag.tasks[name] = curTask
	if ag.lockFlag {
		ag.mu.Unlock()
	}
	atomic.AddUint32(&ag.size, 1)
	curTask.Status = TaskStatusInit
	ag.wg.Add(1)
	taskFn := func() {
		defer func() {
			if pErr := recover(); pErr != nil {
				curTask.retErr = fmt.Errorf("panic: %v", pErr)
			}
			curTask.Status = TaskStatusFinish
			ag.wg.Done()
		}()
		curTask.Status = TaskStatusRunning
		curTask.retVal, curTask.retErr = fn()
	}
	if ag.executor == nil {
		go taskFn()
		return nil
	}
	if err := ag.executor.Submit(taskFn); err != nil {
		go taskFn()
		log.Warnf("workerpool err : %v", err)
	}
	return nil
}

// Close 关闭异步执行组
func (ag *Group) Close() {
	ag.Wait()
	if !atomic.CompareAndSwapInt32((*int32)(&ag.Status), AsyncGroupStatusInit, AsyncGroupStatusClose) {
		return
	}
	if ag.size > maxAsyncCloseSize {
		go ag.doClose()
	} else {
		ag.doClose()
	}
}

// Shutdown 立刻关闭 不进行等待
func (ag *Group) Shutdown() {
	if !atomic.CompareAndSwapInt32((*int32)(&ag.Status), AsyncGroupStatusInit, AsyncGroupStatusClose) {
		return
	}
	go func() {
		ag.doWait(true)
		ag.doClose()
	}()
}

// Wait 等待任务执行完毕
func (ag *Group) Wait() {
	ag.doWait(false)
}

func (ag *Group) doWait(shutdownFlag bool) {
	if !shutdownFlag && atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		log.Warn("use close asyncGroup For wait")
	}
	ag.wg.Wait()
}

// Size 获取任务数量大小
func (ag *Group) Size() int {
	return int(atomic.LoadUint32(&ag.size))
}

// Iterator 获取任务列表
func (ag *Group) Iterator() []*task {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		return nil
	}
	if ag.lockFlag {
		ag.mu.Lock()
		defer ag.mu.Unlock()
	}
	ret := make([]*task, 0, atomic.LoadUint32(&ag.size))
	for _, task := range ag.tasks {
		ret = append(ret, task)
	}
	return ret
}

// GetTask 获取指定Task
func (ag *Group) GetTask(name string) *task {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		return nil
	}
	if ag.lockFlag {
		ag.mu.Lock()
		defer ag.mu.Unlock()
	}
	return ag.tasks[name]
}

func (ag *Group) doClose() {
	// 这里已经关闭 不需要再加锁了
	for _, task := range ag.tasks {
		pool.Put(task)
	}
	ag.mu.Lock()
	defer ag.mu.Unlock()
	ag.tasks = nil
}
