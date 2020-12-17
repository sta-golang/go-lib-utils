package server

import (
	tm "github.com/sta-golang/go-lib-utils/time"
	"time"
)

func init() {
	startTime = time.Now().In(tm.Location())
}

var startTime time.Time

// ServiceUptime 获取当前服务运行时间,纳秒int64.
func ServiceUptime() time.Duration {
	return time.Since(startTime)
}
