package data_structure

import (
	"fmt"
	"testing"
)

type key struct {
	cmp int
}

func (k *key) CompareTo(c Comparable) int {
	k2, _ := c.(*key)
	if k.cmp == k2.cmp {
		return 0
	} else if k.cmp > k2.cmp {
		return 1
	} else {
		return -1
	}
}

func Order(node *binaryTreeNode) {
	if node == nil {
		return
	}
	fmt.Println(node.key, node.data)
	Order(node.leftNode)
	Order(node.rightNode)
}

func TestBinaryTree_GetMin(t *testing.T) {
	testBiTree := &BinaryTree{}
	testBiTree.Insert(&key{3}, 12)
	testBiTree.Insert(&key{1}, 5)
	testBiTree.Insert(&key{2}, 6)
	testBiTree.Put(&key{3}, 99999)
	Order(testBiTree.root)
	fmt.Println("------------")
	testBiTree.Delete(&key{2})
	Order(testBiTree.root)
	fmt.Println("------------")
	testBiTree.Delete(&key{3})
	Order(testBiTree.root)
}
