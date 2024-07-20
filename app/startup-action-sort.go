package app

import "sort"

//对startupActionInfo数组进行排序
type startupActionInfoComparer struct {
	List  []*startupActionInfo
	IsAsc bool
}

// 创建一个新的 startupActionInfoComparer对象
func newStartupActionInfoComparer(list []*startupActionInfo, isAsc bool) startupActionInfoComparer {
	var listComparer = startupActionInfoComparer{
		List:  list,
		IsAsc: isAsc,
	}
	return listComparer
}

//#region sort.Interface实现

func (s startupActionInfoComparer) Len() int {
	return len(s.List)
}

func (s startupActionInfoComparer) Swap(i, j int) {
	s.List[i], s.List[j] = s.List[j], s.List[i]
}

func (s startupActionInfoComparer) Less(i, j int) bool {
	if s.IsAsc {
		return s.List[i].priority < s.List[j].priority
	}
	return s.List[i].priority > s.List[j].priority
}

//#endregion

//对列表进行排序
func (s startupActionInfoComparer) Sort() {
	sort.Sort(s)
}
