package datastructure

import (
	"fmt"
	"testing"
)

func TestNewQueue(t *testing.T) {
	defQueueSize = 2
	q := NewQueue()
	q.Push("hello")
	q.Push("world")
	q.Push("kkk")
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	q.Push("123123")
	q.Push("asdasd")
	q.Push("555")
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
}
