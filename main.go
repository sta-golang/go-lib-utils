package main

import (
	"fmt"
	"github.com/xy63237777/go-lib-utils/log"
	"github.com/xy63237777/go-lib-utils/os/os_windows"
	"github.com/xy63237777/go-lib-utils/source"
	"time"
)

func main() {

	defer func() {
		if er := recover(); er != nil {
			source.Sync()
			panic(er)
		}
	}()
	log.Info(os_windows.GetWindowsSystemInfo())
	log.Error("hello")
	c := cat{
		name: "12312",
		age:  0,
	}
	hello(&c)
	time.Sleep(time.Second * 120)
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
