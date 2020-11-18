package time

import "time"

func init() {
	cstZone = time.FixedZone("CST", 8*3600)
	startTime = time.Now().In(cstZone)
}
var startTime time.Time

// ServiceUptime 获取当前服务运行时间,纳秒int64.
func ServiceUptime() time.Duration {
	return time.Since(startTime)
}

