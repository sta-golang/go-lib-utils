package err

import (
	"github.com/xy63237777/go-lib-utils/algorithm/data_structure"
	"github.com/xy63237777/go-lib-utils/time"
	"strings"
	"sync"
)

var finished bool
var initErrBufferOnce sync.Once
var globalBuffer errBuffer

type errBuffer struct {
	isAsync  bool
	errTable map[string]errorExtra
	list     *data_structure.LinkedList
	help     helper
}

func InitErrBuffer(isAsync bool) {
	if finished {
		return
	}
	initErrBufferOnce.Do(func() {
		finished = true
		globalBuffer = errBuffer{
			isAsync:  isAsync,
			errTable: make(map[string]errorExtra),
			list:     data_structure.NewLinkedList(),
			help:     nil,
		}
		if isAsync {
			globalBuffer.help = &asyncHelp{
				setCh:  make(chan *errorExtra, 128),
				target: &globalBuffer,
			}
		} else {
			globalBuffer.help = &syncHelp{
				target: &globalBuffer,
				lock:   sync.RWMutex{},
			}
		}
	})
}

func PutErr(err Error, extra ...string) {
	ee := errorExtra{
		err: err,
		extras: []timeAndExtra{
			{
				tmStr:    time.GetNowDateTimeStr(),
				extraMsg: strings.Join(extra, " ; "),
			},
		},
	}
	globalBuffer.help.writerErr(ee)
}

func DelErr(err Error) {

}

type errorExtra struct {
	err    Error
	extras []timeAndExtra
}

type timeAndExtra struct {
	tmStr    string
	extraMsg string
}

type helper interface {
	writerErr(err errorExtra)
}

type asyncHelp struct {
	setCh  chan *errorExtra
	target *errBuffer
}

func (ah *asyncHelp) writerErr(err errorExtra) {
	ah.setCh <- &err
}

type syncHelp struct {
	target *errBuffer
	lock   sync.RWMutex
}

func (sh *syncHelp) writerErr(err errorExtra) {
	sh.lock.Lock()
	sh.target.doWriterErrForTable(err)
	sh.lock.Unlock()
}

func (eb *errBuffer) doWriterErrForTable(err errorExtra) {
	key := err.err.getKey()
	if val, ok := eb.errTable[key]; ok {
		val.extras = append(val.extras, err.extras...)
		return
	}
	eb.errTable[key] = err
	eb.list.Add(key)
}
