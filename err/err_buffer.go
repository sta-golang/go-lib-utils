// err_buffer 是一个err的缓冲
// 方便一次性返回数据，但是性能可能不是太高
// 这个是为了方便统计。 如果有更好的想法可以联系微
package err

import (
	"fmt"
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
	errTable map[string]*errorExtra
	list     *data_structure.LinkedList
	help     helper
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
	writerErr(err *errorExtra)
	clean() []string
}

func InitErrBuffer(isAsync bool) {
	if finished {
		return
	}
	initErrBufferOnce.Do(func() {
		finished = true
		globalBuffer = errBuffer{
			isAsync:  isAsync,
			errTable: make(map[string]*errorExtra),
			list:     data_structure.NewLinkedList(),
			help:     nil,
		}
		if isAsync {
			globalBuffer.help = newAsyncHelp(&globalBuffer)
		} else {
			globalBuffer.help = newSyncHelp(&globalBuffer)
		}
	})
}

func PutErr(err Error, extra ...string) {
	ee := &errorExtra{
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

func Clean() []string {
	return globalBuffer.help.clean()
}

func (eb *errBuffer) doWriterErrForTable(err *errorExtra) {
	key := err.err.getKey()
	if val, ok := eb.errTable[key]; ok {
		val.extras = append(val.extras, err.extras...)
		return
	}
	eb.errTable[key] = err
	eb.list.Add(key)
}

func (eb *errBuffer) iterator() []string {
	keys := eb.list.Iterator()
	ret := make([]string, 0, eb.list.Size())
	for i := range keys {
		key := keys[i].(string)
		err := eb.errTable[key]
		ret = append(ret, err.String())
	}
	return ret
}

func (eb *errBuffer) clean() []string {
	ret := eb.iterator()
	eb.errTable = make(map[string]*errorExtra)
	eb.list.Clean()
	return ret
}

func (ee *errorExtra) String() string {
	buff := strings.Builder{}
	buff.WriteString(fmt.Sprintf("* error : code : %d msg : %v", ee.err.Code, ee.err.Err))
	for i := 0; i < len(ee.extras); i++ {
		buff.WriteString(fmt.Sprintf("\n\t\t -- [%s] : %s", ee.extras[i].tmStr, ee.extras[i].extraMsg))
	}
	return buff.String()
}

type asyncHelp struct {
	setCh       chan *errorExtra
	cleanFlagCh chan bool
	errMsgCh    chan []string
	target      *errBuffer
}

func newAsyncHelp(target *errBuffer) helper {
	ah := &asyncHelp{
		setCh:       make(chan *errorExtra, 128),
		cleanFlagCh: make(chan bool),
		errMsgCh:    make(chan []string),
		target:      target,
	}
	go ah.doWorker()
	return ah
}

func (ah *asyncHelp) doWorker() {
	for {
		select {
		case err := <-ah.setCh:
			ah.target.doWriterErrForTable(err)
		case <-ah.cleanFlagCh:
			ah.errMsgCh <- ah.target.clean()
		}
	}
}

func (ah *asyncHelp) clean() []string {
	ah.cleanFlagCh <- true
	return <-ah.errMsgCh
}

func (ah *asyncHelp) writerErr(err *errorExtra) {
	ah.setCh <- err
}

type syncHelp struct {
	target *errBuffer
	lock   sync.Mutex
}

func newSyncHelp(target *errBuffer) helper {
	return &syncHelp{
		target: target,
		lock:   sync.Mutex{},
	}
}

func (sh *syncHelp) writerErr(err *errorExtra) {
	sh.lock.Lock()
	sh.target.doWriterErrForTable(err)
	sh.lock.Unlock()
}

func (sh *syncHelp) clean() []string {
	sh.lock.Lock()
	defer sh.lock.Unlock()
	return sh.target.clean()
}
