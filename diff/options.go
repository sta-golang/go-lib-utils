package diff

type options struct {
	processer Processer
}

// Option 可选参数
type Option func(*options)

var DefaultOption = func() *options {
	return &options{
		processer: NewMapProcesser(),
	}
}

// WithProcesser 设置字段处理
func WithProcesser(p Processer) Option {
	return func(o *options) {
		if p == nil {
			return
		}
		o.processer = p
	}
}
