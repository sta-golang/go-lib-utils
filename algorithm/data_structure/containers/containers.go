package containers

type Container interface {
	Size() int
	Empty() bool
	Clear()
	Iterator() []interface{}
}
