package datastructure

import (
	"fmt"
	"testing"
)

func TestPriorityQueueNoLock(t *testing.T) {
	// 测试无锁版本
	priorityQueue := NewPriorityQueue()
	element1 := &Element{
		value:    "192.10.1.4",
		priority: 10,
	}
	element2 := &Element{
		value:    "192.111.1.4",
		priority: 4,
	}
	element3 := &Element{
		value:    "192.10.1.4",
		priority: 8,
	}
	element4 := &Element{
		value:    "192.10.1.111",
		priority: 8,
	}
	priorityQueue.Push(element1)
	priorityQueue.Push(element2)
	priorityQueue.Push(element3)
	priorityQueue.Push(element4)
	fmt.Println("top: ", priorityQueue.Top())
	priorityQueue.Update(element2, "test update", 11)
	fmt.Println("top: ", priorityQueue.Top())
	fmt.Println(priorityQueue.Pop())
	fmt.Println(priorityQueue.Pop())
	fmt.Println(priorityQueue.Pop())
}

func TestPriorityQueueWithLock(t *testing.T) {
	// 测试锁版本
	priorityQueue := NewPriorityQueue(WithGoroutineSafe())
	element1 := &Element{
		value:    "192.10.1.11",
		priority: 0,
	}
	element2 := &Element{
		value:    "192.111.1.22",
		priority: 8,
	}
	element3 := &Element{
		value:    "192.10.1.33",
		priority: 8,
	}
	element4 := &Element{
		value:    "192.10.1.44",
		priority: 8,
	}
	priorityQueue.Push(element1)
	priorityQueue.Push(element2)
	priorityQueue.Push(element3)
	priorityQueue.Push(element4)
	fmt.Println("top: ", priorityQueue.Top())
	top := priorityQueue.Top()
	priorityQueue.Update(top, "test update", 9)
	fmt.Println("top: ", priorityQueue.Top())
	fmt.Println(priorityQueue.Pop())
	//fmt.Println("top: ", priorityQueue.Top())
	//fmt.Println(priorityQueue.Pop())
	//fmt.Println("top: ", priorityQueue.Top())
	//fmt.Println(priorityQueue.Pop())
}
