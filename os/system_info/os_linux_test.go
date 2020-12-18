//+build !windows

package system_info

import (
	"fmt"
	"testing"
)

func TestUnitSystemInfo_String(t *testing.T) {
	fmt.Println(GetSystemInfo())
}
