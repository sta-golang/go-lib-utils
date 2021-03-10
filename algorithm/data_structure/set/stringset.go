package main
const (
	defCap = 16
)

type StringSet struct {
	set map[string]struct{}
}

func NewHashSet(size ...int) *StringSet {
	cap := 16
	if len(size) > 0 {
		if size[0] > 16 {
			cap = size[0]
		}
	}
	return &StringSet{set: make(map[string]struct{}, cap)}
}


func (hs *StringSet) Add(elements ...string) {
	for i := range elements {
		hs.set[elements[i]] = emptyElement
	}
}

func (hs *StringSet) Remove(elements ...string) {
	for i := range elements {
		delete(hs.set, elements[i])
	}
}

func (hs *StringSet) Contains(val string) bool {
	if _, ok := hs.set[val]; ok {
		return true
	}
	return false
}

func (hs *StringSet) Size() int {
	return len(hs.set)
}

func (hs *StringSet) Empty() bool {
	return hs.Size() == 0
}

func (hs *StringSet) Clear() {
	hs.set = make(map[string]struct{}, defCap)
}

func (hs *StringSet) Iterator() []string {
	ret := make([]string, 0, len(hs.set))
	for key, _ := range hs.set {
		ret = append(ret, key)
	}
	return ret
}
