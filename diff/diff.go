package diff

type diffMap struct {
	o *options
}

// DiffItem diff具体项
type DiffItem struct {
	// 有diff的key
	Key string
	// 这个diff的key是否为source的key
	KeyIsSource bool
	// source 的val
	SourceVal interface{}
	// target 的val
	TargetVal interface{}
}

// New 构造器
var New = func(source map[string]interface{}, opts ...Option) *diffMap {
	o := DefaultOption()
	for _, opt := range opts {
		opt(o)
	}
	dm := &diffMap{
		o: o,
	}
	dm.o.processer.Load(source)
	return dm
}

func (d *diffMap) Diff(target map[string]interface{}) []*DiffItem {
	return d.o.processer.Diff(target)
}
