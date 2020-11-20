package os_linux

import (
	"fmt"
	tm "github.com/xy63237777/go-lib-utils/time"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

const (
	mb = 1024 * 1024
)

// SystemInfo 系统信息
type SystemInfo struct {
	ServerName   string  `json:"server_name"`    //服务器名称
	SystemOs     string  `json:"system_os"`      //操作系统名称
	Runtime      int64   `json:"run_time"`       //服务运行时间,纳秒
	GoroutineNum int     `json:"goroutine_num"`  //goroutine数量
	CPUNum       int     `json:"cpu_num"`        //cpu核数
	CPUUser      float64 `json:"cpu_user"`       //cpu用户态比率
	CPUFree      float64 `json:"cpu_free"`       //cpu空闲比率
	DiskUsed     uint64  `json:"disk_used"`      //已用磁盘空间,字节数
	DiskFree     uint64  `json:"disk_free"`      //可用磁盘空间,字节数
	DiskTotal    uint64  `json:"disk_total"`     //总磁盘空间,字节数
	MemUsed      uint64  `json:"mem_used"`       //已用内存,字节数
	MemSys       uint64  `json:"mem_sys"`        //系统内存占用量,字节数
	MemFree      uint64  `json:"mem_free"`       //剩余内存,字节数
	MemTotal     uint64  `json:"mem_total"`      //总内存,字节数
	AllocGolang  uint64  `json:"alloc_golang"`   //golang内存使用量,字节数
	AllocTotal   uint64  `json:"alloc_total"`    //总分配的内存,字节数
	Lookups      uint64  `json:"lookups"`        //指针查找次数
	Mallocs      uint64  `json:"mallocs"`        //内存分配次数
	Frees        uint64  `json:"frees"`          //内存释放次数
	LastGCTime   uint64  `json:"last_gc_time"`   //上次GC时间,纳秒
	NextGC       uint64  `json:"next_gc"`        //下次GC内存回收量,字节数
	PauseTotalNs uint64  `json:"pause_total_ns"` //GC暂停时间总量,纳秒
	PauseNs      uint64  `json:"pause_ns"`       //上次GC暂停时间,纳秒
}

func (si *SystemInfo) String() string {
	buff := strings.Builder{}
	buff.WriteString("* Server Info : \n")
	buff.WriteString(fmt.Sprintf("\tServerName : %s | SystemOs : %s | 运行时间 : %v | goroutine数量 : %v\n", si.ServerName, si.SystemOs, si.Runtime, si.GoroutineNum))
	buff.WriteString("* CPU Info : \n")
	buff.WriteString(fmt.Sprintf("\tcpu核数 : %v | cpu用户态比率 : %.2f%% | cpu空闲比率 : %.2f%%\n", si.CPUNum, si.CPUUser*100, si.CPUFree*100))
	buff.WriteString("* Disk Info : \n")
	buff.WriteString(fmt.Sprintf("\t已用磁盘空间 : %v(MB) | 可用磁盘空间 : %v(MB) | 总磁盘空间 : %v(MB)\n", si.DiskUsed/mb, si.DiskFree/mb, si.DiskTotal/mb))
	buff.WriteString("* memory Info : \n")
	buff.WriteString(fmt.Sprintf("\t已用内存 : %v(MB) | 系统内存占用量 : %v(MB) | 剩余内存 : %v(MB)\n", si.MemUsed/mb, si.MemSys/mb, si.MemFree/mb))
	buff.WriteString(fmt.Sprintf("\tgolang内存使用量(MB) : %v | 总分配的内存 %v (MB) | 总内存 : %v(MB)\n", si.AllocGolang/mb, si.AllocTotal/mb, si.MemTotal/mb))
	buff.WriteString(fmt.Sprintf("\t指针查找次数 : %v | 内存分配次数 : %v | 内存释放次数 : %v\n", si.Lookups, si.Mallocs, si.Frees))
	buff.WriteString("* gc Info : \n")
	buff.WriteString(fmt.Sprintf("\t上次GC时间 : %v(纳秒) | 下次GC内存回收量 : %v(字节)\n", si.LastGCTime, si.NextGC/mb))
	buff.WriteString(fmt.Sprintf("\tGC暂停时间总量 : %v() | 上次GC暂停时间 : %v()", si.PauseTotalNs, si.PauseNs))
	return buff.String()
}

// GoMemory 获取当前go程序的内存使用,返回字节数.
func GoMemory() uint64 {
	stat := new(runtime.MemStats)

	runtime.ReadMemStats(stat)
	return stat.Alloc
}

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

// Hostname 获取主机名.
func Hostname() (string, error) {
	return os.Hostname()
}

// GetSystemInfo 获取系统运行信息.
func GetSystemInfo() *SystemInfo {
	//运行时信息
	mstat := &runtime.MemStats{}
	runtime.ReadMemStats(mstat)

	//CPU信息
	cpuUser, cpuIdel, cpuTotal := CPUUsage()
	cpuUserRate := float64(cpuUser) / float64(cpuTotal)
	cpuFreeRate := float64(cpuIdel) / float64(cpuTotal)

	//磁盘空间信息
	diskUsed, diskFree, diskTotal := DiskUsage("/")

	//内存使用信息
	memUsed, memFree, memTotal := MemoryUsage(true)

	serverName, _ := os.Hostname()

	return &SystemInfo{
		ServerName: serverName,
		SystemOs:   runtime.GOOS,

		Runtime:      int64(tm.ServiceUptime()),
		GoroutineNum: runtime.NumGoroutine(),
		CPUNum:       runtime.NumCPU(),
		CPUUser:      cpuUserRate,
		CPUFree:      cpuFreeRate,
		DiskUsed:     diskUsed,
		DiskFree:     diskFree,
		DiskTotal:    diskTotal,
		MemUsed:      memUsed,
		MemSys:       mstat.Sys,
		MemFree:      memFree,
		MemTotal:     memTotal,
		AllocGolang:  mstat.Alloc,
		AllocTotal:   mstat.TotalAlloc,
		Lookups:      mstat.Lookups,
		Mallocs:      mstat.Mallocs,
		Frees:        mstat.Frees,
		LastGCTime:   mstat.LastGC,
		NextGC:       mstat.NextGC,
		PauseTotalNs: mstat.PauseTotalNs,
		PauseNs:      mstat.PauseNs[(mstat.NumGC+255)%256],
	}
}
