package singleexecgroup

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 单独执行，如果key相同 那么所有的任务只有一个会执行
// 非常适缓存去读数据库
type Group struct {
	lock    sync.Mutex
	callMap map[string]*call
}

func (seg *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	seg.lock.Lock()
	if seg.callMap == nil {
		seg.callMap = make(map[string]*call)
	}
	if c, ok := seg.callMap[key]; ok {
		seg.lock.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	seg.callMap[key] = c
	seg.lock.Unlock()

	func() {
		defer func() {
			if err := recover(); err != nil {
			}
		}()
		c.val, c.err = fn()
	}()
	c.wg.Done()

	seg.lock.Lock()
	delete(seg.callMap, key)
	seg.lock.Unlock()

	return c.val, c.err
}
