package main

import (
	"github.com/xy63237777/go-lib-utils/log"
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
	lg := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.ERROR)}), time.Second*3)
	lg.Infof("hello")
	go func() {
		for {
			time.Sleep(time.Second * 5)
			lg.Infof("hello")
		}
	}()

	time.Sleep(time.Second * 120)
}
