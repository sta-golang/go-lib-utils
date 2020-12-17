// +build !windows

package system_info

import (
	"github.com/sta-golang/go-lib-utils/os/os_linux"
)

func GetSystemInfo() *SystemInfo {
	return os_linux.GetWindowsSystemInfo()
}
