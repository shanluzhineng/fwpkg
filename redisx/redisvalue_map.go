package redisx

type RedisValueMap map[string]IRedisValue

// 将map[string]IRedis转换为一个对象数组
// filterFn: 转换时是否执行筛选,如果为nil,则不执行任何筛选
func RedisValueMapToSlice[V any](vMap RedisValueMap, filterFn func(*V) bool) ([]*V, error) {
	valueList := make([]*V, 0)
	for _, eachValue := range vMap {
		currentValue := new(V)
		err := _unmarshal(eachValue.Bytes(), currentValue)
		if err != nil {
			return nil, err
		}
		if filterFn != nil && !filterFn(currentValue) {
			continue
		}
		valueList = append(valueList, currentValue)
	}
	return valueList, nil
}

func RedisValueMapToMap[V any](vMap RedisValueMap, filterFn func(*V) bool) (map[string]*V, error) {
	valueMap := make(map[string]*V)
	for eachKey, eachValue := range vMap {
		currentValue := new(V)
		err := _unmarshal(eachValue.Bytes(), currentValue)
		if err != nil {
			return nil, err
		}
		if filterFn != nil && !filterFn(currentValue) {
			continue
		}
		valueMap[eachKey] = currentValue
	}
	return valueMap, nil
}
