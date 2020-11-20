package str

import (
	"fmt"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	str := "hello"
	fmt.Println(BytesToString(StringToBytes(&str)))
}
