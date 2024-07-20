package str

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

var Placeholder struct{}

// 返回一组int32数组的字符串描述
func Int32ArrayToString(source []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(source), " ", delim, -1), "[]")
}

// 返回一组int8数组的字符串描述
func Int8ArrayToString(source []int8, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(source), " ", delim, -1), "[]")
}

// 确保字符串以指定的字符串结尾，如果不以原有的字符串结尾，则自动加上
func EnsureEndWith(s string, suffix string) string {
	sourceString := string(s)
	if suffix == "" {
		//为空，则直接返回本身
		return s
	}
	if strings.HasSuffix(sourceString, suffix) {
		//已经以什么结束
		return s
	}
	return sourceString + suffix
}

// 确保字符串以指定的字符串开始，如果不以原有的字符串开始，则自动加上
func EnsureStartWith(s string, prefix string) string {
	if prefix == "" {
		//为空，则直接返回本身
		return s
	}
	if strings.HasPrefix(s, prefix) {
		//已经以什么开始
		return s
	}
	return prefix + s
}

func Hash32(v string) uint32 {
	h := fnv.New32()
	h.Write([]byte(v))
	return h.Sum32()
}

func Hash64(v string) uint64 {
	h := fnv.New64()
	h.Write([]byte(v))
	return h.Sum64()
}

// 合并2个数组
func Union(first, second []string) []string {
	set := make(map[string]struct{})

	for _, each := range first {
		set[each] = Placeholder
	}
	for _, each := range second {
		set[each] = Placeholder
	}

	merged := make([]string, 0, len(set))
	for k := range set {
		merged = append(merged, k)
	}

	return merged
}

func JoinInt64(elems []int64, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(elems[0], 10)
	}

	var b strings.Builder
	b.WriteString(strconv.FormatInt(elems[0], 10))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.FormatInt(i, 10))
	}
	return b.String()
}
