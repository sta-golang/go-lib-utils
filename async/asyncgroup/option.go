package asyncgroup

import "github.com/sta-golang/go-lib-utils/pool/workerpool"

// Option 选项
type Option func(*Group)

// WithWorkPoolWithSize 使用线程池并指定大小
func WithWorkPoolWithSize(size int) func(*Group) {
	return func(ag *Group) {
		ag.executor = workerpool.New(size)
	}
}

// WithWorkPool 使用指定线程池
func WithWorkPool(wp *workerpool.WorkerPool) func(*Group) {
	return func(ag *Group) {
		ag.executor = wp
	}
}

// WithConcurrentSecurity 设置并发安全.如果此任务组需要共享则可使用此选项
func WithConcurrentSecurity() func(*Group) {
	return func(ag *Group) {
		ag.lockFlag = true
	}
}

// WithTaskSize 设置任务数量大小，如果可以提前预知则可以减少map的扩容次数
func WithTaskSize(size int) func(*Group) {
	return func(ag *Group) {
		ag.tasks = make(map[string]*task, size<<1)
	}
}
