package memory

import (
	"sync"
)

//默认配置
const (
	NodeBufSize    = 2048
	PruningRetSize = 2048
	NodeValSize    = 128
)

//左枝二叉树，left node<current node<right node<parent node
//根节点只有左枝，key为节点的过期时间，因为对于内存缓存gc
//我们需要创建一个不断向上生长的树，即过期要删除的节点都在左下方，刚刚添加的节点都在右上方
type expireLeftTree struct {
	mutex    sync.Mutex  //锁
	root     *removeNode //根节点，最大key
	nodePool sync.Pool
}
type removeNode struct {
	key    int64       //key大小
	values []string    //一个key可以包含多个value
	size   int         //values大小
	left   *removeNode //左节点
	right  *removeNode //右节点
}

//创建树
func newTree() *expireLeftTree {
	return &expireLeftTree{
		mutex: sync.Mutex{},
		root:  nil,
		nodePool: sync.Pool{
			New: func() interface{} {
				return &removeNode{}
			},
		},
	}
}

//添加节点
func (tree *expireLeftTree) Add(key int64, val string) {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()
	//比根节点大，重新创建一个根节点
	if tree.root == nil || key > tree.root.key {
		size := 1
		if tree.root != nil {
			size += tree.root.size
		}
		addNode := tree.newNode(key, val, size)
		addNode.left = tree.root
		tree.root = addNode
		return
	}
	//查询到合适的位置
	tree.root = tree.put(tree.root, key, val)
}

func (tree *expireLeftTree) put(node *removeNode, key int64, val string) *removeNode {
	if node == nil {
		return tree.newNode(key, val, 1)
	}
	localSize := 1
	if key < node.key {
		node.left = tree.put(node.left, key, val)
	} else if key > node.key {
		node.right = tree.put(node.right, key, val)
	} else {
		node.values = append(node.values, val)
		localSize++
	}
	node.size = tree.getSize(node.left) + tree.getSize(node.right) + localSize
	return node
}

func (tree *expireLeftTree) getSize(node *removeNode) int {
	if node == nil {
		return 0
	}
	return node.size
}

//修剪树枝
func (tree *expireLeftTree) Pruning(key int64) []string {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()
	//根节点必须存在
	if tree.root == nil {
		return []string{}
	}
	//如果比根节点大，则全量返回
	var ret []string
	if key > tree.root.key {
		tree.traversalDel(tree.root, &ret)
		tree.root = nil
		return ret
	}
	//开始修剪
	tree.root = tree.doPruning(tree.root, key, &ret)
	return ret
}

func (tree *expireLeftTree) levelOrder() [][]removeNode {
	result := make([][]removeNode, 0)
	if tree.root == nil {
		return result
	}
	tree.doLevelOrder(tree.root, 0, &result)
	return result
}

func (tree *expireLeftTree) doLevelOrder(node *removeNode, level int, res *[][]removeNode) {
	if node == nil {
		return
	}
	if len(*res) < level+1 {
		*res = append(*res, make([]removeNode, 0))
	}
	(*res)[level] = append((*res)[level], *node)
	tree.doLevelOrder(node.left, level+1, res)
	tree.doLevelOrder(node.right, level+1, res)
}

func (tree *expireLeftTree) doPruning(node *removeNode, key int64, ret *[]string) *removeNode {
	if node == nil {
		return nil
	}
	if key == node.key {
		initRes(ret, len(node.values)+tree.getSize(node.left))
		tree.traversalDel(node, ret)
		return node.right
	} else if key < node.key {
		node.left = tree.doPruning(node.left, key, ret)
	} else {
		temp := tree.floor(node.right, key)
		if temp == nil {
			initRes(ret, len(node.values)+tree.getSize(node.left))
			tree.traversalDel(node, ret)
			return node.right
		}
		initRes(ret, len(node.values)+tree.getSize(node.left)+len(temp.values)+tree.getSize(temp.left))
		node.right = tree.doPruning(node.right, temp.key, ret)
		tree.traversalDel(node, ret)
		return node.right
	}
	node.size = tree.getSize(node.left) + tree.getSize(node.right) + len(node.values)
	return node
}

func initRes(res *[]string, size int) {
	if res == nil || cap(*res) > 0 {
		return
	}
	*res = make([]string, 0, size)
}

func (tree *expireLeftTree) floor(node *removeNode, key int64) *removeNode {
	if node == nil {
		return nil
	}
	if key == node.key {
		return node
	}
	if key < node.key {
		return tree.floor(node.left, key)
	}
	temp := tree.floor(node.right, key)
	if temp != nil {
		return temp
	}
	return node
}

//遍历删除节点和他的孩子们
func (tree *expireLeftTree) traversalDel(node *removeNode, ret *[]string) {
	if node == nil {
		return
	}
	tree.doTraversalDel(node, ret)
}

func (tree *expireLeftTree) doTraversalDel(node *removeNode, ret *[]string) {
	if node == nil {
		return
	}
	*ret = append(*ret, node.values...)
	tree.doTraversalDel(node.left, ret)
}

//创建node
func (tree *expireLeftTree) newNode(key int64, val string, size int) (retNode *removeNode) {
	//优先从缓存中获取
	retNode, _ = tree.nodePool.Get().(*removeNode)
	//初始化
	retNode.key = key
	retNode.values = []string{val}
	retNode.left = nil
	retNode.right = nil
	retNode.size = size
	return retNode
}

//释放node
func (tree *expireLeftTree) releaseNode(node *removeNode) {
	tree.nodePool.Put(node)
}
