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

type WorkerPool struct {
	maxWorkers           int32
	currentWorkers       int32
	currentRunningWorker int32
	workerIdleTime       time.Duration
	status               Status
	workerQueue          chan *func()
}

// New 构造方法
func New(maxWorkers int) *WorkerPool {
	return newPool(maxWorkers, DefQueueSize, DefWorkerIdleTime)
}

// NewWithQueueSize 构造方法
func NewWithQueueSize(maxWorkers, queueSize int) *WorkerPool {
	return newPool(maxWorkers, queueSize, DefWorkerIdleTime)
}

// NewWithQueueSizeAndIdleTime 构造方法
func NewWithQueueSizeAndIdleTime(maxWorkers, queueSize int, idle time.Duration) *WorkerPool {
	return newPool(maxWorkers, queueSize, idle)
}

// newPool 正式构造
func newPool(maxWorkers, queueSize int, idle time.Duration) *WorkerPool {
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
	pool := &WorkerPool{
		maxWorkers:     int32(maxWorkers),
		currentWorkers: 0,
		workerIdleTime: idle,
		status:         StatusDispatchRunning,
		workerQueue:    make(chan *func(), queueSize),
	}
	//开启写成池的主循环
	go pool.dispatch()
	return pool
}

// Submit 提交任务
func (wp *WorkerPool) Submit(task func()) error {
	//如果任务为空则直接返回
	if task == nil {
		return nil
	}
	//当协程池的状态为关闭相关时，抛出异常
	if wp.status == StatusClosePretreatment || wp.status == StatusClose || wp.status == StatusCloseWait {
		return UseClosedPoolErr
	}
	//将任务放入任务队列中，准备执行
	select {
	case wp.workerQueue <- &task:
	default:
		return UnableToAddErr
	}
	return nil
}

// SubmitWait 等待任务完成的提交
func (wp *WorkerPool) SubmitWait(task func()) error {
	//如果任务为空则直接返回
	if task == nil {
		return nil
	}
	//如果协程池状态为关闭相关状态则抛出异常
	if wp.status == StatusClosePretreatment || wp.status == StatusClose || wp.status == StatusCloseWait {
		return UseClosedPoolErr
	}
	//构造一个提交完成管道，同时将需要执行的任务封装为doneFunc方法
	doneChan := make(chan bool)
	var doneFunc = func() {
		task()
		close(doneChan)
	}
	//阻塞等待，当被封装成doneFunc的任务执行完成并关闭管道时，会返回
	select {
	case wp.workerQueue <- &doneFunc:
		<-doneChan
	default:
		return UnableToAddErr
	}
	return nil
}

// dispatch 任务循环
func (wp *WorkerPool) dispatch() {
LOOP:
	//当协程池状态为任务循环状态，且当前工作协程数量小于最大运行数量是，则开始调度协程进行任务
	//在成功获取到管道中的任务以后，将协程池中当前运行协程的数量增加
	//当最大协程数目已经达到以后，会退出任务循环，即不启用新的协程运行task任务，保证协程数量不超出
	for atomic.LoadInt32((*int32)(&wp.status)) == int32(StatusDispatchRunning) &&
		atomic.LoadInt32(&wp.currentWorkers) < wp.maxWorkers {
		select {
		case tk, ok := <-wp.workerQueue:
			if !ok {
				break LOOP
			}
			//先将当前运行协程数目自增，保证启动协程后协程数目不会多余最大运行数目
			atomic.AddInt32(&wp.currentWorkers, 1)
			go wp.worker(tk)
		}
	}
}

func (wp *WorkerPool) worker(tk *func()) {
	(*tk)()
	if wp.doWorker() {
		atomic.AddInt32(&wp.currentWorkers, -1)
	}
}

// ReadyQueueLength 获取Queue的长度 此方法并不是线程安全的。
func (wp *WorkerPool) ReadyQueueLength() int {
	return len(wp.workerQueue)
}

func (wp *WorkerPool) doWorker() bool {
	if wp.workerIdleTime > 0 {
		//开启计时器，该计时器用于控制工作协程的过期时间，每次运行完task任务后计时器会被重置
		idle := time.NewTimer(wp.workerIdleTime)
		for atomic.LoadInt32((*int32)(&wp.status)) != int32(StatusClose) {
			select {
			case task, ok := <-wp.workerQueue:
				if !ok {
					break
				}
				(*task)()
				idle.Reset(wp.workerIdleTime)
			//当正在运行的协程长时间没有获取到任务超时后，会将该协程改变为任务分发协程
			//同时会将工作协程减一，同理当该协程在获取到任务，再次将工作协程数开启到最大数量时会退出
			case <-idle.C:
				if atomic.LoadInt32((*int32)(&wp.status)) <= wp.maxWorkers-1 &&
					atomic.CompareAndSwapInt32((*int32)(&wp.status),
						int32(StatusStable), int32(StatusDispatchRunning)) {
					atomic.AddInt32(&wp.currentWorkers, -1)
					wp.dispatch()
					return false
				}
				break
			}
		}
	} else {
		//如果没有设置超时时间，当任务完成后，则自动工作协程自动退出
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

func (wp *WorkerPool) Stop() {
	_ = wp.stop(StatusClose)
}

func (wp *WorkerPool) StopWait() {
	if !wp.stop(StatusCloseWait) {
		return
	}
	for atomic.LoadInt32(&wp.currentWorkers) > 0 {
		time.Sleep(time.Millisecond * 100)
	}
	atomic.StoreInt32((*int32)(&wp.status), int32(StatusClose))
}

func (wp *WorkerPool) Stopped() bool {
	return atomic.LoadInt32((*int32)(&wp.status)) == int32(StatusClose)
}

func (wp *WorkerPool) Status() Status {
	return Status(atomic.LoadInt32((*int32)(&wp.status)))
}

func (wp *WorkerPool) StopGetTasks() []func() {
	wp.stop(StatusClose)
	tasks := make([]func(), 0, wp.ReadyQueueLength())
	for task := range wp.workerQueue {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (wp *WorkerPool) stop(status Status) bool {
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
