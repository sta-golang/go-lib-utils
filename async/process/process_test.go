package process

import (
	"fmt"
	"testing"
)

func TestEndProcess(t *testing.T) {
	myFlag := NewProcessSource()
	kk := NewProcessSource()
	fmt.Println(TryProcess(&myFlag))
	defer EndProcess(&kk)
}
