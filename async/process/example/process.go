package example

import (
	"github.com/sta-golang/go-lib-utils/async/process"
)

var MyExampleSource = process.NewProcessSource()

func Example() {
	if !process.TryProcess(&MyExampleSource) {
		return
	}
	defer process.EndProcess(&MyExampleSource)
}
