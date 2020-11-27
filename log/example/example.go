package example

import (
	"github.com/xy63237777/go-lib-utils/log"
	"github.com/xy63237777/go-lib-utils/source"
	"time"
)

func LogAsyncExample() {
	//----------------------------
	//这个代码最好放到main.go里
	defer func() {
		if er := recover(); er != nil {
			source.Sync()
			panic(er)
		}
	}()
	//---------------------------
	logger := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.ERROR)}), time.Second*3)
	log.SetGlobalLogger(logger)
	//后面就可以调用log

}
