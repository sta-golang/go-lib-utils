package main

import (
	"context"
	"fmt"
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
			log.Infof("hello")
		}
	}()

	testCh := make(chan *string, 5)
	testCh <- nil
	data := <-testCh
	fmt.Println(data)
	ctx := context.Background()
	cancel, cancelFunc := context.WithCancel(ctx)
	go testCtx(cancel)
	go testCtx(cancel)
	go testCtx(cancel)
	go func() {
		time.Sleep(time.Second)
		cancelFunc()
	}()
	time.Sleep(time.Second * 30)
}

func testCtx(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("end")
			return
		default:
			fmt.Println("hello")
			time.Sleep(time.Millisecond * 125)
		}
	}
}
