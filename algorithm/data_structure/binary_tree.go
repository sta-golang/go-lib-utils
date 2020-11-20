package data_structure

type BinaryTree struct {
	root *binaryTreeNode
}

type binaryTreeNode struct {
	key       Comparable
	data      interface{}
	size      int //这个子树的个数
	leftNode  *binaryTreeNode
	rightNode *binaryTreeNode
}

// Comparable 比较接口小于返回负值 等于返回0值 大于返回正值
type Comparable interface {
	CompareTo(Comparable) int
}

func (bt *BinaryTree) GetSize() int {
	return bt.getSize(bt.root)
}

func (bt *BinaryTree) getSize(node *binaryTreeNode) int {
	if node == nil {
		return 0
	}
	return node.size
}

// Put 返回值为替换后的旧值为nil是没有退化为插入
func (bt *BinaryTree) Put(key Comparable, data interface{}) interface{} {
	return bt.put(bt.root, key, data)
}

func (bt *BinaryTree) put(node *binaryTreeNode, key Comparable, data interface{}) interface{} {
	return nil
}

func (bt *BinaryTree) Insert(key Comparable, data interface{}) bool {
	return bt.insert(bt.root, key, data)
}

func (bt *BinaryTree) insert(node *binaryTreeNode, key Comparable, data interface{}) bool {
	return true
}

// Get ...
func (bt *BinaryTree) Get(key Comparable) interface{} {

	return bt.get(bt.root, key)
}

func (bt *BinaryTree) get(node *binaryTreeNode, key Comparable) interface{} {
	return nil
}

func (bt *BinaryTree) Delete(key Comparable) interface{} {
	return nil
}

func (bt *BinaryTree) GetMin() interface{} {
	if bt.root == nil {
		return nil
	}
	return bt.getMin(bt.root).data
}

func (bt *BinaryTree) getMin(node *binaryTreeNode) *binaryTreeNode {
	if node.leftNode == nil {
		return node
	}
	return bt.getMin(node.leftNode)
}

func (bt *BinaryTree) GetMax() interface{} {
	if bt.root == nil {
		return nil
	}
	return bt.getMin(bt.root).data
}

func (bt *BinaryTree) getMax(node *binaryTreeNode) *binaryTreeNode {
	if node.rightNode == nil {
		return node
	}
	return bt.getMin(node.rightNode)
}

func (bt *BinaryTree) SelectKey(k int) Comparable {
	ret := bt.selectKey(bt.root, k)
	if ret == nil {
		return nil
	}
	return ret.key
}

func (bt *BinaryTree) selectKey(node *binaryTreeNode, k int) *binaryTreeNode {
	if node == nil {
		return nil
	}
	t := bt.getSize(node.leftNode)
	if t > k {
		return bt.selectKey(node.leftNode, k)
	} else if t < k {
		return bt.selectKey(node.rightNode, k-t-1)
	}
	return node
}

func (bt *BinaryTree) Rank(key Comparable) int {
	return bt.rank(bt.root, key)
}

func (bt *BinaryTree) rank(node *binaryTreeNode, key Comparable) int {
	if node == nil {
		return 0
	}
	cmp := key.CompareTo(node.key)
	if cmp < 0 {
		return bt.rank(node.leftNode, key)
	} else if cmp > 0 {
		return 1 + bt.getSize(node.leftNode) + bt.rank(node.rightNode, key)
	}
	return bt.getSize(node.leftNode)
}

func (bt *BinaryTree) Floor(key Comparable) Comparable {
	ret := bt.floor(bt.root, key)
	if ret == nil {
		return nil
	}
	return ret.key
}

func (bt *BinaryTree) floor(node *binaryTreeNode, key Comparable) *binaryTreeNode {
	if node == nil {
		return nil
	}
	cmp := key.CompareTo(node.key)
	if cmp == 0 {
		return node
	}
	if cmp < 0 {
		return bt.floor(node.leftNode, key)
	}
	temp := bt.floor(node.rightNode, key)
	if temp != nil {
		return temp
	}
	return node
}

func (bt *BinaryTree) Ceiling(key Comparable) Comparable {
	ret := bt.ceiling(bt.root, key)
	if ret == nil {
		return nil
	}
	return ret.key
}

func (bt *BinaryTree) ceiling(node *binaryTreeNode, key Comparable) *binaryTreeNode {
	if node == nil {
		return nil
	}
	cmp := key.CompareTo(node.key)
	if cmp == 0 {
		return node
	}
	if cmp > 0 {
		return bt.ceiling(node.rightNode, key)
	}
	temp := bt.ceiling(node.leftNode, key)
	if temp != nil {
		return temp
	}
	return nil
}
