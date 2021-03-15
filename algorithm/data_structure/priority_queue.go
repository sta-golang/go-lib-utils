package data_structure

import (
	"container/heap"
	"sync"
)

var (
	defaultLocker FakeLocker
)

type FakeLocker struct {
}

// Lock does nothing
func (l FakeLocker) Lock() {

}

// Unlock does nothing
func (l FakeLocker) Unlock() {

}

type Element struct {
	value    interface{}
	priority int
	index    int
}

type ElementHolder struct {
	elements []*Element
}

type PriorityQueue struct {
	holder *ElementHolder
	locker sync.Locker
}

// Options holds PriorityQueue's options
type Options struct {
	locker sync.Locker
}

// Option is a function type used to set Options
type Option func(option *Options)

// WithGoroutineSafe is used to set the PriorityQueue goroutine-safe
func WithGoroutineSafe() Option {
	return func(option *Options) {
		option.locker = &sync.RWMutex{}
	}
}

// ================= API =================
// New creates a PriorityQueue
func NewPriorityQueue(opts ...Option) *PriorityQueue {
	option := Options{
		locker: defaultLocker,
	}
	for _, opt := range opts {
		opt(&option)
	}
	holder := &ElementHolder{
		elements: make([]*Element, 0, 0),
	}
	return &PriorityQueue{
		holder: holder,
		locker: option.locker,
	}
}

// Push pushes an element to the PriorityQueue
func (pq *PriorityQueue) Push(x *Element) {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	heap.Push(pq.holder, x)
}

// Pop pops an element from the PriorityQueue
func (pq *PriorityQueue) Pop() *Element {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	return heap.Pop(pq.holder).(*Element)
}

// Top get top element from the PriorityQueue
func (pq *PriorityQueue) Top() *Element {
	if pq.holder.Len() == 0 {
		return nil
	}
	return pq.holder.elements[0]
}

// Update update modifies the priority and value of an Item in the queue.
// Note that UPDATE needs to be used in conjunction with TOP,
// or with a known reference to an Element
func (pq *PriorityQueue) Update(element *Element, value string, priority int) {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	element.value = value
	element.priority = priority
	heap.Fix(pq.holder, element.index)
}

// ================= Holder =================
// Push pushes an element to the ElementHolder
func (eh *ElementHolder) Push(x interface{}) {
	n := len(eh.elements)
	element := x.(*Element)
	element.index = n
	eh.elements = append(eh.elements, element)
}

// Pop pops an element from the ElementHolder
func (eh *ElementHolder) Pop() interface{} {
	if len(eh.elements) == 0 {
		return nil
	}
	old := eh.elements
	n := eh.Len()
	element := old[n-1]
	old[n-1] = nil     // avoid memory leak
	element.index = -1 // for safety
	eh.elements = old[:n-1]
	return element
}

// Len returns the amount of elements in ElementHolder
func (eh *ElementHolder) Len() int {
	return len(eh.elements)
}

// Len compare two elements at position i and j , and returns true if elements[i] < elements[j]
func (eh *ElementHolder) Less(i, j int) bool {
	return eh.elements[i].priority > eh.elements[j].priority
}

// Swap swaps two elements at position i and j
func (eh *ElementHolder) Swap(i, j int) {
	eh.elements[i], eh.elements[j] = eh.elements[j], eh.elements[i]
	eh.elements[i].index, eh.elements[j].index = i, j
}
