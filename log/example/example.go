package example

import (
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/source"
	"time"
)

func LogAsyncExample() {
	//----------------------------
	//这个代码最好放到main.go里
	defer func() {
		source.Sync()
		if er := recover(); er != nil {

			panic(er)
		}
	}()
	//---------------------------
	logger := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.ERROR)}), time.Second*3)
	log.SetGlobalLogger(logger)
	//后面就可以调用log

}
