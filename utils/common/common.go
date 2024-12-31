package common

import (
	"runtime"
	"strings"
)

// 当前调用函数名
func CurFuncName() string {
	return GetCallerName(2)
}

// 获取第几层函数的名称
// 0-当前层，1-上一层，2-再上一层，依次向上
func GetCallerName(l int) string {
	pc, _, _, _ := runtime.Caller(l)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, ".")
	//fmt.Printf("第%d层函数,函数名称是:%s\n", l, name)
	return split[len(split)-1]
}
