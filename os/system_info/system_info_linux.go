// +build !windows

package systeminfo

func GetSystemInfo() *SystemInfo {
	return getWindowsSystemInfo()
}
