package data_structure

type Stack struct {
	cnt      int
	stackCap int
	elements []interface{}
}

func NewStack() *Stack {
	c := 16
	return &Stack{
		cnt:      0,
		stackCap: c,
		elements: make([]interface{}, c),
	}
}

func NewStackWithSize(size int) *Stack {
	if size <= 0 {
		size = 16
	}
	return &Stack{
		cnt:      0,
		stackCap: size,
		elements: make([]interface{}, size),
	}
}

func (s *Stack) Empty() bool {
	return s.cnt == 0
}

func (s *Stack) Size() int {
	return s.cnt
}

func (s *Stack) Peek() interface{} {
	if s.Empty() {
		return nil
	}
	return s.elements[s.cnt-1]
}

func (s *Stack) Pop() interface{} {
	if s.Empty() {
		return nil
	}
	ret := s.elements[s.cnt-1]
	s.cnt--
	return ret
}

func (s *Stack) Push(data interface{}) {
	if s.cnt == s.stackCap {
		s.dilatation()
	}
	s.elements[s.cnt] = data
	s.cnt++
}

func (s *Stack) dilatation() {
	newArr := make([]interface{}, s.stackCap<<1)
	for i := 0; i < s.stackCap; i++ {
		newArr[i] = s.elements[i]
	}
	s.stackCap = s.stackCap << 1
	s.elements = newArr
}
