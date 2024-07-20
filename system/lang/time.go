package lang

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func NowToPtr() *time.Time {
	return TimeToPtr(time.Now())
}

// 获取时间的指针
// 这里为啥需要这个函数，因为好多情况下需要使用时间指针，而又无法通过&time.Now()来获取
func TimeToPtr(t time.Time) *time.Time {
	return &t
}

// 将一个时间字符串解析成一个中国时区的时间
func ParseTimeToChinaTimezone(layoutList []string, timeString string) (*time.Time, error) {
	if len(layoutList) <= 0 {
		return nil, errors.New("layoutList不能为空,或者长度不能小于0")
	}
	for _, eachLayout := range layoutList {
		timeValue, err := time.ParseInLocation(eachLayout, timeString, ChinaTimezone)
		if err == nil {
			return &timeValue, nil
		}
	}
	return nil, fmt.Errorf("无效的时间,时间格式必须是%s", strings.Join(layoutList, ","))
}

func DurationHumanText(duration time.Duration) string {
	if duration <= 0 {
		return "0 seconds"
	}
	remainDuration := duration
	days := duration / (24 * time.Hour)
	remainDuration = remainDuration - days*(24*time.Hour)

	hours := remainDuration / time.Hour
	remainDuration = remainDuration - hours*time.Hour

	minutes := remainDuration / time.Minute
	remainDuration = remainDuration - minutes*time.Minute

	seconds := remainDuration / time.Second
	// remainDuration = remainDuration - seconds*time.Second
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		days,
		hours,
		minutes,
		seconds)
}
