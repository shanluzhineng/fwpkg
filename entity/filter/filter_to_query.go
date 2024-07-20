package filter

import "strings"

// FilterToQuery Translate entity.Filter to bson.M
func FilterToQuery(f IFilter) (q map[string]interface{}, err error) {
	if f == nil || f.IsNil() {
		return nil, nil
	}

	q = map[string]interface{}{}
	for _, cond := range f.GetConditions() {
		key := cond.GetKey()
		op := cond.GetOp()
		value := cond.GetValue()
		switch op {
		case FilterOpNotSet:
			// do nothing
		case FilterOpEqual:
			v := cond.GetValue()
			// if v == nil {
			// 	continue
			// }
			// vString, ok := v.(string)
			// if ok && len(vString) <= 0 {
			// 	continue
			// }
			q[key] = v
		case FilterOpNotEqual:
			q[key] = map[string]interface{}{"$ne": value}
		case FilterOpContains, FilterOpRegex, FilterOpSearch:
			regexV := handlerRegexValue(value)
			q[key] = map[string]interface{}{"$regex": regexV, "$options": "i"}
		case FilterOpNotContains:
			q[key] = map[string]interface{}{"$not": map[string]interface{}{"$regex": value}}
		case FilterOpIn:
			q[key] = map[string]interface{}{"$in": value}
		case FilterOpNotIn:
			q[key] = map[string]interface{}{"$nin": value}
		case FilterOpGreaterThan:
			q[key] = map[string]interface{}{"$gt": value}
		case FilterOpGreaterThanEqual:
			q[key] = map[string]interface{}{"$gte": value}
		case FilterOpLessThan:
			q[key] = map[string]interface{}{"$lt": value}
		case FilterOpLessThanEqual:
			q[key] = map[string]interface{}{"$lte": value}
		default:
			return nil, ErrorFilterInvalidOperation
		}
	}
	return q, nil
}

func handlerRegexValue(v interface{}) interface{} {
	s, ok := v.(string)
	if !ok {
		return v
	}
	s = strings.ReplaceAll(s, "+", "\\+")
	return s
}
