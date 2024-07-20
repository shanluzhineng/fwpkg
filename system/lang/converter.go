package lang

// 将bool值转换为指针
func BoolToPtr(v bool) *bool {
	return &v
}

// 将bool指针转换为bool,如果指针为nil，则直接返回false
func PtrToBool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// 将int值转换为指针
func IntToPtr(v int) *int {
	return &v
}

// 将int指针转换为int,如果指针为nil，则直接返回0
func PtrToInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func SliceToInterfaceSlice[T any](list []T) []interface{} {
	result := make([]interface{}, 0)
	if len(list) <= 0 {
		return result
	}
	for index := range list {
		result = append(result, list[index])
	}
	return result
}

func MapToInterfaceMap[T comparable, V any](v map[T]V) map[T]interface{} {
	result := make(map[T]interface{}, 0)
	if len(v) <= 0 {
		return result
	}
	for eachKey, eachValue := range v {
		result[eachKey] = eachValue
	}
	return result
}
