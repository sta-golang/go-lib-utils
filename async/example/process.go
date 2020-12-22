package example

import "github.com/sta-golang/go-lib-utils/async"

var MyExampleSource = async.NewProcessSource()

func Example() {
	if !async.TryProcess(&MyExampleSource) {
		return
	}
	defer async.EndProcess(&MyExampleSource)
}
