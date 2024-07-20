package slicex

import (
	"fmt"
)

// conver []T to []interface{}
func ToInterfaceSlice[T any](list []T) []interface{} {
	result := make([]interface{}, 0)
	if len(list) <= 0 {
		return result
	}
	for index := range list {
		result = append(result, list[index])
	}
	return result
}

func ToSliceV[SV any, DV any](list []SV, valueSelectFunc func(item SV) DV) []DV {
	result := make([]DV, 0)
	for _, eachSV := range list {
		currentDV := valueSelectFunc(eachSV)
		result = append(result, currentDV)
	}
	return result
}

// convert []interface{} to []DV
// if any array element cannot convert to DV,then return error
func MapFromInterface[DV any](list []interface{}) ([]DV, error) {
	result := make([]DV, 0)
	for _, eachSV := range list {
		currentDV, ok := eachSV.(DV)
		if !ok {
			return nil, fmt.Errorf("value must be type:%T", *new(DV))
		}

		result = append(result, currentDV)
	}
	return result, nil
}
