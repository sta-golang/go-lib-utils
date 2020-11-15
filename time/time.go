package time

import "time"

const (
	layoutDate = "2006-01-02"
	layoutTime = "15:04:05"
	layoutDateTime = "2006-01-02 15:04:05"
)

var cstZone = time.FixedZone("CST", 8*3600)

func GetNowDateTimeStr() string{
	return parseTimeToStr(time.Now().In(cstZone), layoutDateTime)
}

func GetNowDateStr() string{
	return parseTimeToStr(time.Now().In(cstZone), layoutDate)
}

func GetNowTimeStr() string{
	return parseTimeToStr(time.Now().In(cstZone), layoutTime)
}

func parseTimeToStr(tm time.Time, layout string) string {
	return tm.Format(layout)
}
