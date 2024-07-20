package mapx

import (
	"encoding/json"
	"sort"
)

// 使用json的方式将一个对象转换为另一个对象
func ConvertObjectTo(src interface{}, dest interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}
	return nil
}

// extract map[k]v from []v
// keySelectFunc cannot be nil
func FromSlice[k comparable, v any](source []v, keySelectFunc func(v) k) map[k]v {
	value := make(map[k]v)
	if len(source) <= 0 {
		return value
	}
	for i := range source {
		key := keySelectFunc(source[i])
		value[key] = source[i]
	}
	return value
}

// extract map[k]v from []interface{}
// keyValueSelectFunc cannot be nil
// k, map[k]v key type
// v, map[k]v value type
// SType, slice element type
func FromSliceWithT[SType any, k comparable, v any](source []SType, keyValueSelectFunc func(SType) (k, v)) map[k]v {
	value := make(map[k]v)
	if len(source) <= 0 {
		return value
	}
	for i := range source {
		key, keyValue := keyValueSelectFunc(source[i])
		value[key] = keyValue
	}
	return value
}

// extract map[k]v from []interface{}
// keyValueSelectFunc cannot be nil
// k, map[k]v key type
// v, map[k]v value type
func FromSliceWith[k comparable, v any](source []interface{}, keyValueSelectFunc func(interface{}) (k, v)) map[k]v {
	value := make(map[k]v)
	if len(source) <= 0 {
		return value
	}
	for i := range source {
		key, keyValue := keyValueSelectFunc(source[i])
		value[key] = keyValue
	}
	return value
}

// extract map value to []v from map[k]v
func ToSlice[k comparable, v any](source map[k]v) []v {
	value := make([]v, 0)
	if len(source) <= 0 {
		return value
	}
	for _, eachItem := range source {
		value = append(value, eachItem)
	}
	return value
}

// extract map value to []string from map[k]v
func ToSliceStringWith[k comparable, v any](source map[k]v, valueFunc func(k) string) []string {
	value := make([]string, 0)
	if len(source) <= 0 {
		return value
	}
	for eachKey := range source {
		value = append(value, valueFunc(eachKey))
	}
	return value
}

// extract []k from map[k]v
func ExtractMapKeys[k comparable, v any](m map[k]v) []k {
	result := make([]k, 0)
	for eachKey := range m {
		result = append(result, eachKey)
	}
	return result
}

// extract []v from map[k]v
func ExtractMapValues[k comparable, v any](m map[k]v) []v {
	result := make([]v, 0)
	for _, eachValue := range m {
		result = append(result, eachValue)
	}
	return result
}

// extract []string from map[string]v
func ExtractMapKeysAsString[v any](m map[string]v) []string {
	result := make([]string, 0)
	for eachKey := range m {
		result = append(result, eachKey)
	}
	return result
}

// extract []string from map[string]v
func ExtractMapSortedKeysAsString[v any](m map[string]v) []string {
	result := make([]string, 0)
	for eachKey := range m {
		result = append(result, eachKey)
	}
	sortedSlice := sort.StringSlice(result)
	sortedSlice.Sort()
	return sortedSlice
}

// extract sorted []k from map[k]v
// asc sorted
func ExtractMapSortedKeys[k comparable, v any](m map[k]v, less func(i, j int) bool) []k {
	result := make([]k, 0)
	for eachKey := range m {
		result = append(result, eachKey)
	}
	sort.Slice(result, less)
	return result
}

func ToInterfaceMap[T comparable, V any](v map[T]V) map[T]interface{} {
	result := make(map[T]interface{}, 0)
	if len(v) <= 0 {
		return result
	}
	for eachKey, eachValue := range v {
		result[eachKey] = eachValue
	}
	return result
}
