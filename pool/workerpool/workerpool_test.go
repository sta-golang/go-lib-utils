package workerpool

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/log"
	"testing"
	"time"
)

func TestWorkerPool_ReadyQueueLength(t *testing.T) {
	pool := New(2)
	_ = pool.Submit(func() {
		fmt.Println(1)
		time.Sleep(time.Millisecond * 50)
	})
	_ = pool.Submit(func() {
		fmt.Println(2)
		time.Sleep(time.Millisecond * 50)
	})
	_ = pool.Submit(func() {
		fmt.Println(3)
		time.Sleep(time.Millisecond * 50)
	})
	_ = pool.Submit(func() {
		fmt.Println(4)
		time.Sleep(time.Millisecond * 50)
	})
	_ = pool.Submit(func() {
		fmt.Println(5)

		time.Sleep(time.Millisecond * 50)
	})
	for pool.ReadyQueueLength() != 0 {
		log.Info(pool.ReadyQueueLength())
		//log.Infof("%d %d", pool.currentRunningWorker, pool.currentWorkers)
		//time.Sleep(time.Millisecond * 50)
	}

}

func TestWorkerPool_StopGetTasks(t *testing.T) {
	pool := New(2)
	_ = pool.Submit(func() {
		fmt.Println(1)
		time.Sleep(time.Second * 1)
		fmt.Println(1, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(2)
		time.Sleep(time.Second * 1)
		fmt.Println(2, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(3)
		time.Sleep(time.Second * 1)
		fmt.Println(3, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(4)
		time.Sleep(time.Second * 1)
		fmt.Println(4, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(5)
		time.Sleep(time.Second * 1)
		fmt.Println(5, " end")
	})
	time.Sleep(time.Millisecond * 50)
	//fmt.Println(pool.StopGetTasks())
	tks := pool.StopGetTasks()
	fmt.Println(tks)
	err := pool.Submit(func() {
		fmt.Println("hello")
	})
	fmt.Println(err)
	for i := range tks {
		tks[i]()
	}
}

func TestWorkerPool_StopWait(t *testing.T) {
	pool := New(2)
	_ = pool.Submit(func() {
		fmt.Println(1)
		time.Sleep(time.Second * 1)
		fmt.Println(1, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(2)
		time.Sleep(time.Second * 1)
		fmt.Println(2, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(3)
		time.Sleep(time.Second * 1)
		fmt.Println(3, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(4)
		time.Sleep(time.Second * 1)
		fmt.Println(4, " end")
	})
	_ = pool.Submit(func() {
		fmt.Println(5)
		time.Sleep(time.Second * 1)
		fmt.Println(5, " end")
	})
	time.Sleep(time.Millisecond * 50)
	//fmt.Println(pool.StopGetTasks())
	pool.StopWait()

}
