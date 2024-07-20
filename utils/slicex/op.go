package slicex

// 删除slice中的元素并且返回新的元素
func RemoveByIndex[V any](slice []V, s int) []V {
	return append(slice[:s], slice[s+1:]...)
}

// 对一个slice去重，并且将不符合条件的项移除
func Distinct[K comparable](s []K, predicateFunc func(item K) bool) []K {
	keys := make(map[K]bool)
	list := []K{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			predicateValue := true
			if predicateFunc != nil {
				predicateValue = predicateFunc(entry)
			}
			if !predicateValue {
				//不符合条件，继续
				continue
			}
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// 排除一个项目中符合条件的成员，并返回一个新的slice
func Except[K any](s []K, exceptFunc func(item K) bool) []K {
	set := make([]K, 0)
	for _, eachItem := range s {
		if exceptFunc(eachItem) {
			//被排除的
			continue
		}
		set = append(set, eachItem)
	}
	return set
}
