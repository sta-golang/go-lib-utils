package main

import (
	"encoding/hex"
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
	fmt.Println((4 << 2) - (4 >> 1))
	pwd := "qq123456"
	//sum := md5.Sum([]byte(pwd))
	str := hex.EncodeToString([]byte(pwd))
	fmt.Println(str)
}
