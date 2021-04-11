package example

import (
	"time"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/source"
	"github.com/sta-golang/go-lib-utils/str"
)

func LogAsyncExample() {
	//----------------------------
	//这个代码最好放到调用的位置
	defer func() {
		source.Sync()
		if er := recover(); er != nil {

			panic(er)
		}
	}()
	//---------------------------
	logger := log.NewFileLogAndAsync(log.DefaultFileLogConfigForAloneWriter(
		[]string{log.GetLevelName(log.INFO), log.GetLevelName(log.ERROR)}), time.Second*3)
	log.SetGlobalLogger(logger)
	//后面就可以调用log

}

// LogYamlExample 通过配置文件的方式来生成Log
// 下面使用yaml 使用json同理
func LogYamlExample() {
	yamlStr := `
log:
  file_dir: ./log
  file_name: sta
  day_age: 15
  log_level: 1 # 1为 INFO 0为DEBUG
  pre_fix: "[STA:Music-Data]"
  alone_writer:
    - DEBUG
    - INFO
    - ERROR
    - FATAL
  re_open: 180
`
	type Config struct {
		LogCfg *log.FileLogConfig `yaml:"log"`
	}
	var cfg Config
	err := codec.API.YamlAPI.UnMarshal(str.StringToBytes(&yamlStr), &cfg)
	log.ConsoleLogger.Debug(err)
	log.ConsoleLogger.Info(cfg.LogCfg)
	logger := log.NewFileLogAndAsync(cfg.LogCfg, time.Second*3)
	log.SetGlobalLogger(logger)
	log.Info("hello")
}

func LogContext() {
	// 在controller层生成ctx 使用一些手段一次请求一个ctx
	ctx := log.LogContextKeyMap(nil, map[string]string{
		"server":   "sta-golang",
		"user":     "TheSevenSky",
		"uuid":     "生成请求的唯一uuid",
		"username": "63237777@qq.com",
	})
	log.InfoContext(ctx, "info")
}
