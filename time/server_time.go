package time

import "time"

var startTime = time.Now()

// ServiceUptime 获取当前服务运行时间,纳秒int64.
func ServiceUptime() time.Duration {
	return time.Since(startTime)
}

