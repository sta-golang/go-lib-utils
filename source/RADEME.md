## 资源管理组件


#### 作用

主要是用来在程序退出的时候同步异步资源

推荐在main.go加上

```go
defer func() {
		if er := recover(); er != nil {
			source.Sync()
			panic(er)
		}
	}()
```