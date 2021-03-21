package datastructure

import (
	"fmt"
	"testing"
)

func TestLinkedList(t *testing.T) {
	list := NewLinkedList()
	list.Add("hello")
	list.Add("world")
	list.Add("golang")
	fmt.Println(list.Iterator())
	fmt.Println(list.Get(0))
	fmt.Println(list.Get(2))
	fmt.Println(list.Get(3))
	fmt.Println(list.Get(-1))
	fmt.Println(list.Get(-2))
	fmt.Println(list.Get(-3))
	fmt.Println(list.Get(-4))
	fmt.Println(list.Iterator())
	list.Clean()
	list.Add("world")
	list.Add("golang")
	fmt.Println(list.Iterator())
}
