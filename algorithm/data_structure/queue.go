package data_structure

var defQueueSize = 16

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
