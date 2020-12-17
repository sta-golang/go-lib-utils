//+build !windows

package os_linux

import (
	"fmt"
	"testing"
)

func TestUnitSystemInfo_String(t *testing.T) {
	fmt.Println(GetSystemInfo())
}
