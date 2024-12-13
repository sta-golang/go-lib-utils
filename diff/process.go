package diff

import (
	"fmt"
	"reflect"

	"github.com/sta-golang/go-lib-utils/log"
)

// Processer 字段处理器
type Processer interface {
	// Load 读取
	Load(source map[string]interface{})
	// Diff diif比较
	Diff(target map[string]interface{}) []*DiffItem
}

// MapProcesser 字段处理器
type MapProcesser struct {
	source map[string]interface{}
}

// NewMapProcesser 构造器
func NewMapProcesser() *MapProcesser {
	return &MapProcesser{}
}

func (mp *MapProcesser) Load(source map[string]interface{}) {
	mp.source = mp.fieldMap(source)
}

func (mp *MapProcesser) Diff(target map[string]interface{}) []*DiffItem {
	return mp.doDiff(mp.fieldMap(target))
}

func (mp *MapProcesser) doDiff(target map[string]interface{}) []*DiffItem {
	ret := make([]*DiffItem, 0)
	for key, val := range mp.source {
		targetVal, ok := target[key]
		if !ok {
			ret = append(ret, &DiffItem{
				Key:         key,
				KeyIsSource: true,
				SourceVal:   val,
				TargetVal:   nil,
			})
			continue
		}
		if reflect.DeepEqual(val, targetVal) {
			continue
		}
		ret = append(ret, &DiffItem{
			Key:         key,
			KeyIsSource: true,
			SourceVal:   val,
			TargetVal:   targetVal,
		})
	}
	for key, val := range target {
		_, ok := mp.source[key]
		if ok {
			continue
		}
		ret = append(ret, &DiffItem{
			Key:         key,
			KeyIsSource: false,
			SourceVal:   nil,
			TargetVal:   val,
		})
	}
	return ret
}

func (mp *MapProcesser) fieldMap(source map[string]interface{}) map[string]interface{} {
	return mp.doMakeFieldMap("", source)
}

func (mp *MapProcesser) doMakeFieldMap(prefix string, source map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range source {
		curKey := key
		if prefix != "" {
			curKey = fmt.Sprintf("%s.%s", prefix, curKey)
		}
		subMap, ok := val.(map[string]interface{})
		if ok {
			mp.mergeMap(ret, mp.doMakeFieldMap(curKey, subMap))
			continue
		}
		ret[curKey] = val
	}
	return ret
}

func (mp *MapProcesser) mergeMap(source, target map[string]interface{}) {
	for key, val := range target {
		if sourceVal, ok := source[key]; ok && !reflect.DeepEqual(sourceVal, val) {
			log.Warn("mergeMap [source] key : {%s} exist val : {%v} [target] : {%v}", key, sourceVal, val)
		}
		source[key] = val
	}
}
