package system_info

func GetSystemInfo() *SystemInfo {
	return getWindowsSystemInfo()
}
