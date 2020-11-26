package main

import (
	"encoding/hex"
	"fmt"
	"github.com/xy63237777/go-lib-utils/algorithm/data_structure"
	"github.com/xy63237777/go-lib-utils/codec"
	"github.com/xy63237777/go-lib-utils/log"
	ow "github.com/xy63237777/go-lib-utils/os/os_windows"
)

func main() {
	m := map[string]int{
		"123": 1, "345": 2,
	}
	bytes, e := codec.API.JsonAPI.Marshal(m)
	fmt.Println(string(bytes), e)
	list := data_structure.NewLinkedList()
	list.Add("hello")
	list.Add("world")
	list.Add("golang")
	fmt.Println((4 << 2) - (4 >> 1))
	pwd := "qq123456"
	//sum := md5.Sum([]byte(pwd))
	str := hex.EncodeToString([]byte(pwd))
	fmt.Println(str)
	log.Error("hello")

	strs := make([]string, 0, 2)
	strs = append(strs, "hello")
	strs = append(strs, "hello")
	strs = append(strs, "hello")
	strs = append(strs, "hello")
	strs = append(strs, "hello")
	fmt.Println(strs)
	fmt.Println(ow.GetWindowsSystemInfo())

	fileLog := log.NewFileLog(log.DefaultFileLogConfigForAloneWriter([]string{
		log.LEVEL_FLAGS[log.WARNING], log.LEVEL_FLAGS[log.ERROR]}))

	fileLog.Info("hello")
	fileLog.Warn("warn")
	fileLog.Error("error", "hahahha")
	fileLog.Errorf("xixixi")
	log.Info("hello")
	log.Infof("test %s", "sta")
}
