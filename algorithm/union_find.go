package algorithm

const (
	ufNotFound = -1
)

// UnionFind 带路径压缩的并查集
// 注意： 这个并不是线程安全的如果多线程访问则需要加锁
type UnionFind struct {
	ids   []int
	sz    []int
	count int
	table map[interface{}]int
}

// NewUnionFind 构造函数
// 参数是一个interface的切片比如传入一个字符串切片
// []interface{}{"hello", "world", "golang"}
func NewUnionFind(items []interface{}) *UnionFind {
	uf := &UnionFind{
		ids:   make([]int, len(items)),
		sz:    make([]int, len(items)),
		count: len(items),
		table: make(map[interface{}]int, (len(items)<<2)-len(items)),
	}
	for i := 0; i < len(items); i++ {
		uf.ids[i] = i
		uf.sz[i] = 1
		uf.table[items[i]] = i
	}
	return uf
}

// 有几个连同分量
func (uf *UnionFind) Count() int {
	return uf.count
}

// 两个分量是不是连通的
func (uf *UnionFind) Connected(p, q interface{}) bool {
	pInt := uf.getIDForTable(p)
	qInt := uf.getIDForTable(q)
	if pInt == ufNotFound || qInt == ufNotFound {
		return false
	}
	return uf.find(pInt) == uf.find(qInt)
}

// 合并两个分量
func (uf *UnionFind) Union(p, q interface{}) bool {
	pInt := uf.getIDForTable(p)
	qInt := uf.getIDForTable(q)
	if pInt == ufNotFound || qInt == ufNotFound {
		return false
	}
	pRoot := uf.find(pInt)
	qRoot := uf.find(qInt)

	if pRoot == qRoot {
		return false
	}
	if uf.sz[pRoot] < uf.sz[qRoot] {
		uf.ids[pRoot] = qRoot
		uf.sz[qRoot] += uf.sz[pRoot]
	} else {
		uf.ids[qRoot] = pRoot
		uf.sz[pRoot] += uf.sz[qRoot]
	}
	uf.count--
	return true
}

func (uf *UnionFind) getIDForTable(target interface{}) int {
	if val, ok := uf.table[target]; ok {
		return val
	}
	return ufNotFound
}

func (uf *UnionFind) find(p int) int {
	if p != uf.ids[p] {
		uf.ids[p] = uf.find(uf.ids[p])
	}
	return uf.ids[p]
}
