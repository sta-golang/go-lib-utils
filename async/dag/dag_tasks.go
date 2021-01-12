package dag

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// TaskState 任务的执行状态
type TaskState int32
type DagState int32

// const 任务执行状态枚举值，0为执行，1执行中，2执行完成
const (
	TaskInit    TaskState = 0
	TaskReady   TaskState = 1
	TaskRunning TaskState = 2
	TaskFinish  TaskState = 3

	DagInit    DagState = 0
	DagReady   DagState = 1
	DagRunning DagState = 2
	DagFinish  DagState = 3

	maxRetry = 3
)

var (
	endTask *task
)

func init() {
	endTask = &task{}
}

type DagTasks struct {
	state      DagState
	wg         sync.WaitGroup
	workerPool *workerpool.WorkerPool
	root       *task
}

type task struct {
	name            string
	fn              func(ctx context.Context, helper TaskHelper) (interface{}, error)
	retVal          interface{}
	retErr          error
	state           TaskState
	finishCnt       int32
	parents         map[*task]bool
	childrenTasks   map[string]*task
	childrenKeyList []string
}

type TaskHelper struct {
	task *task
}

func NewTask(name string, fn func(ctx context.Context, helper TaskHelper) (interface{}, error)) *task {
	return &task{
		name:          name,
		fn:            fn,
		retVal:        nil,
		retErr:        nil,
		state:         TaskInit,
		parents:       make(map[*task]bool),
		childrenTasks: make(map[string]*task),
	}
}

func NewDag(root *task) DagTasks {
	return DagTasks{
		wg:    sync.WaitGroup{},
		root:  root,
		state: DagInit,
	}
}

func (tk *task) AddSubTask(subT *task) {
	if tk.Equals(subT) {
		return
	}
	if _, ok := subT.parents[tk]; ok {
		return
	}
	tk.childrenTasks[subT.name] = subT
	subT.parents[tk] = true
	tk.childrenKeyList = append(tk.childrenKeyList, subT.name)
}

func (tk *task) Equals(other *task) bool {
	if tk.name == other.name {
		return true
	}
	if tk == other {
		return true
	}
	return false
}

func (tk *task) checkDependenceForDfs() bool {
	m := make(map[*task]bool)
	return doCheckDependenceForDfs(&m, tk)
}

func doCheckDependenceForDfs(set *map[*task]bool, tk *task) bool {
	if _, ok := (*set)[tk]; ok {
		return true
	}
	(*set)[tk] = true
	for _, temp := range tk.childrenTasks {
		if doCheckDependenceForDfs(set, temp) {
			return true
		}
	}
	return false
}

func (dt *DagTasks) checkDependenceForDfs() bool {
	return dt.root.checkDependenceForDfs()
}

