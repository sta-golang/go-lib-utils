package memory

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	tree := newTree()
	tree.add(6, "hello6")
	tree.add(3, "hello3")
	tree.add(5, "hello5")
	tree.add(2, "hello2")
	tree.add(8, "hello8")
	tree.add(7, "hello7")
	for i, nodes := range tree.levelOrder() {
		fmt.Println("i = ", i, " nodes ", nodes)
	}
	fmt.Println(tree.pruning(4))
	for i, nodes := range tree.levelOrder() {
		fmt.Println("i = ", i, " nodes ", nodes)
	}
	fmt.Println(tree.pruning(4))
	tree.add(7, "hello7")
	for i, nodes := range tree.levelOrder() {
		fmt.Println("i = ", i, " nodes ", nodes)
	}

	fmt.Println("----------------------------")

	tree = newTree()
	tree.add(87, "hello87")
	tree.add(46, "hello46")
	tree.add(55, "hello55")
	tree.add(48, "hello48")
	tree.add(66, "hello66")
	tree.add(32, "hello32")
	tree.add(50, "hello50")
	for i, nodes := range tree.levelOrder() {
		fmt.Println("i = ", i, " nodes ", nodes)
	}
	fmt.Println(tree.pruning(49))
	for i, nodes := range tree.levelOrder() {
		fmt.Println("i = ", i, " nodes ", nodes)
	}
	var strsss []string
	initRes(&strsss, 5)
	fmt.Println(cap(strsss))
}
