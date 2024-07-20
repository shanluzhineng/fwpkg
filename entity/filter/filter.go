package filter

import (
	"reflect"
)

type IFilter interface {
	GetIsOr() (isOr bool)
	SetIsOr(isOr bool)
	GetConditions() (conditions []FilterCondition)
	SetConditions(conditions []FilterCondition)
	IsNil() (ok bool)
}

type Filter struct {
	IsOr       bool         `form:"is_or" url:"is_or"`
	Conditions []*Condition `json:"conditions"`
}

// #region IFilter members

func (f *Filter) GetIsOr() (isOr bool) {
	return f.IsOr
}

func (f *Filter) SetIsOr(isOr bool) {
	f.IsOr = isOr
}

func (f *Filter) GetConditions() (conditions []FilterCondition) {
	for _, c := range f.Conditions {
		conditions = append(conditions, c)
	}
	return conditions
}

func (f *Filter) SetConditions(conditions []FilterCondition) {
	f.Conditions = make([]*Condition, len(conditions))
	for _, c := range conditions {
		f.Conditions = append(f.Conditions, c.(*Condition))
	}
}

func (f *Filter) IsNil() (ok bool) {
	val := reflect.ValueOf(f)
	return val.IsNil()
}

// #endregion
