# async 并发工具

## dag 有向无环图 ./dag

**如果你想让你的任务的执行速度提高数倍而你的任务又很复杂也许dag会帮到你**

比如说你的某个任务 有a  c 和 b是可以并行的。 但是c的执行又需要依赖于a（比如a是初始化 c需要用到a的数据）

那么就可以称a是c的子任务
![image](https://user-images.githubusercontent.com/45283622/115694704-b4933b80-a393-11eb-9fdc-bac1fce9ac0d.png)


这个时候你就可以两路并行。而且不需要特别麻烦的编码就可以完成
如果你这个任务可能是无数个小任务呢。

那你的任务就是来构建这个图就可以了 并不需要管他的执行是如何调度的！

**test** 也有代码的示例

## group 任务组 ./group

### 并发任务组 async_group.go

平常在写多个并发的时候可能需要 sync.Waitgroup

写起代码会非常难看。而且还需要很多的判断 比如说 任务的返回值 和 返回的错误平时写起来 收集这些东西是很麻烦的

不要担心 async_group 帮你轻松搞定 你只需要看example 就可以轻松完成

### 单程执行器 single_exec_group.go

平时在开发的时候难免会遇到访问数据库，或者调用接口获取想要的数据的情况

但是仔细看看 发现这些返回的数据是同一个。

举个栗子：比如一个商城中的某个商品是个热点数据，那么一分钟内可能有1万个用户去访问

那这1分钟光这个数据的访问可能就是1万个socket并且他们返回的数据都是一样的。
![image](https://user-images.githubusercontent.com/45283622/115695071-0a67e380-a394-11eb-9542-e99015f65f34.png)


数据库可能会进行1万次磁盘io。这个时候不光socket链接耗费的资源，磁盘io耗费的资源浪费。

并且这么大的socket很可能击穿你的数据库。

这个时候有人提到了使用缓存，缓存的手段确实可以解决大量访问数据库的问题。

但是还有可能 同一时间 有1000个用户同时访问 这个时候你的缓存是空的，那么还会进入到数据库。

如果你的数据库单核 或者 双核 的情况下 同时1000个访问 你的数据库可能就直接扛不住了。500可能都扛不住。

这种情况还是有问题的。

所以最好的解决方案是 single_exec_group + 缓存。 single_exec_group 只允许统一时间对于同样的任务只有一个线程去读取数据。
这个时候读取完后再加入到缓存中。

上述的问题 同一秒这1000个线程进来了 发现没有缓存，然后去读数据库。

这1000个线程 只有第一个到达的线程去读取数据库了。 其余999个线程被睡眠

第一个线程读取到的数据 共享给这999个线程 这999个线程拿到数据返回。

之后的线程再访问的时候 就可以访问到缓存中的内容了。 这种情况 只会访问1次数据库
![image](https://user-images.githubusercontent.com/45283622/115695382-5581f680-a394-11eb-8e13-fcacf3e53507.png)


**笔者在测试自己的单核云数据库的时候 500的并发量 使用single_exec_group比不使用的性能提高了160倍**

**并且随着并发量的逐渐提升这个的性能还要提高上限**

示例代码如下

```go
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
```

## 资源处理 ./process

如果你有一些并发的资源，如果遇到停止信号或者错误 什么的 需要同步这些资源
那么就可以使用这个 请查看example
