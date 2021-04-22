package example

import (
	"fmt"
	gp "github.com/sta-golang/go-lib-utils/async/group"
	"time"
)

func exampleAsyncGroup()  {
	group :=  gp.NewAsyncGroup(1)
	defer group.Close()
	err := group.Add("hello", func() (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println("hello 1")
		return "hello", nil
	})
	err = group.Add("hello", func() (interface{}, error) {
		time.Sleep(time.Millisecond * 5000)
		fmt.Println("xixixixi")
		return "2", nil
	})
	if err != nil {
		fmt.Println(err)
		_ = group.Add("hello2", func() (interface{}, error) {
			time.Sleep(time.Millisecond * 7000)
			fmt.Println("xixixixi")
			return "2", nil
		})
	}
	_ = group.Add("hello3", func() (interface{}, error) {
		time.Sleep(time.Millisecond * 2000)
		fmt.Println("ccccccc")
		return nil, nil
	})
	group.Wait()
	for _, tk := range group.Iterator() {
		fmt.Println(tk.Ret())
	}
	fmt.Println(group.GetTask("hello2").Ret())
	fmt.Println(group.GetTask("hello3").Ret())
}

func exampleSingleExec() {
	key := "一个任务一个key相同任务key相同"
	// cache.get 获取到了 直接返回
	type dbRet struct {
	}

	queryDB := func() (dbRet, error) {
		fmt.Println("获取db数据。。。")
		time.Sleep(time.Second) //模拟获取db的值
		return dbRet{}, nil
	}
	seg := gp.SingleExecGroup{}
	retData, err := seg.Do(key, func() (interface{}, error) {
		ret ,err := queryDB()
		if err != nil {
			// log.Error(err)
			return nil, err
		}
		// cache.set(key, ret) // 加入到缓存中
		return ret, nil
	})
	if err != nil {
		// return err
	}
	//return retData
	fmt.Println(retData)
}



