package system_info

import (
	"fmt"
	"strings"
)

const (
	mb = 1024 * 1024
)

type DiskInfo struct {
	DiskUsed  uint64         `json:"disk_used"`  //已用磁盘空间,字节数
	DiskFree  uint64         `json:"disk_free"`  //可用磁盘空间,字节数
	DiskTotal uint64         `json:"disk_total"` //总磁盘空间,字节数
	Children  []ChildrenInfo `json:"children"`
}

type ChildrenInfo struct {
	Path      string `json:"path"`       //路径
	DiskUsed  uint64 `json:"disk_used"`  //已用磁盘空间,字节数
	DiskFree  uint64 `json:"disk_free"`  //可用磁盘空间,字节数
	DiskTotal uint64 `json:"disk_total"` //总磁盘空间,字节数
}

// SystemInfo 系统信息
type SystemInfo struct {
	ServerName   string   `json:"server_name"`    //服务器名称
	SystemOs     string   `json:"system_os"`      //操作系统名称
	Runtime      int64    `json:"run_time"`       //服务运行时间,纳秒
	GoroutineNum int      `json:"goroutine_num"`  //goroutine数量
	CPUNum       int      `json:"cpu_num"`        //cpu核数
	CPUUser      float64  `json:"cpu_user"`       //cpu用户态比率
	CPUFree      float64  `json:"cpu_free"`       //cpu空闲比率
	DiskInfo     DiskInfo `json:"disk_info"`      //磁盘信息
	MemUsed      uint64   `json:"mem_used"`       //已用内存,字节数
	MemSys       uint64   `json:"mem_sys"`        //系统内存占用量,字节数
	MemFree      uint64   `json:"mem_free"`       //剩余内存,字节数
	MemTotal     uint64   `json:"mem_total"`      //总内存,字节数
	AllocGolang  uint64   `json:"alloc_golang"`   //golang内存使用量,字节数
	AllocTotal   uint64   `json:"alloc_total"`    //总分配的内存,字节数
	Lookups      uint64   `json:"lookups"`        //指针查找次数
	Mallocs      uint64   `json:"mallocs"`        //内存分配次数
	Frees        uint64   `json:"frees"`          //内存释放次数
	LastGCTime   uint64   `json:"last_gc_time"`   //上次GC时间,纳秒
	NextGC       uint64   `json:"next_gc"`        //下次GC内存回收量,字节数
	PauseTotalNs uint64   `json:"pause_total_ns"` //GC暂停时间总量,纳秒
	PauseNs      uint64   `json:"pause_ns"`       //上次GC暂停时间,纳秒
}

func (si *SystemInfo) String() string {
	buff := strings.Builder{}
	buff.WriteString("* Server Info : \n")
	buff.WriteString(fmt.Sprintf("\tServerName : %s | SystemOs : %s | 运行时间 : %v | goroutine数量 : %v\n", si.ServerName, si.SystemOs, si.Runtime, si.GoroutineNum))
	buff.WriteString("* CPU Info : \n")
	buff.WriteString(fmt.Sprintf("\tcpu核数 : %v | cpu用户态比率 : %.2f%% | cpu空闲比率 : %.2f%%\n", si.CPUNum, si.CPUUser*100, si.CPUFree*100))
	buff.WriteString("* Disk Info : \n")
	buff.WriteString(fmt.Sprintf("\t已用磁盘空间 : %v(MB) | 可用磁盘空间 : %v(MB) | 总磁盘空间 : %v(MB)\n", si.DiskInfo.DiskUsed/mb, si.DiskInfo.DiskFree/mb, si.DiskInfo.DiskTotal/mb))
	for _, c := range si.DiskInfo.Children {
		buff.WriteString(fmt.Sprintf("\t ====> 路径: %s | 已用磁盘空间 : %v(MB) | 可用磁盘空间 : %v(MB) | 总磁盘空间 : %v(MB)\n", c.Path, c.DiskUsed/mb, c.DiskFree/mb, c.DiskTotal/mb))
	}
	buff.WriteString("* memory Info : \n")
	buff.WriteString(fmt.Sprintf("\t已用内存 : %v(MB) | 系统内存占用量 : %v(MB) | 剩余内存 : %v(MB)\n", si.MemUsed/mb, si.MemSys/mb, si.MemFree/mb))
	buff.WriteString(fmt.Sprintf("\tgolang内存使用量(MB) : %v | 总分配的内存 %v (MB) | 总内存 : %v(MB)\n", si.AllocGolang/mb, si.AllocTotal/mb, si.MemTotal/mb))
	buff.WriteString(fmt.Sprintf("\t指针查找次数 : %v | 内存分配次数 : %v | 内存释放次数 : %v\n", si.Lookups, si.Mallocs, si.Frees))
	buff.WriteString("* gc Info : \n")
	buff.WriteString(fmt.Sprintf("\t上次GC时间 : %v(纳秒) | 下次GC内存回收量 : %v(字节)\n", si.LastGCTime, si.NextGC/mb))
	buff.WriteString(fmt.Sprintf("\tGC暂停时间总量 : %v() | 上次GC暂停时间 : %v()", si.PauseTotalNs, si.PauseNs))
	return buff.String()
}
