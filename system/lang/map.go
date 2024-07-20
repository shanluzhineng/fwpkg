package lang

// 从一个map中提取出key列表
func ExtractMapKeys[T comparable](m map[T]interface{}) []T {
	result := make([]T, 0)
	for eachKey := range m {
		result = append(result, eachKey)
	}
	return result
}
