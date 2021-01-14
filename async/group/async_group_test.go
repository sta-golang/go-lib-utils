package group

import (
	"fmt"
	"testing"
	"time"
)

func TestAsyncGroup_Add(t *testing.T) {
	group := NewAsyncGroup()
	defer group.Close()
	err := group.Add("hello", func() (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println("hello 1")
		return "hello", nil
	})
	err = group.Add("hello", func() (interface{}, error) {
		time.Sleep(time.Millisecond)
		fmt.Println("xixixixi")
		return "2", nil
	})
	if err != nil {
		fmt.Println(err)
		_ = group.Add("hello2", func() (interface{}, error) {
			time.Sleep(time.Millisecond)
			fmt.Println("xixixixi")
			return "2", nil
		})
	}
	_ = group.Add("hello3", func() (interface{}, error) {
		time.Sleep(time.Millisecond * 50)
		fmt.Println("ccccccc")
		return nil, nil
	})
	group.Wait()
	for _, tk := range group.Iterator() {
		fmt.Println(tk.Ret())
	}
	fmt.Println(group.GetTask("hello2").Ret())
	fmt.Println(group.GetTask("hello3").Ret())
}
