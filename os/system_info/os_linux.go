//+build !windows

package system_info

import (
	"github.com/sta-golang/go-lib-utils/server"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// 目前只支持linux
// CPUUsage 获取CPU使用率(仅支持linux),单位jiffies(节拍数).
// user为用户态(用户进程)的运行时间,
// idle为空闲时间,
// total为累计时间.
func CPUUsage() (user, idle, total uint64) {
	contents, _ := ioutil.ReadFile("/proc/stat")

	if len(contents) > 0 {
		lines := strings.Split(string(contents), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if fields[0] == "cpu" {
				//CPU指标：user，nice, system, idle, iowait, irq, softirq
				// cpu  130216 19944 162525 1491240 3784 24749 17773 0 0 0

				numFields := len(fields)
				for i := 1; i < numFields; i++ {
					val, _ := strconv.ParseUint(fields[i], 10, 64)
					total += val // tally up all the numbers to get total ticks
					if i == 1 {
						user = val
					} else if i == 4 { // idle is the 5th field in the cpu line
						idle = val
					}
				}
				break
			}
		}
	}

	return
}

// DiskUsage 获取磁盘/目录使用情况,单位字节.参数path为目录.
// used为已用,
// free为空闲,
// total为总数.
func DiskUsage(path string) (used, free, total uint64) {
	fs := &syscall.Statfs_t{}
	err := syscall.Statfs(path, fs)

	if err == nil {
		total = fs.Blocks * uint64(fs.Bsize)
		free = fs.Bfree * uint64(fs.Bsize)
		used = total - free
	}

	return
}

// MemoryUsage 获取内存使用率(仅支持linux),单位字节.
// 参数 virtual,是否取虚拟内存.
// used为已用,
// free为空闲,
// total为总数.
func MemoryUsage(virtual bool) (used, free, total uint64) {
	if virtual {
		// 虚拟机的内存
		contents, err := ioutil.ReadFile("/proc/meminfo")
		if err == nil {
			lines := strings.Split(string(contents), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) == 3 {
					val, _ := strconv.ParseUint(fields[1], 10, 64) // kB

					if strings.HasPrefix(fields[0], "MemTotal") {
						total = val * 1024
					} else if strings.HasPrefix(fields[0], "MemFree") {
						free = val * 1024
					}
				}
			}

			//计算已用内存
			used = total - free
		}
	} else {
		// 真实物理机内存
		sysi := &syscall.Sysinfo_t{}

		err := syscall.Sysinfo(sysi)
		if err == nil {
			total = sysi.Totalram * uint64(syscall.Getpagesize()/1024)
			free = sysi.Freeram * uint64(syscall.Getpagesize()/1024)
			used = total - free
		}
	}

	return
}

// getLinuxSystemInfo 获取系统运行信息.
func getLinuxSystemInfo() *SystemInfo {
	//运行时信息
	mStat := &runtime.MemStats{}
	runtime.ReadMemStats(mStat)

	//CPU信息
	cpuUser, cpuIdel, cpuTotal := CPUUsage()
	cpuUserRate := float64(cpuUser) / float64(cpuTotal)
	cpuFreeRate := float64(cpuIdel) / float64(cpuTotal)

	//磁盘空间信息
	diskUsed, diskFree, diskTotal := DiskUsage("/")
	di := DiskInfo{
		DiskUsed:  diskUsed,
		DiskFree:  diskFree,
		DiskTotal: diskTotal,
		Children: []ChildrenInfo{
			{
				Path:      "/",
				DiskUsed:  diskUsed,
				DiskFree:  diskFree,
				DiskTotal: diskTotal,
			},
		},
	}
	//内存使用信息
	memUsed, memFree, memTotal := MemoryUsage(true)

	serverName := server.ServerName

	return &SystemInfo{
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
