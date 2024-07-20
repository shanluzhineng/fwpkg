package lang

import (
	"fmt"
	"time"
)

var (
	ChinaTimezone *time.Location
	//默认的time layout
	DefaultTimeLayout string = "2006-01-02 15:04:05"
)

func init() {
	shanghaiTimezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		err = fmt.Errorf("为啥不存在Asia/Shanghai时区数据")
		fmt.Println(err.Error())
		panic(err)
	}
	ChinaTimezone = shanghaiTimezone
}
