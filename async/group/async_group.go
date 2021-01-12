package group

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"github.com/sta-golang/go-lib-utils/str"
	"sync"
)

type TaskStatus int32

var pool *sync.Pool
var once sync.Once

const (
	TaskStatusInit = iota
	TaskStatusRunning
	TaskStatusFinish

	maxAsyncCloseSize = 100
)

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
	wg    sync.WaitGroup
	tasks map[string]*task
}

type task struct {
	Status TaskStatus
	retVal interface{}
	retErr error
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

func NewAsyncGroup(size int) *asyncGroup {
	if pool == nil {
		initPool()
	}
	return &asyncGroup{
		wg:    sync.WaitGroup{},
		tasks: make(map[string]*task, size<<1),
	}
}

func (ag *asyncGroup) Add(fn func() (interface{}, error)) string {
	id := str.XID()
	for _, ok := ag.tasks[id]; ok; {
		id = str.XID()
	}
	curTk := pool.Get().(*task)
	ag.tasks[id] = curTk
	curTk.Status = TaskStatusInit
	ag.wg.Add(1)
	go func(tk *task) {
		defer func() {
			if pErr := recover(); pErr != nil {
				tk.retErr = fmt.Errorf("panic: %v", pErr)
			}
			tk.Status = TaskStatusFinish
			ag.wg.Done()
		}()
		tk.Status = TaskStatusRunning
		tk.retVal, tk.retErr = fn()
	}(curTk)

	return id
}

func (ag *asyncGroup) Close() {
	if len(ag.tasks) > maxAsyncCloseSize {
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
	ag.wg.Wait()
}

func (ag *asyncGroup) Iterator() []*task {
	ret := make([]*task, 0, len(ag.tasks))
	for _, val := range ag.tasks {
		ret = append(ret, val)
	}
	return ret
}

func (ag *asyncGroup) GetTask(requestID string) *task {
	return ag.tasks[requestID]
}

func (ag *asyncGroup) doClose() {
	for _, val := range ag.tasks {
		pool.Put(val)
	}
	ag.tasks = nil
}
