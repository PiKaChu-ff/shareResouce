package utils

import "time"

// 定义时间格式
const TimeFormat string = "2006-01-02 15:04:05.000"

// 获取时间
func GetTime() string {
	return time.Now().Format(TimeFormat)
}
