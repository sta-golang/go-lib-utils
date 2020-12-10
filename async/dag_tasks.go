package async

import (
	"bytes"
	"context"
	"fmt"
	"github.com/xy63237777/go-lib-utils/log"
	"runtime"
	"runtime/debug"
	"sync"
)

// TaskState 任务的执行状态
type TaskState int

// const 任务执行状态枚举值，0为执行，1执行中，2执行完成
const (
	TaskInit    TaskState = 0
	TaskRunning TaskState = 1
	TaskFinish  TaskState = 2
)

var (
	endTask *task
)

func init() {
	endTask = &task{}
}

type DagTasks struct {
	root *task
}

type task struct {
	name            string
	fn              func(ctx context.Context, helper TaskHelper) (interface{}, error)
	retVal          interface{}
	retErr          error
	state           TaskState
	finishCnt       int
	parents         map[*task]bool
	lock            sync.Mutex
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
	return DagTasks{root: root}
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
	go dt.startRun(runCh)
	dt.doExec(ctx, runCh)
	close(runCh)
	return
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
			runCh <- temp
			continue
		}
		dt.doStartRun(temp, runCh)
	}
}

// doExec 执行任务
func (dt *DagTasks) doExec(ctx context.Context, runChan chan *task) {

	for {
		tempTk := <-runChan
		if tempTk == endTask {
			return
		}
		go func(tk *task) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorf("panic:%v", string(debug.Stack()))
				}
			}()
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
				parent.lock.Lock()
				parent.finishCnt++
				if parent.finishCnt == len(parent.childrenTasks) {
					runChan <- parent
				}
				parent.lock.Unlock()
			}

		}(tempTk)
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
func (tk *task) Init() {
	if tk.state == TaskInit {
		return
	}
	tk.state = TaskInit
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

// IsRunning 是否为运行中状态
func (tk *task) IsRunning() bool {
	return tk.state == TaskRunning
}

// IsFinish 是否为运行完成状态
func (tk *task) IsFinish() bool {
	return tk.state == TaskFinish
}
