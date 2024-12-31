package rand

import (
	"bytes"
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

var (
	once sync.Once

	// SeededSecurely is set to true if a cryptographically secure seed
	// was used to initialized rand. When false, the start time is used
	// as a seed.
	SeededSecurely bool
)

// SeedMathRand provides weak, but guaranteed seeding, which is better than
// running with Go's default seed of 1. A call to SeedMathRand() is expected
// to be called via init(), but never a second time.
func SeedMathRand() {
	once.Do(func() {
		n, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			rand.Seed(time.Now().UTC().UnixNano())
			return
		}
		rand.Seed(n.Int64())
		SeededSecurely = true
	})
}

// 所有可用的十六进制字符，小写
var HexStrings []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

// 所有可用的十进制字符
var DecimalStrings []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 随机产生一个字符串
func RandStrings(length int) string {
	var buffer bytes.Buffer
	if !SeededSecurely {
		SeedMathRand()
	}
	for i := 0; i < length; i++ {
		buffer.WriteString(HexStrings[rand.Intn(len(HexStrings))])
	}
	return buffer.String()
}

// 随机产生一个字符串，字符来源于inString参数
func RandStringsIn(length int, inString []string) string {
	if len(inString) <= 0 {
		return ""
	}
	if !SeededSecurely {
		SeedMathRand()
	}
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteString(inString[rand.Intn(len(inString))])
	}
	return buffer.String()
}

// 随机产生一个指定长度的byte数组
func RandByteArray(length int) []int8 {
	if !SeededSecurely {
		SeedMathRand()
	}
	datas := make([]int8, length)
	for i := 0; i < length; i++ {
		//随机生成-128到127之间的值
		datas[i] = int8(rand.Intn(127) - 128)
	}
	return datas
}

// 在2个值之间随机一个值, 上限不能为负数，异常时返回-1
func RandInt32(minValue int, maxValue int) int {
	if !SeededSecurely {
		SeedMathRand()
	}
	if minValue > maxValue {
		return -1
	} else if minValue == maxValue {
		return maxValue
	} else if maxValue <= 0 {
		return -1
	}
	return rand.Intn(maxValue-minValue) + minValue
}
