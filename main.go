package main

import (
	"fmt"
	"github.com/xy63237777/go-lib-utils/codec"
)

func main() {
	m := map[string]int{
		"123":1,"345":2,
	}
	bytes, e := codec.API.JsonAPI.Marshal(m)
	fmt.Println(string(bytes),e)

}