func (tk *task) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("name : %s children : \n", tk.name))
	for _, ch := range tk.childrenTasks {
		buffer.WriteString(ch.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (dt *DagTasks) Do(ctx context.Context, checkDependence bool) (hasDependence bool) {
	// 检查是否存在互相依赖关系
	if checkDependence {
		if dt.checkDependenceForDfs() {
			return true
		}
	}
	if len(dt.root.parents) != 0 {
		return true
	}
	if dt.root == nil {
		return
	}
	runCh := make(chan *task, runtime.NumCPU())
	defer func() {
		if e := recover(); e != nil {
			log.Errorf("panic:%v", string(debug.Stack()))
		}
	}()
	dt.wg.Add(1)
	go dt.startRun(runCh)
	go dt.doExec(ctx, runCh)
	dt.wg.Wait()
	close(runCh)
	return
}

// GetRootTask 获取根节点任务
func (dt *DagTasks) GetRootTask() *task {
	return dt.root
}

// ChangeRootTask 改变根节点任务
func (dt *DagTasks) ChangeRootTask(root *task) error {
	if dt.IsReady() || dt.IsRunning() {
		return errors.New("task is running or ready")
	}
	dt.root = root
	dt.IsInit()
	return nil
}

func (dt *DagTasks) startRun(runCh chan *task) {
	dt.doStartRun(dt.root, runCh)
}

func (dt *DagTasks) doStartRun(tk *task, runCh chan *task) {
	if tk == nil {
		return
	}
	for _, temp := range tk.childrenTasks {
		if len(temp.childrenTasks) == 0 {
			temp.Ready()
			runCh <- temp
			continue
		}
		dt.doStartRun(temp, runCh)
	}
}

// doExec 执行任务
func (dt *DagTasks) doExec(ctx context.Context, runChan chan *task) {

	defer dt.wg.Done()
	for {
		tempTk := <-runChan
		if tempTk == endTask {
			return
		}
		if !tempTk.IsReady() {
			continue
		}
		var fn = func() {
			tk := tempTk
			defer func() {
				if e := recover(); e != nil {
					log.Errorf("panic:%v", string(debug.Stack()))
				}
			}()
			tk.Running()
			if tk.retErr == nil {
				tk.retVal, tk.retErr = tk.fn(ctx, TaskHelper{task: tk})
			}
			if tk.retErr != nil {
				log.Errorf("doTask:%s, err:%v", tk.name, tk.retErr)
			}
			tk.Finish()
			if len(tk.parents) == 0 {
				runChan <- endTask
				return
			}
			for parent := range tk.parents {
				if !parent.IsInit() {
					continue
				}
				atomic.AddInt32(&parent.finishCnt, 1)
				if int(atomic.LoadInt32(&parent.finishCnt)) == len(parent.childrenTasks) && parent.casSetStatus(TaskInit, TaskReady) {
					runChan <- parent
				}
			}
		}
		if dt.workerPool == nil {
			go fn()
		} else {
			var poolErr error
			for i := 0; i < maxRetry; i++ {
				if dt.workerPool == nil {
					poolErr = workerpool.Submit(fn)
					if poolErr == nil {
						break
					}
					time.Sleep(time.Millisecond * 200)
				} else {
					poolErr = dt.workerPool.Submit(fn)
					if poolErr == nil {
						break
					}
					time.Sleep(time.Millisecond * 200)
				}
			}
			if poolErr != nil {
				log.FrameworkLogger.Error("workerPool Err %v ", poolErr)
				break
			}
		}

	}
}

func (tk *task) GetSubTaskRet(key string) (interface{}, error) {
	return tk.childrenTasks[key].GetRet()
}

func (th *TaskHelper) GetSubTaskRet(key string) (interface{}, error) {
	return th.task.GetSubTaskRet(key)
}

func (th *TaskHelper) GetSubTaskSize() int {
	return len(th.task.childrenKeyList)
}

func (th *TaskHelper) GetSubTaskNames() []string {
	return th.task.childrenKeyList
}

func (th *TaskHelper) GetSubTaskRetForIndex(index int) (interface{}, error) {
	if index < 0 || index >= len(th.task.childrenKeyList) {
		return nil, nil
	}
	return th.task.GetSubTaskRet(th.task.childrenKeyList[index])
}

func (tk *task) GetRet() (interface{}, error) {
	if !tk.IsFinish() {
		return nil, nil
	}
	return tk.retVal, tk.retErr
}

// Init init状态设置
func (dt *DagTasks) Init() {
	if dt.state == DagInit {
		return
	}
	dt.state = DagInit
}

// Ready ready状态设置
func (dt *DagTasks) Ready() {
	if dt.state == DagReady {
		return
	}
	dt.state = DagReady
}

// Running running状态设置
func (dt *DagTasks) Running() {
	if dt.state == DagRunning {
		return
	}
	dt.state = DagRunning
}

// Finish finish状态设置
func (dt *DagTasks) Finish() {
	if dt.state == DagFinish {
		return
	}
	dt.state = DagFinish
}

// IsInit 是否为初始化状态
func (dt *DagTasks) IsInit() bool {
	return dt.state == DagInit
}

// IsReady 是否为就绪态
func (dt *DagTasks) IsReady() bool {
	return dt.state == DagReady
}

// IsRunning 是否为运行中状态
func (dt *DagTasks) IsRunning() bool {
	return dt.state == DagRunning
}

// IsFinish 是否为运行完成状态
func (dt *DagTasks) IsFinish() bool {
	return dt.state == DagFinish
}

// casSetStatus cas设置状态
func (tk *task) casSetStatus(old, new TaskState) bool {
	return atomic.CompareAndSwapInt32((*int32)(&tk.state), int32(old), int32(new))
}

// Init init状态设置
func (tk *task) Init() {
	if tk.state == TaskInit {
		return
	}
	tk.state = TaskInit
}

// Ready Ready状态设置
func (tk *task) Ready() {
	if tk.state == TaskReady {
		return
	}
	tk.state = TaskReady
}

// Running running状态设置
func (tk *task) Running() {
	if tk.state == TaskRunning {
		return
	}
	tk.state = TaskRunning
}

// Finish finish状态设置
func (tk *task) Finish() {
	if tk.state == TaskFinish {
		return
	}
	tk.state = TaskFinish
}

// IsInit 是否为初始化状态
func (tk *task) IsInit() bool {
	return tk.state == TaskInit
}

// IsReady 是否为就绪态
func (tk *task) IsReady() bool {
	return tk.state == TaskReady
}

// IsRunning 是否为运行中状态
func (tk *task) IsRunning() bool {
	return tk.state == TaskRunning
}

// IsFinish 是否为运行完成状态
func (tk *task) IsFinish() bool {
	return tk.state == TaskFinish
}
