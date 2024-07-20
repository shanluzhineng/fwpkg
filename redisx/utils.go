package redisx

import "strings"

// 确保字符串以指定的字符串开始，如果不以原有的字符串开始，则自动加上
func ensureStartWith(s string, prefix string) string {
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
