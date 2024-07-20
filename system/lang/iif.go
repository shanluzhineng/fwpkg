package lang

type Empty struct{}

// 类似于c#中的? :赋值运算，返回一个指针
func IfValuePtr[K any](cond bool, trueValue func() *K, falseValue *K) *K {
	if cond {
		return trueValue()
	}
	return falseValue
}

// 类似于c#中的? :赋值运算，返回一个值
func IfValue[K any](cond bool, trueValueCallback func() K, falseValue K) K {
	if cond {
		return trueValueCallback()
	}
	return falseValue
}
