package dag

import (
	"context"

	"sync/atomic"
)

type task struct {
	name            string
	fn              func(ctx context.Context, helper TaskHelper) (interface{}, error)
	retVal          interface{}
	retErr          error
	state           TaskState
	finishCnt       int32
	parents         map[*task]struct{}
	childrenTasks   map[string]*task
	childrenKeyList []string
}

var emptyStruct struct{}

type TaskHelper struct {
	task *task
}

func NewTask(name string, fn func(ctx context.Context, helper TaskHelper) (interface{}, error)) *task {
	if global.taskPool != nil {
		ret := global.taskPool.Get().(*task)
		ret.name = name
		ret.fn = fn
		ret.Clear()
		return ret
	}
	return &task{
		name:          name,
		fn:            fn,
		retVal:        nil,
		retErr:        nil,
		state:         TaskInit,
		parents:       make(map[*task]struct{}),
		childrenTasks: make(map[string]*task),
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

// casSetStatus cas设置状态
func (tk *task) casSetStatus(old, new TaskState) bool {
	return atomic.CompareAndSwapInt32((*int32)(&tk.state), int32(old), int32(new))
}

func (tk *task) AddSubTask(subT *task) {
	if tk.Equals(subT) {
		return
	}
	if _, ok := subT.parents[tk]; ok {
		return
	}
	tk.childrenTasks[subT.name] = subT
	subT.parents[tk] = emptyStruct
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

func (tk *task) Clear() {
	tk.Init()
	tk.retErr = nil
	tk.retVal = nil
	tk.finishCnt = 0
	tk.parents = make(map[*task]struct{})
	tk.childrenTasks = make(map[string]*task)
	tk.childrenKeyList = nil

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
