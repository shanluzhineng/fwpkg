package slicex

// 查找slice中是否在切片中
func In[K comparable](s []K, item K) bool {
	for _, e := range s {
		if e == item {
			return true
		}
	}
	return false
}

// 查找slice中是否在切片中,通过传入的回调来进行比较
func InWithFunc[V any](s []V, equalsFunc func(item V) bool) bool {
	for _, e := range s {
		if equalsFunc(e) {
			return true
		}
	}
	return false
}

// 查找在slice中的索引
func FindIndex[V any](s []V, equalsFunc func(item V) bool) int {
	for index, e := range s {
		if equalsFunc(e) {
			return index
		}
	}
	return -1
}
