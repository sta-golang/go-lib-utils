package group

import (
	"errors"
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"sync"
	"sync/atomic"
)

type TaskStatus int32
type AsyncGroupStatus int32

var pool *sync.Pool
var once sync.Once

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

func initPool() {
	once.Do(func() {
		pool = &sync.Pool{
			New: func() interface{} {
				return newTask()
			},
		}
	})

}

type asyncGroup struct {
	wg       sync.WaitGroup
	tasks    sync.Map
	Status   AsyncGroupStatus
	size     uint32
	executor workerpool.Executor
}

type task struct {
	name   string
	Status TaskStatus
	retVal interface{}
	retErr error
}

type AsyncGroupOption func(*asyncGroup)

func WithWorkPool(size int) func(*asyncGroup) {
	return func(ag *asyncGroup) {
		ag.executor = workerpool.New(size)
	}
}

func WithWorkPoolAndQueueSize(wSize, qSize int) func(*asyncGroup) {
	return func(ag *asyncGroup) {
		ag.executor = workerpool.NewWithQueueSize(wSize, qSize)
	}
}

func newTask() *task {
	return &task{
		Status: TaskStatusInit,
		retVal: nil,
		retErr: nil,
	}
}

func (t *task) Ret() (interface{}, error) {
	if t.Status != TaskStatusFinish {
		return nil, nil
	}
	return t.retVal, t.retErr
}

func NewAsyncGroup(opts ...AsyncGroupOption) *asyncGroup {
	ag := defaultAsyncGroup()
	for _, fn := range opts {
		fn(ag)
	}
	return ag
}

func defaultAsyncGroup() *asyncGroup {
	return &asyncGroup{
		wg:     sync.WaitGroup{},
		tasks:  sync.Map{},
		Status: AsyncGroupStatusInit,
	}
}

func (ag *asyncGroup) Add(name string, fn func() (interface{}, error)) error {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		return UseClosedErr
	}
	if _, ok := ag.tasks.Load(name); ok {
		return NameRepeatErr
	}
	curTk := pool.Get().(*task)
	if _, ok := ag.tasks.LoadOrStore(name, curTk); ok {
		pool.Put(curTk)
		return NameRepeatErr
	}
	atomic.AddUint32(&ag.size, 1)
	curTk.Status = TaskStatusInit
	ag.wg.Add(1)
	taskFn := func() {
		defer func() {
			if pErr := recover(); pErr != nil {
				curTk.retErr = fmt.Errorf("panic: %v", pErr)
			}
			curTk.Status = TaskStatusFinish
			ag.wg.Done()
		}()
		curTk.Status = TaskStatusRunning
		curTk.retVal, curTk.retErr = fn()
	}
	if ag.executor == nil {
		go taskFn()
		return nil
	}
	err := ag.executor.Submit(taskFn)
	return err
}

func (ag *asyncGroup) Close() {
	if !atomic.CompareAndSwapInt32((*int32)(&ag.Status), AsyncGroupStatusInit, AsyncGroupStatusClose) {
		return
	}
	if ag.size > maxAsyncCloseSize {
		err := workerpool.Submit(ag.doClose)
		if err != nil {
			go ag.doClose()
		}
		return
	} else {
		ag.doClose()
	}
}

func (ag *asyncGroup) Wait() {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		log.Warn("use close asyncGroup For wait")
	}
	ag.wg.Wait()
}

func (ag *asyncGroup) Iterator() []*task {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		return nil
	}
	ret := make([]*task, 0, atomic.LoadUint32(&ag.size))
	ag.tasks.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(*task))
		return true
	})
	return ret
}

func (ag *asyncGroup) GetTask(name string) *task {
	if atomic.LoadInt32((*int32)(&ag.Status)) == AsyncGroupStatusClose {
		return nil
	}
	ret, _ := ag.tasks.Load(name)
	return ret.(*task)
}

func (ag *asyncGroup) DestoryTask(tk *task) {
	if pool == nil {
		return
	}
	pool.Put(tk)
}

func (ag *asyncGroup) doClose() {
	ag.tasks.Range(func(key, value interface{}) bool {
		pool.Put(value)
		return true
	})
}
