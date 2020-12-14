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
	var rst interface{}
	bt.root,rst = bt.put(bt.root, key, data)
	return rst
}

func (bt *BinaryTree) put(node *binaryTreeNode, key Comparable, data interface{}) (*binaryTreeNode, interface{}) {
	if node == nil {
		return bt.newBinaryTreeNode(key, data), nil
	}
	var oldData interface{}
	cmp := key.CompareTo(node.key)
	if cmp == 0 {
		if node.data != data {
			oldData = node.data
			node.data = data
			return node, oldData
		} else {
			return node, nil
		}
	} else if cmp < 0 {
		node.leftNode, oldData =  bt.put(node.leftNode, key, data)
	} else {
		node.rightNode, oldData = bt.put(node.leftNode, key, data)
	}
	node.size = bt.getSize(node.leftNode) + bt.getSize(node.rightNode) + 1
	return node, oldData
}

func (bt *BinaryTree) newBinaryTreeNode(key Comparable, data interface{}) *binaryTreeNode {
	return &binaryTreeNode{
		key:       key,
		data:      data,
		size:      1,
		leftNode:  nil,
		rightNode: nil,
	}
}

func (bt *BinaryTree) Insert(key Comparable, data interface{}) bool {
	var rst bool
	bt.root, rst = bt.insert(bt.root, key, data)
	return rst
}

func (bt *BinaryTree) insert(node *binaryTreeNode, key Comparable, data interface{}) (*binaryTreeNode, bool) {
	if node == nil {
		return bt.newBinaryTreeNode(key, data),true
	}
	var rst bool
	cmp := key.CompareTo(node.key)
	if cmp == 0 {
		return node, false
	} else if cmp < 0 {
		node.leftNode, rst  = bt.insert(node.leftNode, key, data)
	} else {
		node.rightNode, rst = bt.insert(node.rightNode, key, data)
	}
	node.size = bt.getSize(node.leftNode) + bt.getSize(node.rightNode) + 1
	return node, rst
}

// Get ...
func (bt *BinaryTree) Get(key Comparable) interface{} {
	return bt.get(bt.root, key)
}

func (bt *BinaryTree) get(node *binaryTreeNode, key Comparable) interface{} {
	if node == nil {
		return nil
	}
	cmp := key.CompareTo(node.key)
	if cmp == 0 {
		return node.data
	} else if cmp < 0 {
		return bt.get(node.leftNode, key)
	} else {
		return bt.get(node.rightNode, key)
	}
}

func (bt *BinaryTree) Delete(key Comparable) {
	bt.root = bt.delete(bt.root,key)
}

func (bt *BinaryTree) delete(node *binaryTreeNode, key Comparable) *binaryTreeNode {
	if node == nil {
		return nil
	}
	cmp := key.CompareTo(node.key)
	if cmp == 0 {
		if node.rightNode == nil {
			return node.leftNode
		} else if node.leftNode == nil {
			return node.rightNode
		}
		temp := node
		node = bt.getMin(node.rightNode)
		node.rightNode = bt.delMin(temp.rightNode)
		node.leftNode = temp.leftNode
	} else if cmp < 0 {
		node.leftNode = bt.delete(node.leftNode, key)
	} else if cmp > 0{
		node.rightNode = bt.delete(node.rightNode, key)
	}
	node.size = bt.getSize(node.leftNode) + bt.getSize(node.rightNode) + 1
	return node
}

func (bt *BinaryTree) DelMin() {
	bt.root = bt.delMin(bt.root)
}

func (bt *BinaryTree) delMin(node *binaryTreeNode) *binaryTreeNode {
	if node.leftNode == nil {
		return node.rightNode
	}
	node.leftNode = bt.delMin(node.leftNode)
	node.size = bt.getSize(node.leftNode) + bt.getSize(node.rightNode) + 1
	return node
}

func (bt *BinaryTree) DelMax() {
	bt.root = bt.delMax(bt.root)
}

func (bt *BinaryTree) delMax(node *binaryTreeNode) *binaryTreeNode {
	if node.rightNode == nil {
		return node.leftNode
	}
	node.rightNode = bt.delMax(node.rightNode)
	node.size = bt.getSize(node.leftNode) + bt.getSize(node.rightNode) + 1
	return node
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
