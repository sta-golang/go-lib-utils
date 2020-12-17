package os_windows

import (
	si "github.com/sta-golang/go-lib-utils/os/system_info"
	"github.com/sta-golang/go-lib-utils/server"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	kernel                  = syscall.NewLazyDLL("Kernel32.dll")
	ProcGetSystemTimes      = kernel.NewProc("GetSystemTimes")
	GlobalMemoryStatusEx    = kernel.NewProc("GlobalMemoryStatusEx")
	GetDiskFreeSpaceExW     = kernel.NewProc("GetDiskFreeSpaceExW")
	GetLogicalDriveStringsW = kernel.NewProc("GetLogicalDriveStringsW")
)

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 //物理内存总量
	ullAvailPhys            uint64 //可用物理内存
	ullTotalPageFile        uint64 //页交换文件最多能放的字节数
	ullAvailPageFile        uint64 //页交换文件中尚未分配给进程的字节数
	ullTotalVirtual         uint64 //用户区总的虚拟地址空间
	ullAvailVirtual         uint64 //用户区当前可用的虚拟地址空间
	ullAvailExtendedVirtual uint64
}

type win32_SystemProcessorPerformanceInformation struct {
	IdleTime       int64 // idle time in 100ns (this is not a filetime).
	KernelTime     int64 // kernel time in 100ns.  kernel time includes idle time. (this is not a filetime).
	UserTime       int64 // usertime in 100ns (this is not a filetime).
	DpcTime        int64 // dpc time in 100ns (this is not a filetime).
	InterruptTime  int64 // interrupt time in 100ns
	InterruptCount uint32
}

func MemoryUsageForWindows() (used, free, total uint64) {
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	_, _, _ = GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	return memInfo.ullTotalPhys - memInfo.ullAvailPhys, memInfo.ullAvailPhys, memInfo.ullTotalPhys
}

func usage(getDiskFreeSpaceExW *syscall.LazyProc, path string) (diskusage, error) {
	lpFreeBytesAvailable := int64(0)
	var info = diskusage{Path: path}
	fromString, err := syscall.UTF16PtrFromString(info.Path)
	if err != nil {
		return info, err
	}
	diskret, _, err := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(fromString)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&(info.Total))),
		uintptr(unsafe.Pointer(&(info.Free))))
	if diskret != 0 {
		err = nil
	}
	return info, err
}

type diskusage struct {
	Path  string `json:"path"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

//硬盘信息
func DiskUsageForWindows() (infos []diskusage) {

	lpBuffer := make([]byte, 254)
	diskret, _, _ := GetLogicalDriveStringsW.Call(
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&lpBuffer[0])))
	if diskret == 0 {
		return
	}
	for _, v := range lpBuffer {
		if v >= 65 && v <= 90 {
			path := string(v) + ":"
			if path == "A:" || path == "B:" {
				continue
			}
			info, err := usage(GetDiskFreeSpaceExW, string(v)+":")
			if err != nil {
				continue
			}
			infos = append(infos, info)
		}
	}
	return infos
}

type fileTiMe struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

func CPUUsageForWindows() (user, idle, total uint64) {
	var lpIdleTime fileTiMe
	var lpKernelTime fileTiMe
	var lpUserTime fileTiMe

	_, _, _ = ProcGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&lpIdleTime)),
		uintptr(unsafe.Pointer(&lpKernelTime)),
		uintptr(unsafe.Pointer(&lpUserTime)),
	)
	user = (uint64(lpUserTime.DwHighDateTime) << 32) + uint64(lpUserTime.DwLowDateTime)
	idle = (uint64(lpIdleTime.DwHighDateTime) << 32) + uint64(lpIdleTime.DwLowDateTime)
	ker := (uint64(lpKernelTime.DwHighDateTime) << 32) + uint64(lpKernelTime.DwLowDateTime)
	total = user + idle + ker
	return
}

// GetLinuxSystemInfo 获取系统运行信息.
func GetWindowsSystemInfo() *si.SystemInfo {
	//运行时信息
	mStat := &runtime.MemStats{}
	runtime.ReadMemStats(mStat)

	//CPU信息
	cpuUser, cpuIdel, cpuTotal := CPUUsageForWindows()
	cpuUserRate := float64(cpuUser) / float64(cpuTotal)
	cpuFreeRate := float64(cpuIdel) / float64(cpuTotal)

	//磁盘空间信息
	diskInfos := DiskUsageForWindows()
	di := si.DiskInfo{
		DiskUsed:  0,
		DiskFree:  0,
		DiskTotal: 0,
		Children:  make([]si.ChildrenInfo, 0, len(diskInfos)),
	}
	for i := range diskInfos {
		info := diskInfos[i]
		cl := si.ChildrenInfo{
			Path:      info.Path,
			DiskUsed:  info.Total - info.Free,
			DiskFree:  info.Free,
			DiskTotal: info.Total,
		}
		di.DiskUsed += info.Total - info.Free
		di.DiskFree += info.Free
		di.DiskTotal += info.Total
		di.Children = append(di.Children, cl)
	}
	//内存使用信息
	memUsed, memFree, memTotal := MemoryUsageForWindows()

	serverName := server.ServerName

	return &si.SystemInfo{
		ServerName: serverName,
		SystemOs:   runtime.GOOS,

		Runtime:      int64(server.ServiceUptime()),
		GoroutineNum: runtime.NumGoroutine(),
		CPUNum:       runtime.NumCPU(),
		CPUUser:      cpuUserRate,
		CPUFree:      cpuFreeRate,
		DiskInfo:     di,
		MemUsed:      memUsed,
		MemSys:       mStat.Sys,
		MemFree:      memFree,
		MemTotal:     memTotal,
		AllocGolang:  mStat.Alloc,
		AllocTotal:   mStat.TotalAlloc,
		Lookups:      mStat.Lookups,
		Mallocs:      mStat.Mallocs,
		Frees:        mStat.Frees,
		LastGCTime:   mStat.LastGC,
		NextGC:       mStat.NextGC,
		PauseTotalNs: mStat.PauseTotalNs,
		PauseNs:      mStat.PauseNs[(mStat.NumGC+255)%256],
	}
}
