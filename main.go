package main

import (
	"fmt"
	"github.com/xy63237777/go-lib-utils/algorithm/data_structure"
	"github.com/xy63237777/go-lib-utils/codec"
)

func main() {
	m := map[string]int{
		"123": 1, "345": 2,
	}
	bytes, e := codec.API.JsonAPI.Marshal(m)
	fmt.Println(string(bytes), e)
	list := data_structure.NewLinkedList()
	list.Add("hello")
	list.Add("world")
	list.Add("golang")

}
