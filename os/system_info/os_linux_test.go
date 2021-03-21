//+build !windows

package systeminfo

import (
	"fmt"
	"testing"
)

func TestUnitSystemInfo_String(t *testing.T) {
	fmt.Println(GetSystemInfo())
}
