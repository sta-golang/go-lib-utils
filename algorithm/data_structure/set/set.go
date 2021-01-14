package set

import "github.com/sta-golang/go-lib-utils/algorithm/data_structure/containers"

type Set interface {
	Add(elements ...interface{})
	Remove(elements ...interface{})
	Contains(val interface{}) bool
	containers.Container
}
