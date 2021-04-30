package data_structure

var defQueueSize = 16

/**
Queue : 队列
使用双指针首尾指针来辅助队列进行出队入队的操作
出队操作并不删除数据而是将队头指针往后移位
入队操作也是将队尾指针移位如果移动的位置超过数组的长度的时候又回到队头
这样的优点是可以复用整个数组不需要重新扩容

重点: 并不是线程安全的，在多线程下会出现脏数据谨慎操作！
 */
type Queue struct {
	elements  []interface{}
	capSize   int
	headIndex int
	tailIndex int
	size      int
}

func NewQueue() *Queue {
	c := defQueueSize
	return &Queue{
		elements:  make([]interface{}, c),
		capSize:   c,
		headIndex: 0,
		tailIndex: 0,
		size:      0,
	}
}

func NewQueueWithSize(size int) *Queue {
	if size <= 0 {
		size = 16
	}
	return &Queue{
		elements:  make([]interface{}, size),
		capSize:   size,
		headIndex: 0,
		tailIndex: 0,
		size:      0,
	}
}

// 入队 如果队头指针和队尾指针重叠了则说明空间满了 需要扩容
func (q *Queue) Push(data interface{}) {
	if q.size != 0 && q.headIndex == q.tailIndex {
		q.dilatation()
	}
	q.elements[q.tailIndex] = data
	q.tailIndex = (q.tailIndex + 1) % q.capSize
	q.size++
}

func (q *Queue) Empty() bool {
	return q.size == 0
}

func (q *Queue) Pop() interface{} {
	if q.Empty() {
		return nil
	}
	ret := q.elements[q.headIndex]
	q.headIndex = (q.headIndex + 1) % q.capSize
	q.size--
	return ret
}

func (q *Queue) Head() interface{} {
	if q.Empty() {
		return nil
	}
	return q.elements[q.headIndex]
}

// Clean 清空操作也只是重置指针位置
func (q *Queue) Clean() {
	q.headIndex = 0
	q.tailIndex = 0
	q.size = 0
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) dilatation() {
	newArr := make([]interface{}, q.capSize<<1)
	index := q.headIndex
	for i := 0; i < q.capSize; i++ {
		newArr[i] = q.elements[index]
		index = (index + 1) % q.capSize
	}
	q.headIndex = 0
	q.tailIndex = q.capSize
	q.capSize = q.capSize << 1
	q.elements = newArr
}
