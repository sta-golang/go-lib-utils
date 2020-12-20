## workerpool

### 快速开始
```go
       num := runtime.NumCPU()
	pool := New(num)
        err := pool.Submit(func() {
		fmt.Println(1)
		time.Sleep(time.Second * 1)
		fmt.Println(1, " end")
	})
	if err != nil {
		//一般可能是因为队列满了或者使用了已经停止的协程池
	}
```