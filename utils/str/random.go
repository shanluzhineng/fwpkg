package str

import (
	"math/rand"
	"time"
)

var HexStrings []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

// 所有可用的十进制字符
var DecimalStrings []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// random stings
func RandStrings(length int) string {
	return RandStringsIn(length, charset)
}

// 随机产生一个字符串，字符来源于inString参数
func RandStringsIn(length int, inString string) string {
	if len(inString) <= 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

// 随机产生一个字符串，字符来源于inString参数
func RandStringsInArray(length int, inString []string) []string {
	if len(inString) <= 0 {
		return []string{}
	}
	rand.Seed(time.Now().UnixNano())

	result := make([]string, length)
	for i := range result {
		result[i] = inString[rand.Intn(len(charset))]
	}

	return result
}
