## log组件

### 快速开始

<hr/>
你只需要导入这个依赖

`"github.com/xy63237777/go-lib-utils/log"`

然后调用

`log.Info("hello")`

`log.Infof("test %s", "sta")`

当然这个的打印则是在控制台打印

#### 在文件打印

`fileLog := log.NewFileLog(log.DefaultFileLogConfig())`

然后调用fileLog.Info fileLog.Error等等就可以

当然如果你觉得这样太麻烦

`log.SetGlobalLogger(fileLog)`

这样像之前
`log.Info("hello")
log.Infof("test %s", "sta")`
调用就可以

如果你需要将ERROR INFO等等分开 那么也非常简单

`fileLog := log.NewFileLog(log.DefaultFileLogConfigForAloneWriter(
 		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.ERROR)}))
 	`
 	
 代码虽然很长 但是你仔细看看就很简单
 参数就只是一个字符串数组
 
 `LEVEL_FLAGS = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}`
 
 它其实是这样的
 也就是上述代码传入的数组等价于
 `fileLog := log.NewFileLog(log.DefaultFileLogConfigForAloneWriter([]string{"INFO","ERROR"}))`
 
 这样所有的FATAL级别和ERROR级别都会在ERROR文件中 所有的WARN和INFO都会在INFO文件中。 如果你设置的级别变成了DEBUG 则所有的DEBUG都在ALL文件中。这个可以读者慢慢试试。
 还有一些高级功能
 
 #### 高级功能
 
 
 ```go
 type FileLogConfig struct {
  	FileDir     string   `yaml:"file_dir"`     // 文件目录 默认为./log
  	FileName    string   `yaml:"file_name"`    // 文件前缀名 默认为sta
  	DayAge      int      `yaml:"day_age"`      // 文件保留日期 默认为7天
  	LogLevel    Level    `yaml:"log_level"`    // 日志等级 默认为INFO 这个可以之后设置
  	Prefix      string   `yaml:"pre_fix"`      // 日志输出前缀 默认为 FOUR-SEASONS: STA
  	MaxSize     int64    `yaml:"max_size"`     // 最大大小 设置0或者辅助 默认失效。 最小为10mb 如果小于10mb则变成16mb
  	AloneWriter []string `yaml:"alone_writer"` // 单独数组的等级，设置后没有出现的向等级低的方向靠
  }
  ```
  
  
  可以通过这个配置类来创建file_log
  
### 输出格式

prefix yyyy-mm-dd hh:MM:ss [level] 代码调用的文件:代码调用的行数 ==> msg

例如

`[FOUR-SEASONS: STA] 2020-11-26 15:19:47 [INFO] /home/thesevensky/gocode/sta-test/main.go:46 ==> test sta`


#### 最佳实践

例如下面代码在任何可能发生err的地方调用打印ERROR
由于此组件有打印代码行号的功能所以可以迅速定位bug
```go
    err := func() 
    if err != nil {
		log.Errorf("xxx Err %v", err)
		return err
	}
```

#### 高级用法

此方法可以极高的提高你的性能。但是它容易丢失一段时间的DEBUG-WARN级别的日志
当然你的单独写入要是有WARN的话就会丢失到INFO级别

具体的请看example

```go
defer func() {
		if er := recover(); er != nil {
			source.Sync()
			panic(er)
		}
	}()
	logger := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.ERROR)}), time.Second*3)
	log.SetGlobalLogger(logger)
```
