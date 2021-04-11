package main

import (
	"fmt"
	"time"

	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/log/example"
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
	//TestLog()
	//example.LogYamlExample()
	example.LogContext()
}

func TestLog() {
	lg := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.WARNING), log.GetLevelName(log.ERROR)}), time.Second*3)
	lg.Warn("hello")
	ctx := log.LogContextKeyMap(nil, map[string]string{
		"user": "thesevensky",
		"id":   "63237777",
	})
	log.ConsoleLogger.Info("hello")
	log.ConsoleLogger.InfoContext(ctx, "hello")
	log.ConsoleLogger.Warn("hello")
	log.ConsoleLogger.Error("hello")
	log.ConsoleLogger.InfoContextf(ctx, "hello %s = %v", "hello1", 1)
	log.ConsoleLogger.InfoContextf(ctx, "hello wo %s %d", " test", 1)
	log.ConsoleLogger.Info("hello")
	log.ConsoleLogger.Infof("hello %s", "test")
	lg.Warnf("hello world")
	lg.WarnContext(ctx, "hello world11111111111111111111111111111")
	lg.WarnContextf(ctx, "hello %s lalal", " test")
	lg.WarnContextf(ctx, "hello world %s", "11111")
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
