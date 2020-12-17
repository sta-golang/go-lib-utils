package time

import (
	"errors"
	"time"
)

const (
	LayoutDate     = "2006-01-02"
	LayoutTime     = "15:04:05"
	LayoutDateTime = "2006-01-02 15:04:05"
)

func init() {
	cstZone = time.FixedZone("CST", 8*3600)
}

var cstZone *time.Location

func Location() *time.Location {
	return cstZone
}

// GetNowDateTimeStr 获取当前时间日期的字符串
func GetNowDateTimeStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutDateTime)
}

// GetNowDateStr 获取当前日期的字符串
func GetNowDateStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutDate)
}

// GetNowTimeStr 获取当前时间的字符串
func GetNowTimeStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutTime)
}

// ParseDateTime 解析时间
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return parseStrToTime(dateTimeStr, LayoutDateTime)
}

// ParseDate ...
func ParseDate(dateStr string) (time.Time, error) {
	return parseStrToTime(dateStr, LayoutDate)
}

// ParseTime ...
func ParseTime(timeStr string) (time.Time, error) {
	return parseStrToTime(timeStr, LayoutTime)
}

func parseStrToTime(str, layout string) (time.Time, error) {
	if str == "" {
		return time.Now(), errors.New("args It can't be empty")
	}
	return time.ParseInLocation(layout, str, cstZone)
}

func parseTimeToStr(tm time.Time, layout string) string {
	return tm.Format(layout)
}

func GetNowTime() time.Time {
	return time.Now().In(cstZone)
}
