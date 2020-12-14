package memory

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"testing"
	"time"
)

func TestMemory(t *testing.T) {
	cache := New(NewConfig(4, 1, 1000, 1024))
	log.GlobalLevel = log.DEBUG
	cache.SetWithRemove("hello1", "world", 1)
	cache.SetWithRemove("hello2", "world", 1)
	cache.SetWithRemove("hello3", "world", 3)
	time.Sleep(time.Second * 60)
	fmt.Println(cache.Get("hello1"))
	fmt.Println(cache.Get("hello2"))
	fmt.Println(cache.Get("hello3"))
}
