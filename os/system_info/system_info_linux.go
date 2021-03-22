// +build !windows

package system_info

func GetSystemInfo() *SystemInfo {
	return getLinuxSystemInfo()
}
