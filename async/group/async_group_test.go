package group

import (
	"fmt"
	"testing"
	"time"
)

func TestAsyncGroup_Add(t *testing.T) {
	group := NewAsyncGroup(10)
	defer group.Close()
	reqID1 := group.Add(func() (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println("hello 1")
		return "hello", nil
	})
	reqID2 := group.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond)
		fmt.Println("xixixixi")
		return "2", nil
	})
	_ = group.Add(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 50)
		fmt.Println("ccccccc")
		return nil, nil
	})
	group.Wait()
	for _, tk := range group.Iterator() {
		fmt.Println(tk.Ret())
	}
	fmt.Println(group.GetTask(reqID1).Ret())
	fmt.Println(group.GetTask(reqID2).Ret())
}
