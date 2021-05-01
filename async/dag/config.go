package dag

import (
	"sync"
)

type dagConfig struct {
	taskPool *sync.Pool
	dagPool  *sync.Pool
}

var global = &dagConfig{}

func Config() *dagConfig {
	return global
}

func (dc *dagConfig) SetTaskPool() {
	dc.taskPool = &sync.Pool{
		New: func() interface{} {
			return &task{}
		},
	}
}

func (dc *dagConfig) SetDagPool() {
	dc.dagPool = &sync.Pool{
		New: func() interface{} {
			return &DagTasks{}
		},
	}
}

func (dc *dagConfig) SetPool() {
	dc.SetDagPool()
	dc.SetTaskPool()
}
