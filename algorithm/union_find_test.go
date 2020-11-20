package algorithm

import (
	"fmt"
	"testing"
)

func TestNewUnionFind(t *testing.T) {
	arr := []interface{}{"hello", "world", "golang"}
	find := NewUnionFind(arr)
	find.Union("hello", "golang")
	//find.Union("hello", "world")
	fmt.Println(find.Count())
	fmt.Println(find.Connected("golang", "world"))
	fmt.Println(find.GetConnects())
}
