package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"sync"
	"sync/atomic"
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
	wg := sync.WaitGroup{}
	start := sync.WaitGroup{}
	ready := sync.WaitGroup{}

	var cnt int32
	cnt = 0
	start.Add(1)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		ready.Add(1)
		go func() {
			defer wg.Done()
			ready.Done()
			start.Wait()
			for j := 0; j < 10000; j++ {
				atomic.AddInt32(&cnt, 1)
			}
		}()
	}
	ready.Wait()
	start.Done()
	wg.Wait()
	fmt.Println(cnt)
	log.Infof("hello %s 123", "hello")
	log.Info("hello")
	log.Info("hello")
	log.Info("hello")
	
	//ReCreate()
	//if runtime.GOOS == "windows" {
	//	fmt.Println(system_info.GetSystemInfo())
	//}
	//time.Sleep(time.Second * 50)
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
