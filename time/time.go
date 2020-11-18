package time

import (
	"errors"
	"time"
)

const (
	layoutDate = "2006-01-02"
	layoutTime = "15:04:05"
	layoutDateTime = "2006-01-02 15:04:05"
)
var cstZone *time.Location

// GetNowDateTimeStr 获取当前时间日期的字符串
func GetNowDateTimeStr() string{
	return parseTimeToStr(time.Now().In(cstZone), layoutDateTime)
}

// GetNowDateStr 获取当前日期的字符串
func GetNowDateStr() string{
	return parseTimeToStr(time.Now().In(cstZone), layoutDate)
}

// GetNowTimeStr 获取当前时间的字符串
func GetNowTimeStr() string{
	return parseTimeToStr(time.Now().In(cstZone), layoutTime)
}

// ParseDateTime 解析时间
func ParseDateTime(dateTimeStr string) (time.Time, error){
	return parseStrToTime(dateTimeStr, layoutDateTime)
}

// ParseDate ...
func ParseDate(dateStr string) (time.Time, error) {
	return parseStrToTime(dateStr, layoutDate)
}

// ParseTime ...
func ParseTime(timeStr string) (time.Time, error) {
	return parseStrToTime(timeStr, layoutTime)
}

func parseStrToTime(str, layout string) (time.Time, error) {
	if str == "" {
		return time.Now(), errors.New("args It can't be empty")
	}
	return time.ParseInLocation(layout,str, cstZone)
}

func parseTimeToStr(tm time.Time, layout string) string {
	return tm.Format(layout)
}
