package datastructure

type linkedListNode struct {
	prev *linkedListNode
	next *linkedListNode
	Data interface{}
}

type LinkedList struct {
	size int
	head *linkedListNode
	tail *linkedListNode
}

// NewLinkedList 创建函数
func NewLinkedList() *LinkedList {
	head := &linkedListNode{}
	tail := &linkedListNode{}
	head.next = tail
	tail.prev = head
	return &LinkedList{
		head: head,
		tail: tail,
		size: 0,
	}
}

// Size 大小
func (ll *LinkedList) Size() int {
	return ll.size
}

// IsEmpty 是否为空
func (ll *LinkedList) IsEmpty() bool {
	return ll.head.next == ll.tail
}

// GetHead 获取头结点的数据
func (ll *LinkedList) GetHead() interface{} {
	if ll.Size() == 0 {
		return nil
	}
	return ll.head.next.Data
}

// AddFirst 添加一个头结点
func (ll *LinkedList) AddFirst(data interface{}) {
	if data == nil {
		return
	}
	node := &linkedListNode{Data: data, next: ll.head.next, prev: ll.head}
	ll.head.next.prev = node
	ll.head.next = node
	ll.size++
}

// RemoveHead 把头结点删除
func (ll *LinkedList) RemoveHead() (interface{}, bool) {
	if ll.Size() == 0 {
		return nil, false
	}
	node := ll.head.next
	node.next.prev = ll.head
	ll.head.next = node.next
	node.next = nil
	node.prev = nil
	return node.Data, true
}

// RemoveTail 将最后一个删除掉
func (ll *LinkedList) RemoveTail() (interface{}, bool) {
	if ll.Size() == 0 {
		return nil, false
	}
	node := ll.tail.prev
	node.prev.next = ll.tail
	ll.tail.prev = node.prev
	node.prev = nil
	node.next = nil
	return node.Data, true
}

// Add 尾插法
func (ll *LinkedList) Add(data interface{}) {
	if data == nil {
		return
	}
	node := &linkedListNode{Data: data, next: ll.tail, prev: ll.tail.prev}
	ll.tail.prev.next = node
	ll.tail.prev = node
	ll.size++
}

// Iterator 获取迭代的切片
func (ll *LinkedList) Iterator() []interface{} {
	ret := make([]interface{}, ll.size)
	i := 0
	for temp := ll.head.next; temp.next != nil; temp = temp.next {
		ret[i] = temp.Data
		i++
	}
	return ret
}

// Get 使用下标来获取数据支持负数
// 下标从0开始。负数则从-1开始 -1代表最后一个 -2代表倒数第二个
func (ll *LinkedList) Get(index int) interface{} {
	if (index > 0 && index >= ll.size) || (index < 0 && -index > ll.size) {
		return nil
	}
	if index < 0 {
		temp := ll.tail.prev
		end := -index
		for i := 1; i < end; i++ {
			temp = temp.prev
		}
		return temp.Data
	}
	temp := ll.head.next
	for i := 0; i < index; i++ {
		temp = temp.next
	}
	return temp.Data
}

func (ll *LinkedList) Clean() {
	ll.size = 0
	ll.head.next = ll.tail
	ll.tail.prev = ll.head
}
