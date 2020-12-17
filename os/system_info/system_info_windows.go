package system_info

import (
	"github.com/sta-golang/go-lib-utils/os/os_windows"
)

func GetSystemInfo() *SystemInfo {
	return os_windows.GetWindowsSystemInfo()
}
