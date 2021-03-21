package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	systeminfo "github.com/sta-golang/go-lib-utils/os/system_info"
	tm "github.com/sta-golang/go-lib-utils/time"
	"runtime"
	"time"
)

func TestTag() {
	i := 0

	i++
MYLoop:
	for ; i < 10; i++ {
		if i == 5 {
			break MYLoop
		}
	}
	fmt.Println(i)
}

func main() {
	arr := make([]int, 1024000)
	arr2 := make([]int, 1024000)
	arr2[3] = 5
	time.Sleep(time.Second)
	arr2 = nil
	runtime.GC()
	time.Sleep(time.Second * 2)
	fmt.Println(tm.ParseDataTimeToStr(tm.GetNowTime().Add(-(time.Hour * 24 * 30))))
	timing := tm.FuncTiming(func() {
		fmt.Println(systeminfo.MemoryUsage())
	})
	fmt.Println(timing)
	info := systeminfo.GetSystemInfo()

	fmt.Println(info)
	arr[50] = 300
}

func ReCreate() {
	lg := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.WARNING), log.GetLevelName(log.ERROR)}), time.Second*3)
	lg.Warn("hello")
}

func hello(t A) {
	fmt.Println(t.show())
}

type A interface {
	show() string
}

type cat struct {
	name string
	age  int
}

func (c *cat) show() string {
	return c.name
}
