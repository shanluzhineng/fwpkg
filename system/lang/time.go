package lang

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// 当前时间字符串 "2006-01-02 15:04:05"。
// 注：代码中未设置时区，与服务器时区有关
func TimeToStr(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}
	return t.Local().Format(DefaultTimeLayout)
}

// 时间转为date字符串，即：年-月-日
func TimeToDateStr(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}
	return t.In(ChinaTimezone).Format(DefaultDateLayout)
}

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

// n天前、后的时间，负数 before，正数 after
func TimeAddDays(n int64) time.Time {
	resTime := time.Now().Add(time.Hour * 24 * time.Duration(n))
	return resTime
}

// n小时前、后的时间，负数 before，正数 after
func TimeAddHours(n int64) time.Time {
	resTime := time.Now().Add(time.Hour * time.Duration(n))
	return resTime
}

// n分钟前、后的时间，负数 before，正数 after
func TimeAddMinutes(n int64) time.Time {
	resTime := time.Now().Add(time.Minute * time.Duration(n))
	return resTime
}

// 时间戳转时间
func TimeFromTimestampMs(milliTs int64) time.Time {
	return time.UnixMilli(milliTs)
}
