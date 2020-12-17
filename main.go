package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/os/system_info"

	"runtime"
	"time"
)

func main() {
	ReCreate()
	if runtime.GOOS != "windows" {
		fmt.Println(system_info.GetSystemInfo())
	}
	time.Sleep(time.Second * 50)
	//var aStr *string
	//
	//target := "world"
	//atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&aStr)), unsafe.Pointer(&target))
	//fmt.Println(*aStr)
	//fmt.Println(target[:len(target)])
	//defer func() {
	//	if er := recover(); er != nil {
	//		source.Sync()
	//		panic(er)
	//	}
	//}()
	//log.Info(os_windows.GetWindowsSystemInfo())
	//log.Error("hello")
	//c := cat{
	//	name: "12312",
	//	age:  0,
	//}
	//hello(&c)
	//time.Sleep(time.Second * 120)
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
