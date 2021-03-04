package utils

import "time"

var (
	// 通用的日期、时间格式模板
	// 因为自带的模板不太好记
	DatetimeFormat   = "2006-01-02 15:04:05"
	DateFormat       = "2006-01-02"
	TimeFormat       = "15:04:05"
	DateMinuteFormat = "2006-01-02 15:04"
)

// 获取毫秒时间戳
func UnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}
