package set

type Set interface {
	Add(elements ...interface{})
	Remove(elements ...interface{})
	Contains(val interface{}) bool
	Size() int
	Empty() bool
	Clear()
}
