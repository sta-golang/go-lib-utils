package workerpool

import (
	"errors"
	"sync/atomic"
	"time"
)

type Status int32

const (
	StatusClose             Status = -1
	StatusCloseWait         Status = -2
	StatusClosePretreatment Status = -99
	StatusDispatchRunning   Status = 1
	StatusStable            Status = 2

	DefWorkerIdleTime = time.Second * 10

	DefQueueSize = 1024
	MaxWorkers   = 4096
	MaxQueueSize = 8192 << 1
)

var UnableToAddErr = errors.New("the queue is full and cannot be added")
var UseClosedPoolErr = errors.New("you are using a closed pool")

type workerPool struct {
	maxWorkers           int32
	currentWorkers       int32
	currentRunningWorker int32
	workerIdleTime       time.Duration
	status               Status
	workerQueue          chan *func()
}

func New(maxWorkers int) *workerPool {
	return newPool(maxWorkers, DefQueueSize, DefWorkerIdleTime)
}

func NewWithQueueSize(maxWorkers, queueSize int) *workerPool {
	return newPool(maxWorkers, queueSize, DefWorkerIdleTime)
}

func NewWithQueueSizeAndIdleTime(maxWorkers, queueSize int, idle time.Duration) *workerPool {
	return newPool(maxWorkers, queueSize, idle)
}

func newPool(maxWorkers, queueSize int, idle time.Duration) *workerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	if maxWorkers > MaxWorkers {
		maxWorkers = MaxWorkers
	}
	if queueSize < 1 {
		queueSize = 1
	}
	if queueSize > MaxQueueSize {
		queueSize = MaxQueueSize
	}
	pool := &workerPool{
		maxWorkers:     int32(maxWorkers),
		currentWorkers: 0,
		workerIdleTime: idle,
		status:         StatusDispatchRunning,
		workerQueue:    make(chan *func(), queueSize),
	}
	go pool.dispatch()
	return pool
}

func (wp *workerPool) Submit(task func()) error {
	if task == nil {
		return nil
	}
	if wp.status == StatusClosePretreatment || wp.status == StatusClose || wp.status == StatusCloseWait {
		return UseClosedPoolErr
	}
	select {
	case wp.workerQueue <- &task:
	default:
		return UnableToAddErr
	}
	return nil
}

func (wp *workerPool) SubmitWait(task func()) error {
	if task == nil {
		return nil
	}
	if wp.status == StatusClosePretreatment || wp.status == StatusClose || wp.status == StatusCloseWait {
		return UseClosedPoolErr
	}
	doneChan := make(chan bool)
	var doneFunc = func() {
		task()
		close(doneChan)
	}
	select {
	case wp.workerQueue <- &doneFunc:
		<-doneChan
	default:
		return UnableToAddErr
	}
	return nil
}

func (wp *workerPool) dispatch() {
LOOP:
	for atomic.LoadInt32((*int32)(&wp.status)) == int32(StatusDispatchRunning) &&
		atomic.LoadInt32(&wp.currentWorkers) < wp.maxWorkers {
		select {
		case tk, ok := <-wp.workerQueue:
			if !ok {
				break LOOP
			}
			atomic.AddInt32(&wp.currentWorkers, 1)
			go wp.worker(tk)
		}
	}
}

func (wp *workerPool) worker(tk *func()) {
	(*tk)()
	if wp.doWorker() {
		atomic.AddInt32(&wp.currentWorkers, -1)
	}
}

// ReadyQueueLength 获取Queue的长度 此方法并不是线程安全的。
func (wp *workerPool) ReadyQueueLength() int {
	return len(wp.workerQueue)
}

func (wp *workerPool) doWorker() bool {
	if wp.workerIdleTime > 0 {
		idle := time.Tick(wp.workerIdleTime)
		for atomic.LoadInt32((*int32)(&wp.status)) != int32(StatusClose) {
			select {
			case task, ok := <-wp.workerQueue:
				if !ok {
					break
				}
				(*task)()
			case <-idle:
				if atomic.LoadInt32((*int32)(&wp.status)) <= wp.maxWorkers-1 &&
					atomic.CompareAndSwapInt32((*int32)(&wp.status), int32(StatusStable), int32(StatusDispatchRunning)) {

					atomic.AddInt32(&wp.currentWorkers, -1)
					wp.dispatch()
					return false
				}
				break
			}
		}
	} else {
		for atomic.LoadInt32((*int32)(&wp.status)) != int32(StatusClose) {
			task, ok := <-wp.workerQueue
			if !ok {
				break
			}
			(*task)()
		}
	}
	return true
}

func (wp *workerPool) Stop() {
	_ = wp.stop(StatusClose)
}

func (wp *workerPool) StopWait() {
	if !wp.stop(StatusCloseWait) {
		return
	}
	for atomic.LoadInt32(&wp.currentWorkers) > 0 {
		time.Sleep(time.Millisecond * 100)
	}
}

func (wp *workerPool) StopGetTasks() []func() {
	wp.stop(StatusClose)
	tasks := make([]func(), 0, wp.ReadyQueueLength())
	for task := range wp.workerQueue {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (wp *workerPool) stop(status Status) bool {
	if atomic.LoadInt32((*int32)(&wp.status)) == int32(StatusClose) || atomic.LoadInt32((*int32)(&wp.status)) == int32(StatusClose) {
		return false
	}
	atomic.StoreInt32((*int32)(&wp.status), int32(StatusClosePretreatment))
	if atomic.CompareAndSwapInt32((*int32)(&wp.status), int32(StatusClosePretreatment), int32(status)) {
		close(wp.workerQueue)
		return true
	}
	return false
}
