// @Author : zk_kiger
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
	Value    interface{}
	Priority int
	Index    int
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

// GetByValue get the element by value
func (pq *PriorityQueue) GetByValue(value interface{}) *Element {
	for _, element := range pq.holder.elements {
		if element.Value == value {
			return element
		}
	}
	return nil
}

// GetByIndex get the element by index
func (pq *PriorityQueue) GetByIndex(index int) *Element {
	element := pq.holder.elements[index]
	if element != nil {
		return element
	}
	return nil
}

// Remove
func (pq *PriorityQueue) Remove(value interface{}) *Element {
	element := pq.GetByValue(value)
	heap.Remove(pq.holder, element.Index)
	return element
}

// Len elements length
func (pq *PriorityQueue) Len() int {
	return pq.holder.Len()
}

// Empty elements is empty (Len == 0)
func (pq *PriorityQueue) Empty() bool {
	return pq.holder.Len() == 0
}

// Update update modifies the priority and value of an Item in the queue.
// Note that UPDATE needs to be used in conjunction with TOP,
// or with a known reference to an Element
func (pq *PriorityQueue) Update(element *Element, value string, priority int) {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	element.Value = value
	element.Priority = priority
	heap.Fix(pq.holder, element.Index)
}

// ================= Holder =================
// Push pushes an element to the ElementHolder
func (eh *ElementHolder) Push(x interface{}) {
	n := len(eh.elements)
	element := x.(*Element)
	element.Index = n
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
	element.Index = -1 // for safety
	eh.elements = old[:n-1]
	return element
}

// Len returns the amount of elements in ElementHolder
func (eh *ElementHolder) Len() int {
	return len(eh.elements)
}

// Len compare two elements at position i and j , and returns true if elements[i] < elements[j]
func (eh *ElementHolder) Less(i, j int) bool {
	return eh.elements[i].Priority > eh.elements[j].Priority
}

// Swap swaps two elements at position i and j
func (eh *ElementHolder) Swap(i, j int) {
	eh.elements[i], eh.elements[j] = eh.elements[j], eh.elements[i]
	eh.elements[i].Index, eh.elements[j].Index = i, j
}
