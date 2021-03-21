package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
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
	fmt.Println(tm.ParseDataTimeToStr(tm.GetNowTime().Add(-(time.Hour * 24 * 30))))
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
