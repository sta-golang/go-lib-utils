package asyncgroup

// TaskStatus 任务状态
type TaskStatus int32

type task struct {
	name   string
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

// Ret 返回值和error
func (t *task) Ret() (interface{}, error) {
	if t.Status != TaskStatusFinish {
		return nil, nil
	}
	return t.retVal, t.retErr
}
