package set

import "github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"

const (
	defCap = 16
)

type HashSet struct {
	set map[interface{}]struct{}
}

func NewHashSet(size ...int) *HashSet {
	cap := 16
	if len(size) > 0 {
		if size[0] > 16 {
			cap = size[0]
		}
	}
	return &HashSet{set: make(map[interface{}]struct{}, cap)}
}


func (hs *HashSet) Add(elements ...interface{}) {
	for i := range elements {
		hs.set[elements[i]] = onceElement
	}
}

func (hs *HashSet) Remove(elements ...interface{}) {
	for i := range elements {
		delete(hs.set, elements[i])
	}
}

func (hs *HashSet) Contains(val interface{}) bool {
	if _, ok := hs.set[val]; ok {
		return true
	}
	return false
}

func (hs *HashSet) Size() int {
	return len(hs.set)
}

func (hs *HashSet) Empty() bool {
	return hs.Size() == 0
}

func (hs *HashSet) Clear() {
	hs.set = make(map[interface{}]struct{}, defCap)
}

func (hs *HashSet) Iterator() []interface{} {
	ret := make([]interface{}, 0, len(hs.set))
	for key, _ := range hs.set {
		ret = append(ret, key)
	}
	return ret
}
